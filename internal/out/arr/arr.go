package arr

import (
	"strconv"
	"strings"

	"github.com/edanko/moses/internal/config"
	"github.com/edanko/moses/internal/models"
)

var newLine = "\r\n"

/*


%P0180x11.0%

<Array Data>
[PartsArray Info]
ProfileSpec=P0180x11.0
Quality=A40
ArrayName=02482_GVRZ0001
CoderId=
WorkedStatus=0
BaseMatLen = 12000
KeepSide = 5
ArrayDir = 0
ArrayOrder = 3
RemainName=
RemainLen = 170
ScrapName=ScrapОстаток_материала02482_GVRZ0001
ScrapLen = 170
[Part01]
Name = 06006-S06006-UZ_SEK_88305_(14)
Posi = 5.0
Dist = 20.0
[Part02]
Name = 06006-S06006-UZ_SEK_88305_(9)
Posi = 2029.9
Dist = 20.0
[Part03]
Name = 06006-S06006-UZ_SEK_88305_(10)
Posi = 4484.9
Dist = 20.0
[Part04]
Name = 06006-S06006-UZ_SEK_88305_(11)
Posi = 6939.9
Dist = 20.0
[Part05]
Name = 06006-S06006-UZ_SEK_88304_(1)
Posi = 9394.9
Dist = 20.0

<Part01>
[Part Info]
ProfileSpec=P0180x11.0
Quality=A40
PartName=06006-S06006-UZ_SEK_88305_(14)
PartLength=2004.9
RealPartLen=2004.9
PartUnitType=0x0301
ExportType=0
FileVerion=1.1

[EndLeft]
Param1 = (27,0,2742,42.0,0.0,0.0,0.0,30.0,42.0,90.0,40.0,0.0,0.0,0.0,0.0,0.0,0,0,0,0)

[EndRight]
Param1 = (21,0,2131,0.0,0.0,0.0,0.0,30.0,0.0,90.0,40.0,0.0,0.0,-227.0,0.0,0.0,0,0,0,0)

<Part02>
[Part Info]
ProfileSpec=P0180x11.0
Quality=A40
PartName=06006-S06006-UZ_SEK_88305_(9)
PartLength=2435.0
RealPartLen=2435.0
PartUnitType=0x0301
ExportType=0
FileVerion=1.1

[EndLeft]
Param1 = (21,0,2131,0.0,0.0,0.0,0.0,30.0,0.0,90.0,40.0,0.0,0.0,273.0,0.0,0.0,0,0,0,0)

[EndRight]
Param1 = (21,0,2131,0.0,0.0,0.0,0.0,30.0,0.0,90.0,40.0,0.0,0.0,116.0,0.0,0.0,0,0,0,0)

<Part03>
[Part Info]
ProfileSpec=P0180x11.0
Quality=A40
PartName=06006-S06006-UZ_SEK_88305_(10)
PartLength=2435.0
RealPartLen=2435.0
PartUnitType=0x0301
ExportType=0
FileVerion=1.1

[EndLeft]
Param1 = (21,0,2131,0.0,0.0,0.0,0.0,30.0,0.0,90.0,40.0,0.0,0.0,273.0,0.0,0.0,0,0,0,0)

[EndRight]
Param1 = (21,0,2131,0.0,0.0,0.0,0.0,30.0,0.0,90.0,40.0,0.0,0.0,116.0,0.0,0.0,0,0,0,0)

<Part04>
[Part Info]
ProfileSpec=P0180x11.0
Quality=A40
PartName=06006-S06006-UZ_SEK_88305_(11)
PartLength=2435.0
RealPartLen=2435.0
PartUnitType=0x0301
ExportType=0
FileVerion=1.1

[EndLeft]
Param1 = (21,0,2131,0.0,0.0,0.0,0.0,30.0,0.0,90.0,40.0,0.0,0.0,273.0,0.0,0.0,0,0,0,0)

[EndRight]
Param1 = (21,0,2131,0.0,0.0,0.0,0.0,30.0,0.0,90.0,40.0,0.0,0.0,116.0,0.0,0.0,0,0,0,0)

<Part05>
[Part Info]
ProfileSpec=P0180x11.0
Quality=A40
PartName=06006-S06006-UZ_SEK_88304_(1)
PartLength=2435.0
RealPartLen=2435.0
PartUnitType=0x0301
ExportType=0
FileVerion=1.1

[EndLeft]
Param1 = (21,0,2131,0.0,0.0,0.0,0.0,30.0,0.0,90.0,40.0,0.0,0.0,273.0,0.0,0.0,0,0,0,0)

[EndRight]
Param1 = (21,0,2131,0.0,0.0,0.0,0.0,30.0,0.0,90.0,40.0,0.0,0.0,116.0,0.0,0.0,0,0,0,0)








*/

func MGF(b *models.InBar) string {
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

		//left, lscallop := Params(b.FullType, strings.Join(p.Left, " "))

		//o.WriteString(left)

		o.WriteString("END_OF_LEFT_END")
		o.WriteString(newLine)

		leftPoint += p.Length

		o.WriteString("RIGHT_END")
		o.WriteString(newLine)
		o.WriteString("LEFT_FARTHEST_POINT=")
		o.WriteString(strconv.FormatFloat(leftPoint, 'f', -1, 64))
		o.WriteString(newLine)

		//right, rscallop := Params(b.FullType, strings.Join(p.Left, " "))

		//o.WriteString(right)

		o.WriteString("END_OF_RIGHT_END")
		o.WriteString(newLine)

		leftPoint += config.Spacing(p.Right[0], b.FullType)

		/* if lscallop != "" {
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
		} */

		if len(p.Icuts) > 0 {
			for _, icut := range p.Icuts {
				o.WriteString("START_OF_ICUT")
				o.WriteString(newLine)

				//icutString := IcutString(icut)
				o.WriteString(icut)

				o.WriteString("END_OF_ICUT")
				o.WriteString(newLine)
			}
		}
	}

	return o.String()
}
