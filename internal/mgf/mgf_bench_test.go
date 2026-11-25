package mgf

import (
	"path/filepath"
	"testing"

	"github.com/edanko/moses/internal/config"
	"github.com/edanko/moses/internal/domain"
	"github.com/edanko/moses/internal/source/txt"
)

func benchmarkNests(b *testing.B) []domain.Nest {
	b.Helper()

	files, err := filepath.Glob("../../testdata/txt/*.txt")
	if err != nil {
		b.Fatal(err)
	}

	nests := make([]domain.Nest, 0, 128)
	for _, file := range files {
		parsed, err := txt.ParseFile(file)
		if err != nil {
			b.Fatalf("parse %s: %v", file, err)
		}
		nests = append(nests, parsed...)
	}
	if len(nests) == 0 {
		b.Fatal("no nests")
	}
	return nests
}

func BenchmarkFormat(b *testing.B) {
	cfg, err := config.Parse([]byte("version: 2\n"))
	if err != nil {
		b.Fatalf("config: %v", err)
	}
	formatter := New(cfg, cfg)
	nests := benchmarkNests(b)

	b.ReportAllocs()
	var sink int
	for b.Loop() {
		for index, nest := range nests {
			content, err := formatter.Format(nest, index+1)
			if err != nil {
				b.Fatal(err)
			}
			sink += len(content)
		}
	}
	_ = sink
}
