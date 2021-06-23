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

func (b *Bar) UsedLengthString() string {
	return strconv.FormatFloat(b.Length, 'f', -1, 64)
}

func (b *Bar) PartsLength() float64 {
	var res float64

	for _, p := range b.Parts {
		res += p.Length
	}

	return res
}

func (b *Bar) PartsLengthString() string {
	return strconv.FormatFloat(b.PartsLength(), 'f', -1, 64)
}

func (b *Bar) CapacityString() string {
	return strconv.FormatFloat(b.Capacity, 'f', -1, 64)
}

func (b *Bar) RemnantLength() float64 {
	if b.HasUsefulRemnant() {
		return b.Capacity - b.Length
	}
	return 0
}

func (b *Bar) RemnantLengthString() string {
	return strconv.FormatFloat(b.RemnantLength(), 'f', -1, 64)
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

func (b *Bar) Scrap() float64 {
	var res float64

	for _, p := range b.Parts {
		res += p.FullLength - p.Length
	}

	if !b.HasUsefulRemnant() {
		res += b.Capacity - b.Length
	}

	return res
}

func (b *Bar) HasUsefulRemnant() bool {
	return b.Capacity-b.Length > viper.GetFloat64("mincut")
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
	o.WriteString(b.PartsLengthString())

	o.WriteString(fmt.Sprintf(" (%.2f %%)", b.PartsLength()/(b.Capacity-b.RemnantLength())*100))

	if b.HasUsefulRemnant() {
		o.WriteString(" / Отход: ")
		o.WriteString(b.RemnantLengthString())
	}

	o.WriteString(" / Лом: ")
	o.WriteString(fmt.Sprintf("%g", b.Scrap()))

	o.WriteString("\r\n")

	if b.HasUsefulRemnant() {
		o.WriteString("Маркировать отход  ")
		o.WriteString(b.NestName())
		o.WriteString("-L")
		o.WriteString(b.RemnantLengthString())
		o.WriteString("\r\n")
	}

	o.WriteString("----------------------------\r\n| Секция | Позиция | Длина |\r\n----------------------------\r\n")

	for _, p := range b.Parts {
		o.WriteString(fmt.Sprintf("| %6s | %7s | %5s |\r\n", p.Section, p.PosNo, p.LengthString()))
	}
	o.WriteString("----------------------------\r\n")

	o.WriteString("\r\n")

	return o.String()
}
