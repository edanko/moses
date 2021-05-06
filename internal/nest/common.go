package nest

import (
	"fmt"
	"strings"

	"github.com/edanko/moses/internal/config"
	"github.com/edanko/moses/internal/models"

	"github.com/spf13/viper"
)

type Nest struct {
	Bars  []*models.Bar
	Parts []*models.Part
}

func New() *Nest {
	return &Nest{
		Bars:  make([]*models.Bar, 0, 1),
		Parts: make([]*models.Part, 0, 1),
	}
}

func (n *Nest) AddBar(bars ...*models.Bar) {
	n.Bars = append(n.Bars, bars...)
}

func (n *Nest) RemoveBar(idx int) {
	n.Bars[idx] = n.Bars[len(n.Bars)-1]
	n.Bars = n.Bars[:len(n.Bars)-1]
}

func (n *Nest) AddPart(parts ...*models.Part) {
	n.Parts = append(n.Parts, parts...)
}

func (n *Nest) TxtOutputString() string {
	res := strings.Builder{}
	res.Grow(1000)

	for i, b := range n.Bars {
		if i != 0 {
			res.WriteString(".\r\n")
		}
		res.WriteString(b.CapacityString())
		res.WriteString("\r\n")

		for _, p := range b.Parts {
			// new line between parts
			res.WriteString("\r\n")

			if viper.GetBool("withsection") {
				res.WriteString(p.Section)
				res.WriteString("-")
				res.WriteString(p.PosNo)
				res.WriteString("\r\n")
			} else {
				res.WriteString(p.PosNo)
				res.WriteString("\r\n")
			}

			res.WriteString(p.LengthString())
			res.WriteString("\r\n")
			res.WriteString(p.LEnd)
			res.WriteString("\r\n")
			res.WriteString(p.REnd)
			res.WriteString("\r\n")

			// icuts
			for _, icut := range p.Icuts {
				res.WriteString(icut)
				res.WriteString("\r\n")
			}
		}
	}
	return res.String()
}

func (n *Nest) TxtFileNameString() string {
	p := n.Bars[0].Parts[0]

	res := &strings.Builder{}

	if p.Project != "" {
		res.WriteString(p.Project)
		res.WriteRune('-')
	}
	res.WriteString(p.Section)
	res.WriteRune('_')
	res.WriteString(strings.ToUpper(p.Dim))
	if p.Quality != "" {
		res.WriteRune('_')
		res.WriteString(strings.ToUpper(p.Quality))
	}

	return res.String()
}

func (n *Nest) NestingListString() string {
	res := strings.Builder{}
	res.Grow(10000)

	for _, b := range n.Bars {
		res.WriteString(b.String())
	}

	return res.String()
}

func (n *Nest) BarListString() string {
	var stockBarsNum int
	var usedRemnantsNum int

	rems := make([]*models.Bar, 0, 10)

	for _, b := range n.Bars {
		if len(b.Parts) > 0 {
			if b.Capacity != config.BarSize(b.Dim()) {
				usedRemnantsNum++
				rems = append(rems, b)
			} else {
				stockBarsNum++
			}
		}
	}

	l := n.Bars[0]

	res := strings.Builder{}
	res.Grow(100)

	res.WriteString(fmt.Sprintf("%10s", strings.ToUpper(l.Dim())))
	res.WriteString(fmt.Sprintf(" |%6s |", strings.ToUpper(l.Quality())))
	if stockBarsNum > 0 {
		res.WriteString(fmt.Sprintf("%3d bars (%5f)", stockBarsNum, config.BarSize(l.Dim())))
	}
	if stockBarsNum > 0 && usedRemnantsNum > 0 {
		res.WriteString(" and ")
	}
	if usedRemnantsNum > 0 {
		res.WriteString(fmt.Sprintf("%3d remnants (%s", usedRemnantsNum, rems[0].CapacityString()))
		for i := 1; i < len(rems); i++ {
			res.WriteString(" + ")
			res.WriteString(rems[i].CapacityString())
		}
		res.WriteString(")")
	}
	res.WriteString("\r\n")

	return res.String()
}
