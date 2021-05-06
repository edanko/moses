package in

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/edanko/moses/internal/models"
	"github.com/edanko/moses/internal/scan"
)

// Parse txt file to InBar struct
func ProcessTxt(fname string) []*models.InBar {
	var readBarStatus string
	f, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	str, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	s := scan.NewScanner(string(str))

	bars := make([]*models.InBar, 0)

	i := 1
	bar, readBarStatus := readBar(&s, fname)
	bar.I = i
	bar.SetNestName()
	bars = append(bars, bar)

	for {
		i++
		if readBarStatus == "eof" {
			break
		}
		if readBarStatus == "next" {
			var bar *models.InBar
			bar, readBarStatus = readBar(&s, fname)
			bar.I = i
			bar.SetNestName()
			bars = append(bars, bar)
		}
	}
	return bars
}

//func readBar(s *scanner.Scanner, file string) (models.InBar, string) {
func readBar(s *scan.Scan, file string) (*models.InBar, string) {
	var readBarStatus string

	fbase := strings.TrimSuffix(filepath.Base(file), ".txt")

	bar := models.NewInBar()

	bar.SetSection(fbase)
	bar.SetProject()
	bar.SetFullType(fbase)
	bar.SetMaterial(fbase)

	// read raw length
	var err error
	bar.RawLength, err = strconv.ParseFloat(s.ReadNextLine(), 64)
	if err != nil {
		panic(err)
	}
	//skip newline between raw length and first part
	_ = s.ReadNextLine()

	var readPartStatus string

	part, readPartStatus := readPart(s)
	bar.AddInPart(part)

	for {
		if readPartStatus == "eof" {
			readBarStatus = "eof"
			break
		}
		if readPartStatus == "next" {
			var part *models.InPart

			part, readPartStatus = readPart(s)
			// not add new part if newlines at the end of input files
			if part.PosNo == "" || part.PosNo == "eof" { // case for double (and more) newline at end of input file
				continue
			}
			bar.AddInPart(part)
			//parts = append(parts, &part)
			continue
		}
		if readPartStatus == "nextbar" {
			readBarStatus = "next"
			break
		}
	}

	sort.Sort(models.InPartSlice(bar.Parts))

	//bar.SetNestName()

	return bar, readBarStatus
}

func readPart(s *scan.Scan) (*models.InPart, string) {
	part := &models.InPart{}
	// read id
	part.PosNo = strings.TrimSpace(s.ReadNextLine())
	// read length
	part.Length, _ = strconv.ParseFloat(strings.TrimSpace(s.ReadNextLine()), 64)
	// read left endcut
	part.Left = strings.Fields(s.ReadNextLine())
	// read right endcut
	part.Right = strings.Fields(s.ReadNextLine())
	// read the rest until eof
	for {
		text := s.ReadNextLine()
		switch text {
		// eof
		case "eof":
			return part, "eof"
		// read next part
		case "":
			return part, "next"
		// read next bar
		case ".":
			return part, "nextbar"
		// icuts inv for doc
		case "inv":
			part.IcutXInv = true
			continue
		default:
			part.Icuts = append(part.Icuts, text)
		}
	}
}
