package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Bar struct {
	ID            int
	Capacity      float64
	Length        float64
	Parts         []*Part
	UsedRemnantID int
	UsedRemnant   string
	MarkRemnant   string
}

type BarSlice []*Bar

func (bs BarSlice) Len() int { return len(bs) }
func (bs BarSlice) Less(i, j int) bool {
	return bs[i].Capacity < bs[j].Capacity
}
func (bs BarSlice) Swap(i, j int) {
	bs[i], bs[j] = bs[j], bs[i]
}

func NewBar(c float64) *Bar {
	return &Bar{
		ID:            0,
		Capacity:      c,
		Length:        0,
		UsedRemnantID: 0,
		MarkRemnant:   "",
		Parts:         make([]*Part, 0),
	}
}

func (b *Bar) NestName() string {
	p := b.Parts[0]
	res := strings.Builder{}

	res.WriteString(p.Section)
	res.WriteString("-")
	res.WriteString(strings.ToUpper(p.Dim))
	res.WriteString("-")

	if b.Quality() != "" {
		res.WriteString(strings.ToUpper(b.Quality()))
		res.WriteString("-")
	}

	res.WriteString(strconv.Itoa(b.ID))

	return res.String()
}

func (b *Bar) Dim() string {
	return b.Parts[0].Dim
}

func (b *Bar) LengthString() string {
	return strconv.FormatFloat(b.Length, 'f', -1, 64)
}

func (b *Bar) CapacityString() string {
	return strconv.FormatFloat(b.Capacity, 'f', -1, 64)
}

func (b *Bar) RemnantLength() float64 {
	return b.Capacity - b.Length
}

func (b *Bar) RemnantLengthString() string {
	return strconv.FormatFloat(b.Capacity-b.Length, 'f', -1, 64)
}

func (b *Bar) Project() string {
	return b.Parts[0].Project
}

func (b *Bar) Section() string {
	return b.Parts[0].Section
}

func (b *Bar) Quality() string {
	return b.Parts[0].Quality
}

func (b *Bar) PutPart(p *Part) {
	b.Parts = append(b.Parts, p)
	b.Length += p.FullLength
}

func (b *Bar) String() string {
	o := strings.Builder{}
	o.Grow(1000)

	o.WriteString(strings.ToUpper(b.Dim()))
	o.WriteString(" ")
	o.WriteString(strings.ToUpper(b.Quality()))
	o.WriteString("\r\n")

	o.WriteString("Имя программы      ")
	o.WriteString(b.NestName())
	o.WriteString("\r\n")

	if b.UsedRemnantID != 0 {
		o.WriteString("Использовать отход ")
		o.WriteString(b.UsedRemnant)
		o.WriteString("\r\n")
	}

	o.WriteString("Длина: ")
	o.WriteString(b.CapacityString())
	o.WriteString(" / ")
	o.WriteString("Исп.: ")
	o.WriteString(b.LengthString())

	if b.RemnantLength() > 0 {
		o.WriteString(" / Отход: ")
		o.WriteString(b.RemnantLengthString())
	}

	o.WriteString(fmt.Sprintf(" (%.2f %%)", float64(b.Length)/float64(b.Capacity)*100))

	o.WriteString("\r\n")

	if b.Capacity-b.Length > viper.GetFloat64("mincut") {
		o.WriteString("Маркировать отход  ")
		o.WriteString(b.NestName())
		o.WriteString("-L")
		o.WriteString(b.RemnantLengthString())
		o.WriteString("\r\n")
	}

	if viper.GetBool("withsection") {
		o.WriteString("------------------------\r\n| section|  id  |length|\r\n------------------------\r\n")
	} else {
		o.WriteString("---------------\r\n|  id  |length|\r\n---------------\r\n")
	}

	for _, p := range b.Parts {
		if viper.GetBool("withsection") {
			o.WriteString(fmt.Sprintf("|%8s|%6s|%6s|\r\n", p.Section, p.PosNo, p.LengthString()))
		} else {
			o.WriteString(fmt.Sprintf("|%6s|%6s|\r\n", p.PosNo, p.LengthString()))
		}
	}
	if viper.GetBool("withsection") {
		o.WriteString("------------------------\r\n")
	} else {
		o.WriteString("---------------\r\n")
	}
	o.WriteString("\r\n")

	return o.String()
}
