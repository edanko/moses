// Package dxf parses Cadmatic profile-report DXF drawings.
package dxf

import (
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"uuid"

	"github.com/edanko/dxf"
	"github.com/edanko/dxf/entity"

	"github.com/edanko/moses/internal/domain"
	"github.com/edanko/moses/internal/source/scan"
)

var (
	xCoordRE = regexp.MustCompile(`x=(\d+)`)
	angleRE  = regexp.MustCompile(`angle=(-?\d+)`)
)

type rawPart struct {
	project  string
	section  string
	posNo    string
	quality  string
	dim      string
	quantity int
	length   float64
	lEnd     string
	rEnd     string
	icuts    []string
}

// ParseFile reads one DXF profile report.
func ParseFile(path string) (*domain.Part, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", path, err)
	}
	defer func() { _ = file.Close() }()

	profile, err := Parse(file)
	if err != nil {
		return nil, fmt.Errorf("parse %q: %w", path, err)
	}
	return profile, nil
}

// Parse reads one DXF profile report. It returns nil, nil for drawings that
// describe an unsupported profile family.
func Parse(r io.Reader) (*domain.Part, error) {
	drawing, err := dxf.FromReader(r)
	if err != nil && !strings.Contains(err.Error(), "unknown entity type") {
		return nil, fmt.Errorf("read dxf: %w", err)
	}

	if drawing == nil {
		return nil, fmt.Errorf("read dxf: empty drawing")
	}

	var frag11, layer0 strings.Builder
	for _, ent := range drawing.Entities() {
		text, ok := ent.(*entity.Text)
		if !ok {
			continue
		}
		switch ent.Layer().Name() {
		case "FRAG-11":
			frag11.WriteString(text.Value)
			frag11.WriteByte('\n')
		case "0":
			layer0.WriteString(text.Value)
			layer0.WriteByte('\n')
		}
	}

	fragment := frag11.String()
	if len(fragment) < 10 {
		return nil, fmt.Errorf(`drawing has no profile data in layer "FRAG-11"`)
	}
	if strings.HasPrefix(fragment, "LENGTH") {
		fragment = invertFragment(fragment)
	}

	part, err := parseHeader(fragment)
	if err != nil {
		return nil, err
	}
	if part == nil {
		return nil, nil
	}

	if err := parseLayer0(layer0.String(), part); err != nil {
		return nil, err
	}
	return part.toDomain()
}

func invertFragment(fragment string) string {
	lines := strings.Split(fragment, "\n")
	lines = lines[:len(lines)-1]
	for i := len(lines)/2 - 1; i >= 0; i-- {
		opposite := len(lines) - 1 - i
		lines[i], lines[opposite] = lines[opposite], lines[i]
	}
	return strings.Join(lines, "\n") + "\n"
}

func parseHeader(fragment string) (*rawPart, error) {
	scanner := scan.New(fragment)
	part := &rawPart{}

	part.project = strings.TrimSpace(scanner.ReadLine())
	part.section = strings.TrimSpace(scanner.ReadLine())
	_ = scanner.ReadLine() // date
	part.posNo = strings.TrimSpace(scanner.ReadLine())

	quantityText := strings.TrimSpace(scanner.ReadLine())
	quantity, err := strconv.Atoi(quantityText)
	if err != nil {
		return nil, fmt.Errorf("invalid quantity %q: %w", quantityText, err)
	}
	part.quantity = quantity

	_ = scanner.ReadLine() // panel
	part.quality = strings.TrimSpace(scanner.ReadLine())
	part.dim = strings.TrimSpace(scanner.ReadLine())

	if len(part.dim) < 2 {
		return nil, fmt.Errorf("invalid profile type %q", part.dim)
	}
	switch strings.ToLower(part.dim[:2]) {
	case "rp", "fb":
	default:
		return nil, nil
	}

	var ends [2][2][]string
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 6; k++ {
				if i == 0 && j == 0 && k == 0 {
					nextLine := strings.TrimSpace(scanner.ReadLine())
					if strings.Contains(nextLine, "wl") || strings.Contains(nextLine, "vk") {
						_ = scanner.ReadLine()
						nextLine = strings.TrimSpace(scanner.ReadLine())
					}
					ends[j][i] = append(ends[j][i], nextLine)
					continue
				}
				ends[j][i] = append(ends[j][i], strings.TrimSpace(scanner.ReadLine()))
			}
		}
	}

	part.lEnd, err = processEndCut(ends[0][0], ends[0][1])
	if err != nil {
		return nil, fmt.Errorf("left end: %w", err)
	}
	part.rEnd, err = processEndCut(ends[1][0], ends[1][1])
	if err != nil {
		return nil, fmt.Errorf("right end: %w", err)
	}

	_ = scanner.ReadLine() // direction
	lengthLine := scanner.ReadLine()
	lengthParts := strings.SplitN(lengthLine, "=", 2)
	if len(lengthParts) != 2 {
		return nil, fmt.Errorf("invalid length line %q", lengthLine)
	}
	lengthText := strings.TrimSpace(lengthParts[1])
	part.length, err = strconv.ParseFloat(lengthText, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid length %q: %w", lengthText, err)
	}

	return part, nil
}

func parseLayer0(layer0 string, part *rawPart) error {
	if len(layer0) <= 10 {
		return nil
	}

	scanner := scan.New(layer0)
	var previousLine, line, over string
	isBended := false

	addLegLen := func(angle float64) float64 {
		var factor float64
		switch strings.ToLower(part.dim) {
		case "rp100x6":
			factor = 1.75
		case "rp120x6.5":
			factor = 2.09
		case "rp140x7":
			factor = 2.44
		case "rp160x8":
			factor = 2.79
		case "rp180x9":
			factor = 3.14
		case "rp200x10":
			factor = 3.49
		case "rp220x11":
			factor = 3.84
		case "rp240x12":
			factor = 4.19
		}
		return math.Round(factor * angle)
	}

	rightOver := strconv.FormatFloat(part.length-300, 'f', -1, 64)
	rightOver2 := strconv.FormatFloat(part.length-600, 'f', -1, 64)
	rightOver3 := strconv.FormatFloat(part.length-599, 'f', -1, 64)

	for {
		previousLine = line
		line = scanner.ReadLine()

		if line == "" {
			if isBended {
				if strings.Contains(over, "left") {
					part.lEnd = "ba angle=0"
				}
				if strings.Contains(over, "right") {
					part.rEnd = "ba angle=0"
				}
				part.icuts = removeHCCuts(part.icuts)
				break
			}

			if strings.Contains(over, "left") {
				part.length -= 300
				if part.lEnd != "" {
					angle, err := parseEndAngle(part.lEnd)
					if err != nil {
						return err
					}
					if err := shiftCutX(part.icuts, -300); err != nil {
						return err
					}
					if angle > 0 {
						part.length += addLegLen(angle)
						if err := shiftCutX(part.icuts, addLegLen(angle)); err != nil {
							return err
						}
					}
				}
			}
			if strings.Contains(over, "right") {
				part.length -= 300
				if part.rEnd != "" {
					angle, err := parseEndAngle(part.rEnd)
					if err != nil {
						return err
					}
					if angle > 0 {
						part.length += addLegLen(angle)
					}
				}
			}
			break
		}

		if strings.Contains(line, "angle") {
			continue
		}

		switch line {
		case "hc":
			nextLine := scanner.ReadLine()
			if line == nextLine {
				hc := "hc x=" + previousLine + " " + scanner.ReadLine() + " "
				_ = scanner.ReadLine()
				hc += scanner.ReadLine() + " "
				_ = scanner.ReadLine()
				hc += scanner.ReadLine()
				_ = scanner.ReadLine()
				part.icuts = append(part.icuts, strings.ToLower(hc))
				continue
			}
			part.icuts = append(part.icuts, "hc x="+previousLine+" "+strings.ToLower(nextLine+" "+scanner.ReadLine()+" "+scanner.ReadLine()))
			continue
		case "hin":
			nextLine := scanner.ReadLine()
			if line == nextLine {
				hin := "hin x=" + previousLine + " " + scanner.ReadLine() + " "
				_ = scanner.ReadLine()
				hin += scanner.ReadLine()
				_ = scanner.ReadLine()
				part.icuts = append(part.icuts, strings.ToLower(hin))
				_ = scanner.ReadLine()
				continue
			}
			part.icuts = append(part.icuts, "hin x="+previousLine+" "+strings.ToLower(nextLine+" "+scanner.ReadLine()))
			_ = scanner.ReadLine()
			continue
		case "hole":
			nextLine := scanner.ReadLine()
			if line == nextLine {
				part.icuts = append(part.icuts, strings.ToLower("hole x="+previousLine+" "+scanner.ReadLine()))
				_ = scanner.ReadLine()
				continue
			}
			part.icuts = append(part.icuts, "hole x="+previousLine+" "+strings.ToLower(nextLine))
			continue
		case "600", "599":
			over += "left"
			continue
		}

		if line == rightOver2 || line == rightOver3 {
			over += "right"
			continue
		}
		if line == "300" || line == rightOver {
			continue
		}

		if strings.Contains(line, " ") {
			bendText := strings.TrimSpace(scanner.ReadLine())
			bend, err := strconv.Atoi(bendText)
			if err != nil {
				return fmt.Errorf("invalid bend %q: %w", bendText, err)
			}
			if bend > 12 {
				isBended = true
			}
		}
	}

	return nil
}

func parseEndAngle(end string) (float64, error) {
	match := angleRE.FindStringSubmatch(end)
	if len(match) != 2 {
		return 0, nil
	}
	angle, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid angle in %q: %w", end, err)
	}
	return angle, nil
}

func shiftCutX(cuts []string, delta float64) error {
	for i, cut := range cuts {
		match := xCoordRE.FindStringSubmatch(cut)
		if len(match) != 2 {
			continue
		}
		x, err := strconv.ParseFloat(match[1], 64)
		if err != nil {
			return fmt.Errorf("invalid x in %q: %w", cut, err)
		}
		cuts[i] = strings.Replace(cut, "x="+match[1], "x="+strconv.FormatFloat(x+delta, 'f', -1, 64), 1)
	}
	return nil
}

func removeHCCuts(cuts []string) []string {
	result := make([]string, 0, len(cuts))
	for _, cut := range cuts {
		if strings.HasPrefix(cut, "hc") {
			continue
		}
		result = append(result, cut)
	}
	return result
}

func processEndCut(web, flange []string) (string, error) {
	at := func(values []string, index int) string {
		if index < 0 || index >= len(values) {
			return ""
		}
		return values[index]
	}
	webValue := func(index int) string { return at(web, index) }
	flangeValue := func(index int) string { return at(flange, index) }

	switch webValue(0) {
	case "bxag", "bxaf":
		return webValue(0) + " angle=" + webValue(1) + " hprof=" + webValue(2) + " h=" + webValue(3) + " l=" + webValue(4) + " r=" + webValue(3), nil

	case "ba", "ct", "cw", "baf":
		if flangeValue(0) == "bz" {
			return "babz angle=" + webValue(1) + " h2=" + flangeValue(1) + " angle2=" + flangeValue(2), nil
		}
		return "ba angle=" + webValue(1), nil

	case "cup":
		if webValue(1) == "0" {
			return "ba angle=" + webValue(2), nil
		}
		if (webValue(1) == "30" || webValue(1) == "45" || webValue(1) == "40" || webValue(1) == "50") && webValue(2) == "60" {
			return "bc h=" + webValue(1) + " angle=" + webValue(2), nil
		}
		return "", fmt.Errorf("unknown cup end cut: web=%v flange=%v", web, flange)

	case "bk":
		if flangeValue(0) == "bz" {
			return "bkbz h=" + webValue(1) + " angle=" + webValue(2) + " h2=" + flangeValue(1) + " angle2=" + flangeValue(2), nil
		}
		return "bk h=" + webValue(1) + " angle=" + webValue(2), nil

	case "bc", "cun":
		return "bc h=" + webValue(1) + " angle=" + webValue(2), nil

	case "bck":
		return "bck angle=" + webValue(1) + " h=" + webValue(2) + " angle2=" + webValue(3) + " b=" + webValue(4), nil

	case "bcr":
		return "bcr angle=" + webValue(1) + " h=" + webValue(2) + " angle2=" + webValue(3) + " r=" + webValue(4), nil

	case "cx":
		return "cx h=" + webValue(1) + " angle=" + webValue(2), nil

	case "bv":
		if flangeValue(0) == "bz" {
			return "bvbz angle=" + webValue(1) + " r=" + webValue(2) + " h=" + webValue(3) + " b=" + webValue(4) + " h2=" + flangeValue(1) + " angle2=" + flangeValue(2), nil
		}
		return "bv angle=" + webValue(1) + " r=" + webValue(2) + " h=" + webValue(3) + " b=" + webValue(4), nil

	case "bw":
		if webValue(2) != "0" {
			return "bl r=" + webValue(1) + " angle=" + webValue(2), nil
		}
		return "bw r=" + webValue(1) + " angle=" + webValue(2) + " bulb=" + webValue(3), nil

	case "btn":
		return "bl r=" + webValue(3) + " angle=" + webValue(1), nil

	case "bmn":
		return "bmn r=" + webValue(1) + " angle=" + webValue(2) + " b=" + webValue(4), nil

	case "br":
		return "br h=" + webValue(1) + " angle=" + webValue(2) + " bulb=" + webValue(3), nil

	case "bvf":
		return "bvf angle=" + webValue(1) + " r=" + webValue(2) + " h=" + webValue(3) + " b=" + webValue(4), nil

	case "bvg":
		if webValue(1) == "0" && webValue(2) == "0" && webValue(3) == "0" && webValue(4) == "0" {
			return "ba angle=" + webValue(1), nil
		}
		return "bvg angle=" + webValue(1) + " r=" + webValue(2) + " h=" + webValue(3) + " b=" + webValue(4), nil

	case "blg", "bl":
		if flangeValue(0) == "bz" {
			return "blbz r=" + webValue(1) + " angle=" + webValue(2) + " h2=" + flangeValue(1) + " angle2=" + flangeValue(2), nil
		}
		return "bl r=" + webValue(1) + " angle=" + webValue(2), nil

	default:
		return "", fmt.Errorf("unknown end cut: web=%v flange=%v", web, flange)
	}
}

func (p *rawPart) toDomain() (*domain.Part, error) {
	profileType, err := domain.ParseProfileDesignation(p.dim)
	if err != nil {
		return nil, err
	}
	left, err := domain.ParseEndCut(p.lEnd)
	if err != nil {
		return nil, fmt.Errorf("left end: %w", err)
	}
	right, err := domain.ParseEndCut(p.rEnd)
	if err != nil {
		return nil, fmt.Errorf("right end: %w", err)
	}

	holes := make([]domain.Cutout, 0, len(p.icuts))
	for _, cut := range p.icuts {
		if strings.TrimSpace(cut) == "" {
			continue
		}
		hole, err := domain.ParseCutout(cut)
		if err != nil {
			return nil, fmt.Errorf("cut %q: %w", cut, err)
		}
		holes = append(holes, hole)
	}

	profile := domain.NewPart(domain.PartParams{
		ID:          uuid.NewV4(),
		Project:     p.project,
		Section:     p.section,
		PositionNo:  p.posNo,
		Designation: profileType,
		Grade:       p.quality,
		Length:      p.length,
		Quantity:    p.quantity,
		LeftEndCut:  left,
		RightEndCut: right,
		Cutouts:     holes,
	})
	return &profile, nil
}
