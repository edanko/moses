package models

import (
	"sort"

	"github.com/emirpasic/gods/trees/avltree"
)

type byMaterial map[string]*avltree.Tree
type byProfType map[string]byMaterial

func NewPartsStorage() byProfType {
	return make(byProfType)
}

func (a byProfType) Project() string {
	var prj string
loop:
	for _, t := range a.Dims() {
		for _, m := range a[t].Qualitys() {
			partsIterator := a[t][m].Iterator()
			for partsIterator.Next() {
				prj = partsIterator.Value().(*Part).Project
				break loop
			}
		}
	}
	return prj
}

func (a byProfType) Sections() []string {

	sections := make(map[string]bool)

	for _, t := range a.Dims() {
		for _, m := range a[t].Qualitys() {
			partsIterator := a[t][m].Iterator()
			for partsIterator.Next() {
				sections[partsIterator.Value().(*Part).Section] = true
			}
		}
	}

	var res []string
	for sec := range sections {
		res = append(res, sec)
	}

	return res
}

func (m byMaterial) qualitys() []string {
	var res []string
	for q := range m {
		if q != "" {
			res = append(res, q)
		}
	}
	return res
}

func (m byMaterial) Qualitys() []string {
	qs := m.qualitys()

	if !sort.StringsAreSorted(qs) {
		sort.Strings(qs)
	}
	return qs
}

func (a byProfType) dims() []string {
	var res []string
	for d := range a {
		if d != "" {
			res = append(res, d)
		}
	}
	return res
}

func (a byProfType) Dims() []string {
	dims := a.dims()

	if !sort.StringsAreSorted(dims) {
		sort.Strings(dims)
	}
	return dims
}

func (a byProfType) Add(p *Part) {
	d := p.Dim
	if _, ok := a[d]; !ok {
		a[d] = make(byMaterial)
	}

	if _, ok := a[d][p.Quality]; !ok {
		a[d][p.Quality] = avltree.NewWithStringComparator()
	}

	a[d][p.Quality].Put(p.Section+"-"+p.PosNo, p)
}

func (a byProfType) Remove(p *Part) {
	a[p.Dim][p.Quality].Remove(p.Section + "-" + p.PosNo)
}
