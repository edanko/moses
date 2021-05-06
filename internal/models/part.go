package models

import (
	"strconv"
	"strings"

	"github.com/edanko/moses/internal/config"
)

type Part struct {
	Project string
	Section string
	PosNo   string

	Quality string
	Shape   string
	Dim     string

	Length     float64
	FullLength float64

	Quantity int

	LEnd string
	REnd string

	LWebMacro    string
	LWebBevel    string
	LFlangeMacro string
	LFlangeBevel string

	RWebMacro    string
	RWebBevel    string
	RFlangeMacro string
	RFlangeBevel string

	TraceBevel string

	BendingCurve [][]float64

	Icuts []string

	IsBended bool
}

type PartSlice []*Part

func (ps PartSlice) Len() int { return len(ps) }
func (ps PartSlice) Less(i, j int) bool {
	return ps[i].FullLength > ps[j].FullLength
}
func (ps PartSlice) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func (p *Part) LengthString() string {
	return strconv.FormatFloat(p.Length, 'f', -1, 64)
}

func (p *Part) PartHeight() float64 {
	h := strings.Split(p.Dim, "x")[0]

	num, err := strconv.ParseFloat(h[2:], 64)
	if err != nil {
		panic(err)
	}

	return num
}

func (p *Part) SetFullLength() {
	left := config.Spacing(strings.Split(p.LEnd, " ")[0], p.Dim)
	right := config.Spacing(strings.Split(p.REnd, " ")[0], p.Dim)

	p.FullLength = left + p.Length + right
}

func (p *Part) InvertIcuts() {
	for _, i := range p.Icuts {
		ss := strings.Split(i, " ")

		for _, s := range ss {
			if strings.Contains(s, "x=") {
				oldLen, err := strconv.ParseFloat(s[2:], 64)
				if err != nil {
					break
				}

				newLen := p.Length - oldLen
				newLenStr := "x=" + strconv.FormatFloat(newLen, 'f', -1, 64)

				i = strings.Replace(i, s, newLenStr, 1)
			}
		}
	}
}
