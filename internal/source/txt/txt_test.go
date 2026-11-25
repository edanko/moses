package txt

import (
	"path/filepath"
	"strings"
	"testing"
)

const fixtures = "../../../testdata/txt"

func TestParseFilename(t *testing.T) {
	meta, err := ParseFilename("056001-05101_RP140X7_A40.txt")
	if err != nil {
		t.Fatalf("filename: %v", err)
	}
	if meta.Project != "056001" || meta.Section != "05101" {
		t.Fatalf("project/section = %q/%q", meta.Project, meta.Section)
	}
	if meta.Designation.String() != "RP140X7" || meta.Grade != "A40" {
		t.Fatalf("profile/quality = %q/%q", meta.Designation.String(), meta.Grade)
	}
}

func TestParseFilenameRequiresQuality(t *testing.T) {
	if _, err := ParseFilename("056001-05101_RP140X7.txt"); err == nil {
		t.Fatal("expected an error for a file name without quality")
	}
}

func TestParseFixtures(t *testing.T) {
	files, err := filepath.Glob(filepath.Join(fixtures, "*.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 5 {
		t.Fatalf("fixtures = %d, want 5", len(files))
	}

	totalNests := 0
	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			nests, err := ParseFile(file)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			if len(nests) == 0 {
				t.Fatal("no nests parsed")
			}
			totalNests += len(nests)

			for _, nest := range nests {
				if nest.StockBar().Length() <= 0 {
					t.Fatalf("stock length = %v, want positive", nest.StockBar().Length())
				}
				previous := -1.0
				for _, profile := range nest.Parts() {
					if profile.Length() < previous {
						t.Fatalf("profiles are not sorted ascending: %v then %v", previous, profile.Length())
					}
					previous = profile.Length()
				}
			}
		})
	}

	if totalNests != 71 {
		t.Fatalf("total nests = %d, want 71", totalNests)
	}
}

func TestParseSyntheticInversion(t *testing.T) {
	input := strings.Join([]string{
		"3000",
		"",
		"1001",
		"1000",
		"21 v1=0 v2=0",
		"21 v1=0 v2=0",
		"i1 x=100 r=30",
		"inv",
		".",
	}, "\r\n")

	meta, err := ParseFilename("056001-05101_RP160X8_A40.txt")
	if err != nil {
		t.Fatal(err)
	}
	nests, err := Parse(strings.NewReader(input), meta)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(nests) != 1 {
		t.Fatalf("nests = %d, want 1", len(nests))
	}
	holes := nests[0].Parts()[0].Cutouts()
	if len(holes) != 1 {
		t.Fatalf("holes = %d, want 1", len(holes))
	}
	if got := holes[0].Param("x"); got != 900 {
		t.Fatalf("inverted x = %v, want 900", got)
	}
}

func TestParseMalformed(t *testing.T) {
	meta, err := ParseFilename("056001-05101_RP160X8_A40.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Parse(strings.NewReader("not-a-number\n"), meta); err == nil {
		t.Fatal("expected an error for a malformed bar length")
	}
	if _, err := Parse(strings.NewReader(""), meta); err == nil {
		t.Fatal("expected an error for empty input")
	}
}
