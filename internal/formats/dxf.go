package formats

import (
	"fmt"
	"io"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/edanko/moses/internal/models"
	"github.com/edanko/moses/internal/scan"

	"github.com/edanko/dxf"
	"github.com/edanko/dxf/entity"
)

var (
	xCoordRE = regexp.MustCompile(`x=(\d+)`)
	angleRE  = regexp.MustCompile(`angle=(-?\d+)`)
)

func ProcessMssDxf(r io.Reader) *models.Part {
	drawing, err := dxf.FromReader(r)
	if err != nil {
		if !strings.Contains(err.Error(), "unknown entity type") {
			panic(err)
		}
	}

	var frag11 string
	var layer0 string

	for _, ent := range drawing.Entities() {
		switch e := ent.(type) {
		case (*entity.Text):
			switch ent.Layer().Name() {
			case "FRAG-11":
				frag11 += e.Value + "\n"
			case "0":
				layer0 += e.Value + "\n"
			}
		}
	}

	if len(frag11) < 10 {
		panic("PANIC: file doen't have enough information in \"flag11\" layer!")
	}

	// inv test
	if frag11[:6] == "LENGTH" {

		//fmt.Println("DEBUG: inverted direction of all items in the dxf!")

		a := strings.Split(frag11, "\n")

		a = a[:len(a)-1]

		for i := len(a)/2 - 1; i >= 0; i-- {
			opp := len(a) - 1 - i
			a[i], a[opp] = a[opp], a[i]
		}

		var invertedFrag11 string
		for i := range a {
			invertedFrag11 += a[i] + "\n"
		}
		frag11 = invertedFrag11
		//fmt.Println(frag11)
	}

	p := models.Part{}

	//p.Source = path.Base(fname)

	s1 := scan.NewScanner(frag11)

	p.Project = strings.TrimSpace(s1.ReadLine())
	if p.Project == "" {
		fmt.Println("project doesn't read from dxf")
	}

	p.Section = strings.TrimSpace(s1.ReadLine())
	if p.Section == "" {
		fmt.Println("section doesn't read from dxf")
	}

	_ = s1.ReadLine() // skip date

	p.PosNo = strings.TrimSpace(s1.ReadLine())
	if p.PosNo == "" {
		fmt.Println("part id doesn't read from dxf")
	}

	/*if p.IsExcluded() {
		// skip
		return &models.Part{}
	}*/

	p.Quantity, err = strconv.Atoi(strings.TrimSpace(s1.ReadLine()))
	if err != nil {
		panic(err)
	}

	_ = s1.ReadLine() // skip panel

	p.Quality = strings.TrimSpace(s1.ReadLine())
	if p.Quality == "" {
		fmt.Println("quality doensn't read from dxf")
	}

	p.Dim = strings.TrimSpace(s1.ReadLine())
	if p.Dim == "" {
		fmt.Println("profile type doensn't read from dxf")
	}

	// only for rp and fb profiles
	switch p.Dim[:2] {
	default:
		// skip
		return &models.Part{}
	// just continue
	case "rp":
	case "fb":
	}

	var End [2][2][]string

	// i - web, flange
	// j - left, right
	// k - endcut name + 5 params
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 6; k++ {
				if i == 0 && j == 0 && k == 0 {
					nextLine := strings.TrimSpace(s1.ReadLine())
					if strings.Contains(nextLine, "wl") || strings.Contains(nextLine, "vk") {
						//skip length of wl or vk
						s1.ReadLine()
						nextLine = strings.TrimSpace(s1.ReadLine())
					}
					End[j][i] = append(End[j][i], strings.TrimSpace(nextLine))
					continue
				}
				line := strings.TrimSpace(s1.ReadLine())
				End[j][i] = append(End[j][i], line)
			}
		}
	}

	p.LEnd = processEndCut(&p, End[0][0], End[0][1])
	p.REnd = processEndCut(&p, End[1][0], End[1][1])

	_ = s1.ReadLine() // direction ?

	tmpLength := strings.Trim(strings.Split(s1.ReadLine(), "=")[1], " ")
	if tmpLength == "" {
		fmt.Println("length doensn't read from dxf")
	}

	p.Length, err = strconv.ParseFloat(tmpLength, 64)
	if err != nil {
		panic(err)
	}

	//p.ExcessCheck()
	//p.ProcessEndCuts()

	if len(layer0) > 10 {
		s2 := scan.NewScanner(layer0)

		var prevLine string
		var line string
		var over string
		var isBended bool

		iPartLen := p.Length
		rigthOver := strconv.FormatFloat(iPartLen-300, 'f', -1, 64)
		//rigthOver4 := strconv.Itoa(iPartLen - 330)

		rightOver2 := strconv.FormatFloat(iPartLen-600, 'f', -1, 64)
		rightOver3 := strconv.FormatFloat(iPartLen-599, 'f', -1, 64) // wtf??

		addLegLen := func(angle float64) float64 {
			var add float64

			ang := float64(angle)

			switch p.Dim {
			case "rp100x6":
				add = math.Round(1.75 * ang)
			case "rp120x6.5":
				add = math.Round(2.09 * ang)
			case "rp140x7":
				add = math.Round(2.44 * ang)
			case "rp160x8":
				add = math.Round(2.79 * ang)
			case "rp180x9":
				add = math.Round(3.14 * ang)
			case "rp200x10":
				add = math.Round(3.49 * ang)
			case "rp220x11":
				add = math.Round(3.84 * ang)
			case "rp240x12":
				add = math.Round(4.19 * ang)
			}
			return add
		}

		for {
			prevLine = line
			line = s2.ReadLine()

			if line == "" {
				if isBended {
					//fmt.Println(p.Section, p.ID, "bended")

					if strings.Contains(over, "left") {
						p.LEnd = "ba angle=0"
					}
					if strings.Contains(over, "right") {
						p.REnd = "ba angle=0"
					}

					for idx, curIcut := range p.Icuts {
						if curIcut[:2] == "hc" {
							if len(p.Icuts) > 1 {
								if idx >= len(p.Icuts) {
									p.Icuts = p.Icuts[:idx-1]
								} else {
									p.Icuts = append(p.Icuts[:idx], p.Icuts[idx+1:]...)
								}
							} else {
								p.Icuts = make([]string, 0)
							}
						}
					}
				} else {
					if strings.Contains(over, "left") {
						p.Length -= 300
						if p.LEnd != "" {

							angle, err := strconv.ParseFloat(angleRE.FindAllStringSubmatch(p.LEnd, -1)[0][1], 64)
							if err != nil {
								panic(err)
							}

							//remove left excess from icuts
							for idx, curIcut := range p.Icuts {
								xCoordStr := xCoordRE.FindAllStringSubmatch(curIcut, -1)[0][1]
								xCoord, err := strconv.Atoi(xCoordStr)
								if err != nil {
									panic(err)
								}
								newXCoord := xCoord - 300

								p.Icuts[idx] = strings.Replace(curIcut, "x="+xCoordStr, "x="+strconv.Itoa(newXCoord), 1)
							}

							if angle > 0 {
								p.Length += addLegLen(angle)

								// add angle len to icuts
								for idx, curIcut := range p.Icuts {
									xCoordStr := xCoordRE.FindAllStringSubmatch(curIcut, -1)[0][1]
									xCoord, err := strconv.ParseFloat(xCoordStr, 64)
									if err != nil {
										panic(err)
									}
									//newXCoord := xCoord + int(add)
									newXCoord := xCoord + addLegLen(angle)

									p.Icuts[idx] = strings.Replace(curIcut, "x="+xCoordStr, "x="+strconv.FormatFloat(newXCoord, 'f', -1, 64), 1)
								}
							}
						}
					}
					if strings.Contains(over, "right") {
						p.Length -= 300
						if p.REnd != "" {
							angle, err := strconv.ParseFloat(angleRE.FindAllStringSubmatch(p.REnd, -1)[0][1], 64)
							if err != nil {
								panic(err)
							}
							if angle > 0 {
								p.Length += addLegLen(angle)
							}
						}
					}

				}
				break
			}

			// skip unneccesary bending markings
			if strings.Contains(line, "angle") {
				continue
			}

			if line == "hc" {
				nextLine := s2.ReadLine()
				if line == nextLine {
					var hcStr string
					hcStr += "hc x=" + prevLine + " " + s2.ReadLine() + " "
					s2.ReadLine()
					hcStr += s2.ReadLine() + " "
					s2.ReadLine()
					hcStr += s2.ReadLine()
					s2.ReadLine()

					hcStr = strings.ToLower(hcStr)
					p.Icuts = append(p.Icuts, hcStr)
					continue
				}
				p.Icuts = append(p.Icuts, "hc x="+prevLine+" "+strings.ToLower(nextLine+" "+s2.ReadLine()+" "+s2.ReadLine()))
				continue
			}

			if line == "hin" {
				nextLine := s2.ReadLine()
				if line == nextLine {
					var hinStr string
					hinStr += "hin x=" + prevLine + " " + s2.ReadLine() + " "
					s2.ReadLine()
					hinStr += s2.ReadLine()
					s2.ReadLine()

					hinStr = strings.ToLower(hinStr)
					p.Icuts = append(p.Icuts, hinStr)
					// skip last "A=" param in hin definition
					_ = s2.ReadLine()
					continue
				}
				p.Icuts = append(p.Icuts, "hin x="+prevLine+" "+strings.ToLower(nextLine+" "+s2.ReadLine()))
				// skip last "A=" param in hin definition
				_ = s2.ReadLine()
				continue
			}

			if line == "hole" {
				nextLine := s2.ReadLine()
				if line == nextLine {
					var holeStr string
					holeStr += "hole x=" + prevLine + " " + s2.ReadLine()
					s2.ReadLine()

					holeStr = strings.ToLower(holeStr)
					p.Icuts = append(p.Icuts, holeStr)
					continue
				}
				p.Icuts = append(p.Icuts, "hole x="+prevLine+" "+strings.ToLower(nextLine))
				continue
			}

			if line == "600" || line == "599" { // wtf?
				over += "left"
				continue
			}

			if line == rightOver2 || line == rightOver3 {
				over += "right"
				continue
			}

			if line == "300" || line == rigthOver {
				continue
			}

			if strings.Contains(line, " ") {
				bend, err := strconv.Atoi(strings.TrimSpace(s2.ReadLine()))
				if err != nil {
					panic(err)
				}
				if bend > 12 {
					isBended = true
					continue
				}
			}
		}
	}

	p.SetFullLength()

	return &p
}

func processEndCut(p *models.Part, web []string, flange []string) string {

	addBxaLen := func() {
		var a float64

		webL, err := strconv.ParseFloat(web[2], 64)
		if err != nil {
			panic(err)
		}

		diff := p.PartHeight() - webL

		switch diff {
		case 20:
			a += 2
		case 40:
			a += 4
		case 60:
			a += 7
		case 80:
			a += 11
		case 100:
			a += 14
		case 120:
			a += 16
		default:
			panic("some unknown hprof here!!!!")
		}

		p.FullLength += a
	}

	switch web[0] {
	case "bxag":
		addBxaLen()
		return "bxag angle=" + web[1] + " hprof=" + web[2] + " h=" + web[3] + " l=" + web[4] + " r=" + web[3]

	case "bxaf":
		addBxaLen()
		return "bxaf angle=" + web[1] + " hprof=" + web[2] + " h=" + web[3] + " l=" + web[4] + " r=" + web[3]

	case "ba":
		fallthrough
	case "ct":
		fallthrough
	case "cw":
		fallthrough
	case "baf":
		if flange[0] == "bz" {
			return "babz angle=" + web[1] + " h2=" + flange[1] + " angle2=" + flange[2]
		}
		return "ba angle=" + web[1]

	case "cup":
		if web[1] == "0" {
			return "ba angle=" + web[2]
		}
		if (web[1] == "30" || web[1] == "45" || web[1] == "40" || web[1] == "50") && web[2] == "60" {
			return "bc h=" + web[1] + " angle=" + web[2]
		}

		log.Fatalln("some unknown type of cup endcut here, not ba with angle type and not bc type")

	case "bk":
		if flange[0] == "bz" {
			return "bkbz h=" + web[1] + " angle=" + web[2] + " h2=" + flange[1] + " angle2=" + flange[2]
		}
		return "bk h=" + web[1] + " angle=" + web[2]

	case "bc":
		fallthrough
	case "cun":
		return "bc h=" + web[1] + " angle=" + web[2]

	case "bck":
		return "bck angle=" + web[1] + " h=" + web[2] + " angle2=" + web[3] + " b=" + web[4]

	case "bcr":
		return "bcr angle=" + web[1] + " h=" + web[2] + " angle2=" + web[3] + " r=" + web[4]

	case "cx":
		return "cx h=" + web[1] + " angle=" + web[2]

	case "bv":
		if flange[0] == "bz" {
			return "bvbz angle=" + web[1] + " r=" + web[2] + " h=" + web[3] + " b=" + web[4] + " h2=" + flange[1] + " angle2=" + flange[2]
		}
		return "bv angle=" + web[1] + " r=" + web[2] + " h=" + web[3] + " b=" + web[4]

	case "bw":
		if web[2] != "0" {
			return "bl r=" + web[1] + " angle=" + web[2]
		}
		return "bw r=" + web[1] + " angle=" + web[2] + " bulb=" + web[3]

	case "btn":
		return "bl r=" + web[3] + " angle=" + web[1]

	case "bmn":
		return "bmn r=" + web[1] + " angle=" + web[2] + " b=" + web[4]

	case "br":
		return "br h=" + web[1] + " angle=" + web[2] + " bulb=" + web[3]

	case "bvf":
		return "bvf angle=" + web[1] + " r=" + web[2] + " h=" + web[3] + " b=" + web[4]

	case "bvg":
		if web[1] == "0" && web[2] == "0" && web[3] == "0" && web[4] == "0" {
			return "ba angle=" + web[1]
		}
		return "bvg angle=" + web[1] + " r=" + web[2] + " h=" + web[3] + " b=" + web[4]

	case "blg":
		fallthrough
	case "bl":
		if flange[0] == "bz" {
			return "blbz r=" + web[1] + " angle=" + web[2] + " h2=" + flange[1] + " angle2=" + flange[2]
		}
		return "bl r=" + web[1] + " angle=" + web[2]
	default:
		fmt.Printf("unknown endcut - web %s, flange %s\n", web[0], flange[0])
	}
	return ""
}
