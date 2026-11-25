package aveva

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"
)

func writeFixture(t *testing.T, directory string) string {
	t.Helper()

	lst := filepath.Join(directory, "parts.lst")
	if err := os.WriteFile(lst, []byte("Type/Dim. PP/160x8\nQual : A40\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	csvPath := filepath.Join(directory, "parts.csv")
	file, err := os.Create(csvPath)
	if err != nil {
		t.Fatal(err)
	}
	writer := csv.NewWriter(file)

	header := make([]string, 31)
	for i := range header {
		header[i] = "h"
	}
	if err := writer.Write(header); err != nil {
		t.Fatal(err)
	}

	row := make([]string, 31)
	row[0] = "05102-3P-4835A"
	row[3] = "1000"
	row[5] = "21"
	row[6] = "1"
	row[19] = "23"
	row[20] = "2"
	row[26] = "100"
	row[27] = "80"
	if err := writer.Write(row); err != nil {
		t.Fatal(err)
	}
	if err := writer.Write(row); err != nil {
		t.Fatal(err)
	}

	row[0] = "05102-3P-4836B"
	row[3] = "2000"
	if err := writer.Write(row); err != nil {
		t.Fatal(err)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	return csvPath
}

func TestParseFile(t *testing.T) {
	csvPath := writeFixture(t, t.TempDir())

	profiles, err := ParseFile(csvPath, "056001")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(profiles) != 2 {
		t.Fatalf("profiles = %d, want 2", len(profiles))
	}

	first := profiles[0]
	if first.Project() != "056001" || first.Section() != "05102" || first.PositionNo() != "4835" {
		t.Fatalf("identity = %q %q %q", first.Project(), first.Section(), first.PositionNo())
	}
	if first.DesignationName() != "RP160X8" || first.Grade() != "A40" {
		t.Fatalf("dimension = %s %s", first.Designation(), first.Grade())
	}
	if first.Length() != 1000 || first.Quantity() != 2 {
		t.Fatalf("length/quantity = %v/%v", first.Length(), first.Quantity())
	}
	if first.LeftEndKind() != "21" || first.RightEndKind() != "23" {
		t.Fatalf("ends = %q/%q", first.LeftEndKind(), first.RightEndKind())
	}
	if first.LeftEndCut().Param("a") != 1 || first.RightEndCut().Param("a") != 2 {
		t.Fatalf("end params = %v/%v", first.LeftEndCut().Param("a"), first.RightEndCut().Param("a"))
	}
	if first.RightEndCut().Param("v1") != 10 || first.RightEndCut().Param("v2") != 10 {
		t.Fatalf("right angles = %v/%v, want 10/10", first.RightEndCut().Param("v1"), first.RightEndCut().Param("v2"))
	}
	if first.LeftEndCut().Param("v1") != 0 || first.LeftEndCut().Param("v2") != 0 {
		t.Fatalf("right angles leaked into the left end: %v/%v", first.LeftEndCut().Param("v1"), first.LeftEndCut().Param("v2"))
	}
}

func TestMissingLST(t *testing.T) {
	directory := t.TempDir()
	csvPath := filepath.Join(directory, "parts.csv")
	if err := os.WriteFile(csvPath, []byte("header\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ParseFile(csvPath, "056001"); err == nil {
		t.Fatal("expected an error for a missing LST file")
	}
}

func TestRequiresProject(t *testing.T) {
	if _, err := ParseFile("parts.csv", ""); err == nil {
		t.Fatal("expected an error for an empty project")
	}
}
