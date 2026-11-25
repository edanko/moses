package app

import (
	"context"
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"

	"github.com/edanko/dxf/color"
	"github.com/edanko/dxf/drawing"
	"github.com/edanko/dxf/entity"

	"github.com/edanko/moses/internal/config"
	"github.com/edanko/moses/internal/source/txt"
)

func testConfig(t *testing.T) *config.Config {
	t.Helper()
	cfg, err := config.Load("../config/testdata/v2/pipeline.yaml")
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	return cfg
}

func TestConvert10510Gen(t *testing.T) {
	directory := t.TempDir()
	inDir := filepath.Join(directory, "in")
	if err := os.MkdirAll(inDir, 0o755); err != nil {
		t.Fatal(err)
	}
	copyFiles(t, "../../testdata/gen", inDir)

	cfg := testConfig(t)
	options := ConvertOptions{
		Config:       cfg,
		InDir:        inDir,
		OutDir:       filepath.Join(directory, "out"),
		RemnantsPath: filepath.Join(directory, "remnants.csv"),
		Algorithm:    "best-fit-decreasing",
		Project:      "056001",
		WriteExcel:   true,
	}

	summary, err := Convert10510Gen(context.Background(), options)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	if summary.Files != 9 {
		t.Fatalf("files = %d, want 9", summary.Files)
	}
	if summary.Parts != 239 {
		t.Fatalf("parts = %d, want 239", summary.Parts)
	}
	if summary.Nests < 60 {
		t.Fatalf("nests = %d, want at least 60", summary.Nests)
	}
	if summary.Utilization < 85 {
		t.Fatalf("utilization = %.2f, want at least 85", summary.Utilization)
	}

	outputs := []string{
		"056001-05102_bar.txt",
		"056001-05102_nst.txt",
		"056001-05102.xlsx",
	}
	for _, name := range outputs {
		if _, err := os.Stat(filepath.Join(options.OutDir, name)); err != nil {
			t.Fatalf("missing %s: %v", name, err)
		}
	}

	texts, err := filepath.Glob(filepath.Join(options.OutDir, "txt", "*.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if len(texts) != 5 {
		t.Fatalf("txt files = %d, want 5", len(texts))
	}

	roundTripped := 0
	for _, path := range texts {
		nests, err := txt.ParseFile(path)
		if err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}
		roundTripped += len(nests)
	}
	if roundTripped != summary.Nests {
		t.Fatalf("round-tripped nests = %d, want %d", roundTripped, summary.Nests)
	}

	remnants, err := os.ReadFile(options.RemnantsPath)
	if err != nil {
		t.Fatalf("read remnants: %v", err)
	}
	if len(remnants) == 0 {
		t.Fatal("remnant file is empty")
	}
}

func TestConvert10510Synth(t *testing.T) {
	directory := t.TempDir()
	inDir := filepath.Join(directory, "in")
	if err := os.MkdirAll(inDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeAvevaFixture(t, inDir)

	cfg := testConfig(t)
	options := ConvertOptions{
		Config:       cfg,
		InDir:        inDir,
		OutDir:       filepath.Join(directory, "out"),
		RemnantsPath: filepath.Join(directory, "remnants.csv"),
		Algorithm:    "best-fit-decreasing",
		Project:      "056001",
		WriteExcel:   true,
	}

	summary, err := Convert10510(context.Background(), options)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	if summary.Files != 1 || summary.Parts != 1 {
		t.Fatalf("summary = %+v, want one file and one aggregated part", summary)
	}
	if summary.Nests != 1 {
		t.Fatalf("nests = %d, want 1", summary.Nests)
	}

	texts, err := filepath.Glob(filepath.Join(options.OutDir, "txt", "*.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if len(texts) != 1 {
		t.Fatalf("txt files = %d, want 1", len(texts))
	}
	content, err := os.ReadFile(texts[0])
	if err != nil {
		t.Fatal(err)
	}
	if len(content) == 0 {
		t.Fatal("txt output is empty")
	}
}

func writeAvevaFixture(t *testing.T, directory string) {
	t.Helper()

	if err := os.WriteFile(filepath.Join(directory, "parts.lst"), []byte("Type/Dim. PP/160x8\nQual : A40\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	file, err := os.Create(filepath.Join(directory, "parts.csv"))
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
	row[3] = "5000"
	row[5] = "21"
	row[19] = "21"
	if err := writer.Write(row); err != nil {
		t.Fatal(err)
	}
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
}

func TestConvert10022(t *testing.T) {
	directory := t.TempDir()
	inDir := filepath.Join(directory, "in")
	if err := os.MkdirAll(filepath.Join(inDir, "parts"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeDXFFixture(t, filepath.Join(inDir, "parts", "part.dxf"))

	cfg := testConfig(t)
	options := ConvertOptions{
		Config:       cfg,
		InDir:        inDir,
		OutDir:       filepath.Join(directory, "out"),
		RemnantsPath: filepath.Join(directory, "remnants.csv"),
		Algorithm:    "knapsack",
	}

	summary, err := Convert10022(context.Background(), options)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	if summary.Files != 1 || summary.Parts != 1 || summary.Nests != 1 {
		t.Fatalf("summary = %+v", summary)
	}

	texts, err := filepath.Glob(filepath.Join(options.OutDir, "txt", "*.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if len(texts) != 1 {
		t.Fatalf("txt files = %d, want 1", len(texts))
	}
	if _, err := os.Stat(filepath.Join(options.OutDir, "056001-05102_bar.txt")); err != nil {
		t.Fatalf("missing bar report: %v", err)
	}
}

func writeDXFFixture(t *testing.T, path string) {
	t.Helper()

	document, err := drawing.New()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := document.AddLayer("FRAG-11", color.White, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := document.AddLayer("0", color.White, nil); err != nil {
		t.Fatal(err)
	}

	lines := []string{
		"056001", "05102", "date", "4835", "1", "panel", "A40", "RP160X8",
		"ba", "30", "0", "0", "0", "0",
		"ba", "30", "0", "0", "0", "0",
		"", "0", "0", "0", "0", "0",
		"", "0", "0", "0", "0", "0",
		"direction", "LENGTH=1000",
	}
	for index, line := range lines {
		text := entity.NewText()
		text.SetLayerName("FRAG-11")
		text.Value = line
		text.SetCoord([]float64{0, float64(index), 0})
		text.Height = 1
		if err := document.AddEntity(text); err != nil {
			t.Fatal(err)
		}
	}
	if err := document.SaveAs(path); err != nil {
		t.Fatal(err)
	}
}
