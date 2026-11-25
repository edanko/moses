package app

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/edanko/moses/internal/config"
)

func BenchmarkConvert10510Gen(b *testing.B) {
	cfg, err := config.Parse([]byte("version: 2\n"))
	if err != nil {
		b.Fatalf("config: %v", err)
	}

	outDir := b.TempDir()
	remnantsPath := filepath.Join(b.TempDir(), "remnants.csv")
	options := ConvertOptions{
		Config:       cfg,
		InDir:        "../../testdata/gen",
		OutDir:       outDir,
		RemnantsPath: remnantsPath,
		Algorithm:    "best-fit-decreasing",
		Project:      "056001",
		WriteExcel:   true,
	}

	b.ReportAllocs()
	for b.Loop() {
		if err := os.Remove(remnantsPath); err != nil && !os.IsNotExist(err) {
			b.Fatal(err)
		}
		if _, err := Convert10510Gen(context.Background(), options); err != nil {
			b.Fatal(err)
		}
	}
}
