package models

import (
	"strconv"
	"strings"
)

// InBar describes input profile bar
type InBar struct {
	NestName  string
	Project   string
	Section   string
	FullType  string
	Quality   string
	RawLength float64
	I         int
	Parts     []*InPart
}

// InPart describes one part on a bar
type InPart struct {
	PosNo    string
	Length   float64
	Left     []string
	Right    []string
	IcutXInv bool
	Icuts    []string
}

type InPartSlice []*InPart

func (ips InPartSlice) Len() int { return len(ips) }
func (ips InPartSlice) Less(i, j int) bool {
	return ips[i].Length < ips[j].Length
}
func (ips InPartSlice) Swap(i, j int) {
	ips[i], ips[j] = ips[j], ips[i]
}

func NewInBar() *InBar {
	return &InBar{
		Project:   "",
		Section:   "",
		FullType:  "",
		Quality:   "",
		RawLength: 0,
		I:         0,
		Parts:     make([]*InPart, 0),
	}
}

func (b *InBar) AddInPart(p *InPart) {
	b.Parts = append(b.Parts, p)
}

func (b *InBar) SetProject() {
	p := strings.Split(b.Section, "-")
	if len(p) > 1 {
		b.Project = p[0]
		b.Section = p[1]
	}
}

func (b *InBar) SetNestName() {
	res := strings.Builder{}
	res.Grow(50)

	if b.Project != "" {
		res.WriteString(b.Project)
		res.WriteString("-")
	}

	res.WriteString(b.Section)
	res.WriteString("-")
	res.WriteString(strings.ToUpper(b.FullType))
	res.WriteString("-")

	if b.Quality != "" {
		res.WriteString(strings.ToUpper(b.Quality))
		res.WriteString("-")
	}
	res.WriteString(strconv.Itoa(b.I))

	b.NestName = res.String()
}

func (b *InBar) SetSection(s string) {
	p := strings.Split(s, "_")
	b.Section = p[0]
}

func (b *InBar) SetFullType(s string) {
	p := strings.Split(s, "_")

	res := strings.ToLower(p[1])
	switch res {
	case "10":
		res = "rp100x6"
	case "12":
		res = "rp120x6.5"
	case "14a":
		res = "rp140x7"
	case "16a":
		res = "rp160x8"
	case "16b":
		res = "rp160x10"
	case "18a":
		res = "rp180x9"
	case "18b":
		res = "rp180x11"
	case "20a":
		res = "rp200x10"
	case "20b":
		res = "rp200x11"
	case "22a":
		res = "rp220x11"
	case "22b":
		res = "rp220x13"
	case "24a":
		res = "rp240x12"
	case "24b":
		res = "rp240x14"
	}
	b.FullType = res

}

func (b *InBar) SetMaterial(s string) {
	p := strings.Split(s, "_")

	if len(p) > 2 {
		b.Quality = p[2]
	}
}

func (b *InBar) Height() float64 {
	s := strings.Split(b.FullType[2:], "x")[0]

	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}

	return num
}

func (b *InBar) HeightString() string {
	return strings.Split(b.FullType[2:], "x")[0]
}

func (b *InBar) Thickness() float64 {
	s := strings.Split(b.FullType, "x")[1]

	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return num
}
func (b *InBar) ThicknessString() string {
	return strings.Split(b.FullType, "x")[1]
}

func (b *InBar) Type() string {
	return b.FullType[:2]
}
