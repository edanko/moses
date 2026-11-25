package config

import (
	"errors"
	"github.com/edanko/moses/internal/domain"
	"os"
	"path/filepath"
	"testing"
)

func TestParseDefaults(t *testing.T) {
	cfg, err := Parse([]byte("version: 2\n"))
	if err != nil {
		t.Fatalf("parse defaults: %v", err)
	}

	if cfg.Defaults.StockLength != 12000 {
		t.Fatalf("stock length = %v, want 12000", cfg.Defaults.StockLength)
	}
	if !cfg.UseRemnants() {
		t.Fatal("remnants should be enabled by default")
	}
	if cfg.MinRemnantLength() != 1100 {
		t.Fatalf("min remnant = %v, want 1100", cfg.MinRemnantLength())
	}

	spec, err := cfg.Spec("RP160X10")
	if err != nil {
		t.Fatalf("spec: %v", err)
	}
	if spec.Shape != domain.ShapeHP || spec.B != 160 || spec.C != 28 || spec.S != 10 || spec.R != 7 {
		t.Fatalf("unexpected RP spec: %+v", spec)
	}

	fbSpec, err := cfg.Spec("FB80x8")
	if err != nil {
		t.Fatalf("fb spec: %v", err)
	}
	if fbSpec.Shape != domain.ShapeFB || fbSpec.B != 80 || fbSpec.S != 8 {
		t.Fatalf("unexpected FB spec: %+v", fbSpec)
	}

	angleSpec, err := cfg.Spec("AE250x90x10x15")
	if err != nil {
		t.Fatalf("angle spec: %v", err)
	}
	if angleSpec.Shape != domain.ShapeL || angleSpec.A != 250 || angleSpec.T != 15 {
		t.Fatalf("unexpected angle spec: %+v", angleSpec)
	}
}

func TestResolveOverrides(t *testing.T) {
	cfg, err := Parse([]byte(`
version: 2
defaults:
  stockLength: 13000
remnants:
  enabled: false
families:
  hp:
    stockLength: 13000
    endcutSpacing: 7
    endcuts:
      "21":
        profiles:
          rp160x8: 9
`))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	tests := []struct {
		name     string
		profile  string
		endcut   string
		expected float64
	}{
		{name: "profile override", profile: "RP160X8", endcut: "21", expected: 9},
		{name: "family default", profile: "RP160X10", endcut: "23", expected: 7},
		{name: "other family", profile: "FB80X8", endcut: "21", expected: 5},
		{name: "angle family", profile: "AE250X90X10X15", endcut: "21", expected: 5},
		{name: "pp alias", profile: "PP160X8", endcut: "21", expected: 9},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := cfg.Spacing(test.profile, test.endcut)
			if err != nil {
				t.Fatalf("spacing: %v", err)
			}
			if got != test.expected {
				t.Fatalf("spacing = %v, want %v", got, test.expected)
			}
		})
	}

	if got, err := cfg.StockLength("RP160X10"); err != nil || got != 13000 {
		t.Fatalf("stock length = %v, %v; want 13000, nil", got, err)
	}
	if got, err := cfg.StockLength("FB200X20"); err != nil || got != 6000 {
		t.Fatalf("fb stock length = %v, %v; want 6000, nil", got, err)
	}
	if cfg.UseRemnants() {
		t.Fatal("remnants should be disabled by the override")
	}
}

func TestRejectsUnknownFieldsAndVersions(t *testing.T) {
	if _, err := Parse([]byte("version: 2\nunknown: true\n")); !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("unknown field error = %v, want ErrInvalidConfig", err)
	}
	if _, err := Parse([]byte("version: 1\n")); !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("version error = %v, want ErrInvalidConfig", err)
	}
}

func TestNormalizesProfileKeys(t *testing.T) {
	cfg, err := Parse([]byte(`
version: 2
profiles:
  "rp260*12.0":
    family: hp
    shape: HP
    b: 260
    c: 44
    s: 12
    r: 10
`))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if _, err := cfg.Spec("RP260X12"); err != nil {
		t.Fatalf("normalized spec: %v", err)
	}
}

func TestRejectsDuplicateNormalizedProfiles(t *testing.T) {
	if _, err := Parse([]byte(`
version: 2
profiles:
  rp160x8: {family: hp, shape: HP, b: 160, c: 28, s: 8, r: 7}
  "RP160X8": {family: hp, shape: HP, b: 160, c: 28, s: 8, r: 7}
`)); !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("duplicate error = %v, want ErrInvalidConfig", err)
	}
}

func TestMissingSectionGeometry(t *testing.T) {
	cfg, err := Parse([]byte("version: 2\n"))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if _, err := cfg.Spec("RP999X99"); !errors.Is(err, ErrNoSectionGeometry) {
		t.Fatalf("spec error = %v, want ErrNoSectionGeometry", err)
	}
}

func TestLoadFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte("version: 2\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.Source() != path {
		t.Fatalf("source = %q, want %q", cfg.Source(), path)
	}
}
