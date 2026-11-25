// Package gen parses AVEVA marine GEN profile files.
package gen

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"uuid"

	aveva "github.com/edanko/aveva-go"

	"github.com/edanko/moses/internal/domain"
)

// ParseFile reads a GEN profile file and returns its parts.
func ParseFile(path, project string) ([]domain.Part, error) {
	if strings.TrimSpace(project) == "" {
		return nil, fmt.Errorf("project is required")
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", path, err)
	}
	defer func() { _ = file.Close() }()

	document, err := aveva.ReadProfileDocument(file)
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", path, err)
	}
	if document.Common == nil {
		return nil, fmt.Errorf("%q: no common data", path)
	}
	if len(document.Profiles) == 0 {
		return nil, fmt.Errorf("%q: no profiles", path)
	}

	common := document.Common
	shape, err := shapePrefix(common.Shape)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", path, err)
	}

	dimension := strings.Replace(common.Dimension, "*", "X", 1)
	dimension = strings.Replace(dimension, ".0", "", 1)
	profileType, err := domain.ParseProfileDesignation(shape + dimension)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", path, err)
	}

	quality := common.Quality
	if quality == "D40-B" {
		quality = "D40_B"
	}

	profiles := make([]domain.Part, 0, len(document.Profiles))
	for _, part := range document.Profiles {
		if part.LeftEnd == nil || part.RightEnd == nil {
			return nil, fmt.Errorf("%q position %q: missing end cut", path, part.PosNo)
		}

		leftText := endText(part.LeftEnd)
		rightText := endText(part.RightEnd)
		left, err := domain.ParseEndCut(leftText)
		if err != nil {
			return nil, fmt.Errorf("%q position %q left end: %w", path, part.PosNo, err)
		}
		right, err := domain.ParseEndCut(rightText)
		if err != nil {
			return nil, fmt.Errorf("%q position %q right end: %w", path, part.PosNo, err)
		}

		holes := make([]domain.Cutout, 0, len(part.Holes))
		for _, cut := range part.Holes {
			if cut == nil || cut.Name == "" {
				continue
			}
			hole, err := parseCut(cut.Name, cut.DistOrigin)
			if err != nil {
				return nil, fmt.Errorf("%q position %q cut %q: %w", path, part.PosNo, cut.Name, err)
			}
			holes = append(holes, hole)
		}

		profile := domain.NewPart(domain.PartParams{
			ID:          uuid.NewV4(),
			Project:     project,
			Section:     part.BlockNo,
			PositionNo:  part.PosNo,
			Designation: profileType,
			Grade:       quality,
			Length:      part.TLengthManual,
			Quantity:    1,
			LeftEndCut:  left,
			RightEndCut: right,
			Cutouts:     holes,
		})
		profiles = append(profiles, profile)
	}

	return profiles, nil
}

func shapePrefix(shape string) (string, error) {
	if len(shape) < 2 {
		return "", fmt.Errorf("invalid profile shape %q", shape)
	}
	switch shape[:2] {
	case "PP":
		return "RP", nil
	case "FB":
		return "FB", nil
	default:
		return "", fmt.Errorf("unsupported profile shape %q", shape)
	}
}

func endText(end *aveva.End) string {
	var builder strings.Builder
	builder.WriteString(strconv.Itoa(end.EndcutType))

	bevel := ""
	if end.AngleTS != 0 || end.AngleOS != 0 {
		switch end.BevelCode {
		case 177, 154, 180, 227, 143, -245:
		case -131:
			bevel = "-131"
		default:
			if end.BevelCode != 0 {
				bevel = "-" + strconv.Itoa(end.BevelCode)
			}
		}
	}
	builder.WriteString(bevel)

	appendParam(&builder, "a", end.A)
	appendParam(&builder, "b", end.B)
	appendParam(&builder, "c", end.C)
	appendParam(&builder, "r1", end.R1)
	appendParam(&builder, "r2", end.R2)
	if end.V1 != 0 {
		appendValue(&builder, "v1", end.V1-90)
	}
	if end.V2 != 0 {
		appendValue(&builder, "v2", 90-end.V2)
	}
	appendParam(&builder, "v3", end.V3)
	appendParam(&builder, "v4", end.V4)
	appendParam(&builder, "snipe", end.Ks)

	return builder.String()
}

func appendParam(builder *strings.Builder, name string, value float64) {
	if value == 0 {
		return
	}
	appendValue(builder, name, value)
}

func appendValue(builder *strings.Builder, name string, value float64) {
	builder.WriteString(" ")
	builder.WriteString(name)
	builder.WriteString("=")
	builder.WriteString(strconv.FormatFloat(value, 'f', -1, 64))
}

func parseCut(definition string, origin float64) (domain.Cutout, error) {
	if strings.Contains(definition, "=") {
		return domain.ParseCutout(definition)
	}
	params := map[string]float64{"x": origin}
	hole := domain.NewCutout(strings.ToLower(strings.TrimSpace(definition)), params)
	return *hole, nil
}
