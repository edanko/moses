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
