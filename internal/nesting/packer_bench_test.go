package nesting

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/edanko/moses/internal/config"
	"github.com/edanko/moses/internal/domain"
	"github.com/edanko/moses/internal/source/gen"
)

func benchmarkProfiles(b *testing.B) ([]domain.Part, config.Config) {
	b.Helper()

	cfg, err := config.Parse([]byte("version: 2\n"))
	if err != nil {
		b.Fatalf("config: %v", err)
	}

	files, err := filepath.Glob("../../testdata/gen/*.gen")
	if err != nil {
		b.Fatal(err)
	}
	profiles := make([]domain.Part, 0, 256)
	for _, file := range files {
		parsed, err := gen.ParseFile(file, "056001")
		if err != nil {
			b.Fatalf("parse %s: %v", file, err)
		}
		profiles = append(profiles, parsed...)
	}
	if len(profiles) == 0 {
		b.Fatal("no profiles")
	}
	return profiles, *cfg
}

func BenchmarkPack(b *testing.B) {
	profiles, cfg := benchmarkProfiles(b)
	options := PackOptions{
		Spacing:     cfg.Spacing,
		StockLength: cfg.StockLength,
		Algorithm:   "knapsack",
	}

	b.ReportAllocs()
	for b.Loop() {
		if _, err := Pack(context.Background(), profiles, nil, options); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPackSingleGroup(b *testing.B) {
	profiles, cfg := benchmarkProfiles(b)
	group := make([]domain.Part, 0, len(profiles))
	for _, profile := range profiles {
		if profile.DesignationName() == "RP220X11" {
			group = append(group, profile)
		}
	}
	if len(group) == 0 {
		b.Fatal("no RP220X11 profiles")
	}

	options := PackOptions{
		Spacing:     cfg.Spacing,
		StockLength: cfg.StockLength,
		Algorithm:   "knapsack",
	}

	b.ReportAllocs()
	for b.Loop() {
		if _, err := Pack(context.Background(), group, nil, options); err != nil {
			b.Fatal(err)
		}
	}
}
func BenchmarkPackBestFit(b *testing.B) {
	profiles, cfg := benchmarkProfiles(b)
	options := PackOptions{
		Spacing:     cfg.Spacing,
		StockLength: cfg.StockLength,
		Algorithm:   "best-fit-decreasing",
	}

	b.ReportAllocs()
	for b.Loop() {
		if _, err := Pack(context.Background(), profiles, nil, options); err != nil {
			b.Fatal(err)
		}
	}
}
