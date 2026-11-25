package gen

import (
	"path/filepath"
	"testing"
)

const fixtures = "../../../testdata/gen"

func TestParseFixtures(t *testing.T) {
	files, err := filepath.Glob(filepath.Join(fixtures, "*.gen"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no fixtures")
	}

	total := 0
	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			profiles, err := ParseFile(file, "056001")
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			if len(profiles) == 0 {
				t.Fatal("no profiles")
			}
			total += len(profiles)

			for _, profile := range profiles {
				if profile.Project() != "056001" {
					t.Fatalf("project = %q", profile.Project())
				}
				if profile.Section() != "05102" {
					t.Fatalf("section = %q", profile.Section())
				}
				if profile.Grade() == "" {
					t.Fatal("empty quality")
				}
				if profile.LeftEndKind() == "" || profile.RightEndKind() == "" {
					t.Fatalf("missing end cut: %q / %q", profile.LeftEndKind(), profile.RightEndKind())
				}
				if profile.Length() <= 0 {
					t.Fatalf("length = %v", profile.Length())
				}
			}
		})
	}
	if total == 0 {
		t.Fatal("no profiles parsed")
	}
}

func TestRequiresProject(t *testing.T) {
	if _, err := ParseFile(filepath.Join(fixtures, "05102-P2_6700_1S_S4S.gen"), ""); err == nil {
		t.Fatal("expected an error for an empty project")
	}
}
