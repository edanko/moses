// Package txt parses the pre-nested production TXT format.
package txt

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"uuid"

	"github.com/edanko/moses/internal/domain"
	"github.com/edanko/moses/internal/source/scan"
)

// Meta is the metadata encoded in a TXT file name:
// <project>-<section>_<profile>_<quality>.txt.
type Meta struct {
	Project     string
	Section     string
	Designation domain.ProfileDesignation
	Grade       string
}

// ParseFile reads and parses a TXT file using the metadata from its name.
func ParseFile(path string) ([]domain.Nest, error) {
	meta, err := ParseFilename(path)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", path, err)
	}
	defer func() { _ = file.Close() }()

	nests, err := Parse(file, meta)
	if err != nil {
		return nil, fmt.Errorf("parse %q: %w", path, err)
	}
	return nests, nil
}

// ParseFilename extracts the project, section, profile type and quality from
// the production file name.
func ParseFilename(path string) (Meta, error) {
	stem := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	parts := strings.Split(stem, "_")
	if len(parts) < 2 {
		return Meta{}, fmt.Errorf("file name %q must be <section>_<profile>[_<quality>].txt", path)
	}

	project := ""
	section := strings.TrimSpace(parts[0])
	if before, after, ok := strings.Cut(section, "-"); ok {
		project = before
		section = after
	}

	profileType, err := domain.ParseProfileDesignation(parts[1])
	if err != nil {
		return Meta{}, fmt.Errorf("file name %q: %w", path, err)
	}

	quality := ""
	if len(parts) > 2 {
		quality = strings.TrimSpace(parts[2])
	}
	if quality == "" {
		return Meta{}, fmt.Errorf("file name %q does not contain a material quality", path)
	}

	return Meta{
		Project:     project,
		Section:     section,
		Designation: profileType,
		Grade:       quality,
	}, nil
}

// Parse reads bars from r. The caller supplies the metadata parsed from the
// file name.
func Parse(r io.Reader, meta Meta) ([]domain.Nest, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	scanner := scan.New(string(data))
	nests := make([]domain.Nest, 0)

	for {
		nest, status, err := readBar(&scanner, meta)
		if err != nil {
			return nil, err
		}
		if nest != nil {
			nests = append(nests, *nest)
		}
		if status == "eof" {
			break
		}
		if status != "next" {
			return nil, fmt.Errorf("unexpected bar status %q", status)
		}
	}

	if len(nests) == 0 {
		return nil, fmt.Errorf("no bars found")
	}
	return nests, nil
}

func readBar(scanner *scan.Scan, meta Meta) (*domain.Nest, string, error) {
	rawLengthText := strings.TrimSpace(scanner.ReadNextLine())
	if rawLengthText == "eof" {
		return nil, "eof", nil
	}
	rawLength, err := strconv.ParseFloat(rawLengthText, 64)
	if err != nil {
		return nil, "", fmt.Errorf("invalid bar length %q: %w", rawLengthText, err)
	}

	// A blank line separates the raw bar length from the first part.
	_ = scanner.ReadNextLine()

	profiles := make([]domain.Part, 0)
	for {
		profile, status, err := readPart(scanner, meta)
		if err != nil {
			return nil, "", err
		}
		if profile != nil {
			profiles = append(profiles, *profile)
		}
		switch status {
		case "next":
			continue
		case "nextbar":
			if len(profiles) == 0 {
				return nil, "", fmt.Errorf("bar has no profiles")
			}
			sort.Sort(byLength(profiles))
			stock := domain.NewStockBar(domain.StockBarParams{
				ID:          uuid.NewV4(),
				Project:     meta.Project,
				Designation: meta.Designation,
				Grade:       meta.Grade,
				Length:      rawLength,
				Type:        domain.StockTypeBar,
			})
			return domain.NewNest(domain.NestParams{
				ID:       uuid.NewV4(),
				StockBar: stock,
				Parts:    profiles,
			}), "next", nil
		case "eof":
			if len(profiles) == 0 {
				return nil, "eof", nil
			}
			sort.Sort(byLength(profiles))
			stock := domain.NewStockBar(domain.StockBarParams{
				ID:          uuid.NewV4(),
				Project:     meta.Project,
				Designation: meta.Designation,
				Grade:       meta.Grade,
				Length:      rawLength,
				Type:        domain.StockTypeBar,
			})
			return domain.NewNest(domain.NestParams{
				ID:       uuid.NewV4(),
				StockBar: stock,
				Parts:    profiles,
			}), "eof", nil
		default:
			return nil, "", fmt.Errorf("unexpected part status %q", status)
		}
	}
}

func readPart(scanner *scan.Scan, meta Meta) (*domain.Part, string, error) {
	posNo := strings.TrimSpace(scanner.ReadNextLine())
	if posNo == "eof" {
		return nil, "eof", nil
	}

	lengthText := strings.TrimSpace(scanner.ReadNextLine())
	length, err := strconv.ParseFloat(lengthText, 64)
	if err != nil {
		return nil, "", fmt.Errorf("position %q: invalid length %q: %w", posNo, lengthText, err)
	}

	left, err := domain.ParseEndCut(scanner.ReadNextLine())
	if err != nil {
		return nil, "", fmt.Errorf("position %q left end: %w", posNo, err)
	}
	right, err := domain.ParseEndCut(scanner.ReadNextLine())
	if err != nil {
		return nil, "", fmt.Errorf("position %q right end: %w", posNo, err)
	}
	if left.Kind() == "eof" || right.Kind() == "eof" {
		return nil, "", fmt.Errorf("position %q: truncated end cut", posNo)
	}

	holes := make([]domain.Cutout, 0)
	inverted := false

	for {
		line := scanner.ReadNextLine()
		switch line {
		case "eof":
			profile := newProfile(meta, posNo, length, left, right, holes, inverted)
			return &profile, "eof", nil
		case "":
			profile := newProfile(meta, posNo, length, left, right, holes, inverted)
			return &profile, "next", nil
		case ".":
			profile := newProfile(meta, posNo, length, left, right, holes, inverted)
			return &profile, "nextbar", nil
		case "inv":
			inverted = true
		default:
			hole, err := domain.ParseCutout(line)
			if err != nil {
				return nil, "", fmt.Errorf("position %q cut %q: %w", posNo, line, err)
			}
			holes = append(holes, hole)
		}
	}
}

func newProfile(
	meta Meta,
	posNo string,
	length float64,
	left, right domain.EndCut,
	holes []domain.Cutout,
	inverted bool,
) domain.Part {
	if inverted {
		for i := range holes {
			holes[i].SetParam("x", length-holes[i].Param("x"))
		}
	}

	profile := domain.NewPart(domain.PartParams{
		ID:          uuid.NewV4(),
		Project:     meta.Project,
		Section:     meta.Section,
		PositionNo:  posNo,
		Designation: meta.Designation,
		Grade:       meta.Grade,
		Length:      length,
		Quantity:    1,
		LeftEndCut:  left,
		RightEndCut: right,
		Cutouts:     holes,
	})
	return profile
}

type byLength []domain.Part

func (p byLength) Len() int           { return len(p) }
func (p byLength) Less(i, j int) bool { return p[i].Length() < p[j].Length() }
func (p byLength) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
