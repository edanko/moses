package formats

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/edanko/moses/entities"
	"github.com/edanko/moses/internal/scan"
)

var (
	project, section, dimension, quality string
)

func ProcessTxt(fname string) ([]*entities.Nest, error) {
	f, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	str, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	s := scan.NewScanner(string(str))

	spl := strings.Split(string(str), "\n.")
	fmt.Println(spl)

	fbase := strings.Split(strings.TrimSuffix(filepath.Base(fname), ".txt"), "_")

	section = fbase[0]
	p := strings.Split(section, "-")
	if len(p) > 1 {
		project = p[0]
		section = p[1]
	}

	dimension = renameDim(fbase[1])
	if len(fbase) > 2 {
		quality = fbase[2]
	} else {
		project = "08380"
		quality = "E32W"
	}

	nests := make([]*entities.Nest, 0)

	var readBarStatus string

	i := 1
	nest, readBarStatus := readBar(&s)
	//nest.I = i
	//nest.SetNestName()

	nest.Project = project
	//nest.Section = section
	nest.Bar.Dim = dimension
	nest.Bar.Quality = quality

	nests = append(nests, nest)

	for {
		i++
		if readBarStatus == "eof" {
			break
		}
		if readBarStatus == "next" {
			nest, _ := readBar(&s)
			//nest.I = i
			//nest.SetNestName()
			nests = append(nests, nest)
		}
	}
	return nests, err
}

func readBar(s *scan.Scan) (*entities.Nest, string) {
	var readBarStatus string
	nest := &entities.Nest{}
	nest.Bar = &entities.Bar{}

	var err error
	nest.Bar.Length, err = strconv.ParseFloat(s.ReadNextLine(), 64)
	if err != nil {
		panic(err)
	}
	// newline between raw length and first part
	_ = s.ReadNextLine()

	var readPartStatus string

	part, readPartStatus := readPart(s)
	nest.Profiles = append(nest.Profiles, part)

	for {
		if readPartStatus == "eof" {
			readBarStatus = "eof"
			break
		}
		if readPartStatus == "next" {
			part, _ := readPart(s)
			// not add new part if newlines at the end of input files
			if part.PosNo == "" || part.PosNo == "eof" { // case for double (and more) newline at end of input file
				continue
			}
			nest.Profiles = append(nest.Profiles, part)
			continue
		}
		if readPartStatus == "nextbar" {
			readBarStatus = "next"
			break
		}
	}

	return nest, readBarStatus
}

// TODO: merge same parts
func readPart(s *scan.Scan) (*entities.Profile, string) {
	part := &entities.Profile{}
	part.Project = project
	part.Section = section
	part.Quality = quality
	part.Dim = dimension
	part.Quantity = 1

	//part.Source = path.Base(fname)

	part.PosNo = strings.TrimSpace(s.ReadNextLine())
	part.Length, _ = strconv.ParseFloat(strings.TrimSpace(s.ReadNextLine()), 64)
	part.L = &entities.End{}
	part.L.Name = s.ReadNextLine()

	part.R = &entities.End{}
	part.R.Name = s.ReadNextLine()

	var inv bool

	for {
		text := s.ReadNextLine()
		switch text {
		// eof
		case "eof":
			if inv {
				part.InvertHolesX()
			}
			return part, "eof"
		// read next part
		case "":
			if inv {
				part.InvertHolesX()
			}
			return part, "next"
		// read next bar
		case ".":
			if inv {
				part.InvertHolesX()
			}
			return part, "nextbar"
		// icuts inv for doc
		case "inv":
			inv = true
			continue
		default:

			h := &entities.Hole{}
			h.Name = text

			part.Holes = append(part.Holes, h)
		}
	}
}

func renameDim(dim string) string {
	dim = strings.ToUpper(dim)
	switch dim {
	case "10":
		dim = "RP100*6.0"
	case "12":
		dim = "RP120*6.5"
	case "14A":
		dim = "RP140*7.0"
	case "16A":
		dim = "RP160*8.0"
	case "16B":
		dim = "RP160*10.0"
	case "18A":
		dim = "RP180*9.0"
	case "18B":
		dim = "RP180*11.0"
	case "20A":
		dim = "RP200*10.0"
	case "20B":
		dim = "RP200*11.0"
	case "22A":
		dim = "RP220*11.0"
	case "22B":
		dim = "RP220*13.0"
	case "24A":
		dim = "RP240*12.0"
	case "24B":
		dim = "RP240*14.0"
	}
	return dim
}
