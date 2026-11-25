package app

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/edanko/moses/internal/config"
)

func TestParserMatchesGolden(t *testing.T) {
	directory := t.TempDir()
	inDir := filepath.Join(directory, "in")
	if err := os.MkdirAll(inDir, 0o755); err != nil {
		t.Fatal(err)
	}
	copyFiles(t, "../../testdata/txt", inDir)

	cfg, err := config.Load("../config/testdata/v2/parser.yaml")
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	outDir := filepath.Join(directory, "out")
	summary, err := RunParser(context.Background(), ParserOptions{
		Config:     cfg,
		InDir:      inDir,
		OutDir:     outDir,
		DontRename: true,
	})
	if err != nil {
		t.Fatalf("run parser: %v", err)
	}
	if summary.Files != 5 || summary.Nests != 71 || summary.MGFs != 71 {
		t.Fatalf("summary = %+v, want 5 files / 71 nests / 71 mgf", summary)
	}

	goldens, err := filepath.Glob("../../testdata/golden/mgf/*.mgf")
	if err != nil {
		t.Fatal(err)
	}
	if len(goldens) != 71 {
		t.Fatalf("goldens = %d, want 71", len(goldens))
	}

	for _, golden := range goldens {
		name := filepath.Base(golden)
		produced := filepath.Join(outDir, "056001", "05101", name)
		expected, err := os.ReadFile(golden)
		if err != nil {
			t.Fatal(err)
		}
		actual, err := os.ReadFile(produced)
		if err != nil {
			t.Fatalf("missing %s: %v", name, err)
		}
		if !bytes.Equal(actual, expected) {
			t.Fatalf("%s differs from golden (%d vs %d bytes)", name, len(actual), len(expected))
		}
	}
}

func copyFiles(t *testing.T, from, to string) {
	t.Helper()
	entries, err := os.ReadDir(from)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(from, entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(to, entry.Name()), data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
}
