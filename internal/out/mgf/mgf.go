package mgf

import (
	"strconv"
	"strings"

	"github.com/edanko/moses/internal/config"
	"github.com/edanko/moses/internal/models"
)

var newLine = "\r\n"

func StrightCut(b *models.InBar) string {

	check := func(rightEndcut string) bool {
		switch rightEndcut {
		case "fv", "fv.1", "b-102", "b-103", "b-104", "b-106", "sb1", "sb2", "sb6", "sb9", "sb9-01":
			return true
		default:
			return false
		}
	}

	o := strings.Builder{}

	for i, p := range b.Parts {
		if check(p.Right[0]) {
			if len(b.Parts) != i+1 {
				o.WriteString(strconv.Itoa(i + 1))
			}
		}
	}
	return o.String()
}

func Mgf(b *models.InBar) string {
	o := strings.Builder{}

	// header
	o.WriteString("TYPE_OF_GENERIC_FILE=MOSES_NESTED_PROFILE")
	o.WriteString(newLine)
	o.WriteString("VERSION=1.0")
	o.WriteString(newLine)
	o.WriteString("USAGE=PLASMA")
	o.WriteString(newLine)

	// common data
	o.WriteString("COMMON_DATA")
	o.WriteString(newLine)
	o.WriteString("NEST_NAME=")
	o.WriteString(b.NestName)
	o.WriteString(newLine)

	o.WriteString(config.Norm(b.FullType))
	o.WriteString(newLine)

	o.WriteString("RAW_LENGTH=")
	if b.RawLength > 1800 {
		o.WriteString(strconv.FormatFloat(b.RawLength, 'f', -1, 64))
		o.WriteString(newLine)
	} else {
		o.WriteString("1800")
		o.WriteString(newLine)
	}

	var used float64
	for _, p := range b.Parts {
		used += config.Spacing(p.Left[0], b.FullType)
		used += p.Length
		used += config.Spacing(p.Right[0], b.FullType)
	}
	rest := b.RawLength - used

	o.WriteString("USED_LENGTH=")
	o.WriteString(strconv.FormatFloat(used, 'f', -1, 64))
	o.WriteString(newLine)

	o.WriteString("REST_LENGTH=")
	o.WriteString(strconv.FormatFloat(rest, 'f', -1, 64))
	o.WriteString(newLine)

	o.WriteString("TEXT_HEIGHT=20")
	o.WriteString(newLine)

	o.WriteString("TEXT_WIDTH=16")
	o.WriteString(newLine)
	o.WriteString("TEXT_PLANE=0")
	o.WriteString(newLine)
	o.WriteString("TEXT_PLACING=3")
	o.WriteString(newLine)
	o.WriteString("TEXT_U=100")
	o.WriteString(newLine)
	o.WriteString("TEXT_V=80")
	o.WriteString(newLine)
	o.WriteString("MATERIAL=St37-2")
	o.WriteString(newLine)
	o.WriteString("NO_OF_PROFS=")
	o.WriteString(strconv.Itoa(len(b.Parts)))
	o.WriteString(newLine)

	o.WriteString("END_OF_COMMON_DATA")
	o.WriteString(newLine)

	var leftPoint float64

	for i, p := range b.Parts {
		o.WriteString("PROFILE_DATA")
		o.WriteString(newLine)

		o.WriteString("TLENGTH=")
		o.WriteString(strconv.FormatFloat(p.Length, 'f', -1, 64))
		o.WriteString(newLine)

		o.WriteString("IDENT_STRING=")
		o.WriteString(b.Project + "-" + b.Section + "-" + p.PosNo)
		o.WriteString(newLine)

		o.WriteString("NO_OF_MARKS=0")
		o.WriteString(newLine)
		o.WriteString("NO_OF_ICUTS=")
		o.WriteString(strconv.Itoa(len(p.Icuts)))
		o.WriteString(newLine)

		o.WriteString("NO_OF_PARTS=")
		o.WriteString(strconv.Itoa(i + 1))
		o.WriteString(newLine)

		o.WriteString("END_OF_PROFILE_DATA")
		o.WriteString(newLine)

		o.WriteString("LEFT_END")
		o.WriteString(newLine)

		leftPoint += config.Spacing(p.Left[0], b.FullType)

		o.WriteString("LEFT_CLOSEST_POINT=")
		o.WriteString(strconv.FormatFloat(leftPoint, 'f', -1, 64))
		o.WriteString(newLine)

		left, lscallop := Params(b.FullType, strings.Join(p.Left, " "))

		o.WriteString(left)

		o.WriteString("END_OF_LEFT_END")
		o.WriteString(newLine)

		leftPoint += p.Length

		o.WriteString("RIGHT_END")
		o.WriteString(newLine)
		o.WriteString("LEFT_FARTHEST_POINT=")
		o.WriteString(strconv.FormatFloat(leftPoint, 'f', -1, 64))
		o.WriteString(newLine)

		right, rscallop := Params(b.FullType, strings.Join(p.Right, " "))

		o.WriteString(right)

		o.WriteString("END_OF_RIGHT_END")
		o.WriteString(newLine)

		leftPoint += config.Spacing(p.Right[0], b.FullType)

		if lscallop != "" {
			o.WriteString("LEFT_SCALLOP")
			o.WriteString(newLine)

			o.WriteString(lscallop)

			o.WriteString("END_OF_LEFT_SCALLOP")
			o.WriteString(newLine)
		}

		if rscallop != "" {
			o.WriteString("RIGHT_SCALLOP")
			o.WriteString(newLine)

			o.WriteString(rscallop)

			o.WriteString("END_OF_RIGHT_SCALLOP")
			o.WriteString(newLine)
		}

		if len(p.Icuts) > 0 {
			for _, icut := range p.Icuts {
				o.WriteString("START_OF_ICUT")
				o.WriteString(newLine)

				icutString := IcutString(icut)
				o.WriteString(icutString)

				o.WriteString("END_OF_ICUT")
				o.WriteString(newLine)
			}
		}
	}

	return o.String()
}

/* func NewMGF(b *models.InBar) *MGF {
	for i := range b.Parts {
		ip := b.Parts[i]
		op := &MgfPart{}


		left := m.processTemplate(ip.Left[0], b.GetType(), mkMap(ip.Left))
		left2 := strings.Split(left, "@")
		op.LeftEnd = left2[0]
		if len(left2) > 1 {
			op.LeftScallop = left2[1][2:]
		}

		switch ip.Right[0] {
		case "sb1":
			if ip.Length < 1800 && len(ip.Icuts) == 3 {
				ip.Right[0] = "sb1.1"
			}
		case "sb3":
			ip.Right[0] = "sb3.1"
		case "sb8":
			ip.Right[0] = "sb8.1"
		case "sb20":
			ip.Right[0] = "sb20.1"
		case "fv":
			ip.Right[0] = "fv.1"
		}

		right := m.processTemplate(ip.Right[0], b.GetType(), mkMap(ip.Right))
		right2 := strings.Split(right, "@")
		op.RightEnd = right2[0]
		if len(right2) > 1 {
			op.RightScallop = right2[1][2:]
		}

		//leftPoint := stoi(op.LeftPoint)
		length := ip.Length

		// special case. increase length for bxaf and bxag template
		addBxLen := func(in []string) {
			hprof := mkMap(in)

			diff := b.GetHeight() - hprof["hprof"].(float64)
			switch diff {
			case 20:
				length += 2
			case 40:
				length += 4
			case 60:
				length += 7
			case 80:
				length += 11
			case 100:
				length += 14
			case 120:
				length += 16
			}
		}

		switch ip.Left[0] {
		case "bxaf", "bxag":
			addBxLen(ip.Left)
		}
		switch ip.Right[0] {
		case "bxaf", "bxag":
			addBxLen(ip.Right)
		}

		// ------
		//op.RightPoint = strconv.Itoa(leftPoint + length)
		//op.Length = strconv.Itoa(length)

		if ip.Icuts != nil {
			op.IcutXInv = ip.IcutXInv
			op.Icuts = strconv.Itoa(len(ip.Icuts))
			op.Icut = make([]string, len(ip.Icuts))

			for k := range ip.Icuts {
				if op.IcutXInv {
					tmp := strings.Fields(ip.Icuts[k])
					tmpMap := mkMap(tmp)

					// x in map = part length minus x coord in map
					l, err := strconv.ParseFloat(op.Length, 64)
					if err != nil {
						panic(err)
					}
					x := fmt.Sprintf("%.0f", l-tmpMap["x"].(float64))
					tmpMap["x"] = x

					op.Icut[k] = m.processTemplate(tmp[0], "i", tmpMap)
				} else {
					tmp := strings.Fields(ip.Icuts[k])
					op.Icut[k] = m.processTemplate(tmp[0], "i", mkMap(tmp))
				}
			}
		} else {
			op.Icuts = "0"
		}
		m.Profiles = append(m.Profiles, op)

		// check for last cut length
		//lastCutLength = config.GetSpacing(b.Parts[i].Right[0], b.FullType)
		lastCutLength = 0
		// check for over length
		l, _ := strconv.ParseFloat(op.RightPoint, 64)
		if l > b.RawLength {
			fmt.Printf(" [x] error @ bar ХЗ (len:%g), after part id %s (len:%f, need extra:%g)\n", b.RawLength, ip.PosNo, ip.Length, l+lastCutLength-b.RawLength)

			fmt.Printf(" [x] parts:")
			for _, p := range m.Profiles {
				fmt.Printf(" %s (len:%s)", p.ID, p.Length)
			}
			fmt.Println()
			fmt.Println()
		}
	}

} */
