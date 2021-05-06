package in

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/edanko/gen2dxf/pkg/gen"
	"github.com/edanko/moses/internal/models"
)

const (
	prj = "056001"
)

func ProcessGen(filename string) []*models.Part {
	var profs []*models.Part

	g := gen.ParseProfileFile(filename)

	for _, p := range g.Profiles {

		part := &models.Part{}

		part.Project = prj
		part.Section = p.BlockNo
		part.PosNo = p.PosNo

		part.Length = p.TlengthManual

		var shape string

		dim := strings.Replace(g.CommonData.Dimension, "*", "X", 1)
		dim = strings.Replace(dim, ".0", "", 1)

		switch g.CommonData.Shape[:2] {
		case "PP":
			shape = "RP"
		case "FB":
			shape = "FB"
		}

		part.Shape = shape

		part.Dim = shape + dim
		part.Quality = g.CommonData.Quality

		if part.Quality == "D40-B" {
			part.Quality = "D40_B"
		}

		bev := func(e *gen.End) (bev string) {
			if e.AngleTs == 0 && e.AngleOs == 0 {
				return
			}

			if e.BevelCode == 0 {
				return
			}

			switch e.BevelCode {
			case 177, 154, 180, 227, 143, -245:
				bev = ""

			case -131:
				bev = "-131"

			default:
				bev = "-" + strconv.Itoa(e.BevelCode)
			}

			return
		}

		part.LEnd = strconv.Itoa(p.LeftEnd.EndcutType) + bev(p.LeftEnd)

		if p.LeftEnd.A != 0 {
			part.LEnd += " a=" + strconv.FormatFloat(p.LeftEnd.A, 'f', -1, 64)
		}
		if p.LeftEnd.B != 0 {
			part.LEnd += " b=" + strconv.FormatFloat(p.LeftEnd.B, 'f', -1, 64)
		}
		if p.LeftEnd.C != 0 {
			part.LEnd += " c=" + strconv.FormatFloat(p.LeftEnd.C, 'f', -1, 64)
		}
		if p.LeftEnd.R1 != 0 {
			part.LEnd += " r1=" + strconv.FormatFloat(p.LeftEnd.R1, 'f', -1, 64)
		}
		if p.LeftEnd.R2 != 0 {
			part.LEnd += " r2=" + strconv.FormatFloat(p.LeftEnd.R2, 'f', -1, 64)
		}
		if p.LeftEnd.V1 != 0 {
			part.LEnd += " v1=" + strconv.FormatFloat(p.LeftEnd.V1-90, 'f', -1, 64)
		}
		if p.LeftEnd.V2 != 0 {
			part.LEnd += " v2=" + strconv.FormatFloat(90-p.LeftEnd.V2, 'f', -1, 64)
		}
		if p.LeftEnd.V3 != 0 {
			part.LEnd += " v3=" + strconv.FormatFloat(p.LeftEnd.V3, 'f', -1, 64)
		}
		if p.LeftEnd.V4 != 0 {
			part.LEnd += " v4=" + strconv.FormatFloat(p.LeftEnd.V4, 'f', -1, 64)
		}
		if p.LeftEnd.Ks != 0 {
			part.LEnd += " snipe=" + strconv.FormatFloat(p.LeftEnd.Ks, 'f', -1, 64)
		}

		part.REnd = strconv.Itoa(p.RightEnd.EndcutType) + bev(p.RightEnd)

		if p.RightEnd.A != 0 {
			part.REnd += " a=" + strconv.FormatFloat(p.RightEnd.A, 'f', -1, 64)
		}
		if p.RightEnd.B != 0 {
			part.REnd += " b=" + strconv.FormatFloat(p.RightEnd.B, 'f', -1, 64)
		}
		if p.RightEnd.C != 0 {
			part.REnd += " c=" + strconv.FormatFloat(p.RightEnd.C, 'f', -1, 64)
		}
		if p.RightEnd.R1 != 0 {
			part.REnd += " r1=" + strconv.FormatFloat(p.RightEnd.R1, 'f', -1, 64)
		}
		if p.RightEnd.R2 != 0 {
			part.REnd += " r2=" + strconv.FormatFloat(p.RightEnd.R2, 'f', -1, 64)
		}
		if p.RightEnd.V1 != 0 {
			part.REnd += " v1=" + strconv.FormatFloat(p.RightEnd.V1-90, 'f', -1, 64)
		}
		if p.RightEnd.V2 != 0 {
			part.REnd += " v2=" + strconv.FormatFloat(90-p.RightEnd.V2, 'f', -1, 64)
		}
		if p.RightEnd.V3 != 0 {
			part.REnd += " v3=" + strconv.FormatFloat(p.RightEnd.V3, 'f', -1, 64)
		}
		if p.RightEnd.V4 != 0 {
			part.REnd += " v4=" + strconv.FormatFloat(p.RightEnd.V4, 'f', -1, 64)
		}
		if p.RightEnd.Ks != 0 {
			part.REnd += " snipe=" + strconv.FormatFloat(p.RightEnd.Ks, 'f', -1, 64)
		}

		for _, h := range p.HolesNotchesCutouts {
			if h.Name != "" {
				fmt.Println(h.Type, h.Name, "x=", h.DistOrigin)
			}
			part.Icuts = append(part.Icuts, h.Name)
		}

		part.SetFullLength()

		profs = append(profs, part)
	}
	return profs
}
