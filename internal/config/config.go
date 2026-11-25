// Package config owns the versioned moses configuration: loading, merging,
// validation and the spacing/stock/profile lookups used by the pipeline.
package config

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"github.com/edanko/moses/internal/domain"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Version is the supported configuration schema version.
const Version = 2

var (
	// ErrNoSpacing reports that no endcut spacing could be resolved.
	ErrNoSpacing = errors.New("config: no spacing")
	// ErrNoStockLength reports that no stock length could be resolved.
	ErrNoStockLength = errors.New("config: no stock length")
	// ErrNoSectionGeometry reports that no geometry spec could be resolved.
	ErrNoSectionGeometry = errors.New("config: no profile spec")
	// ErrUnknownFamily reports that a profile type has no configured family.
	ErrUnknownFamily = errors.New("config: unknown profile family")
	// ErrInvalidConfig reports a schema or value problem.
	ErrInvalidConfig = errors.New("config: invalid configuration")
)

// Defaults holds values used when a family or profile does not override them.
type Defaults struct {
	StockLength   float64 `yaml:"stockLength"`
	EndcutSpacing float64 `yaml:"endcutSpacing"`
}

// Remnants holds remnant CSV settings.
type Remnants struct {
	Enabled   bool    `yaml:"enabled"`
	MinLength float64 `yaml:"minLength"`
	File      string  `yaml:"file"`
}

// Family groups profile types that share stock and endcut spacing defaults.
type Family struct {
	Prefixes      []string          `yaml:"prefixes"`
	StockLength   *float64          `yaml:"stockLength,omitempty"`
	EndcutSpacing *float64          `yaml:"endcutSpacing,omitempty"`
	Endcuts       map[string]Endcut `yaml:"endcuts,omitempty"`
}

// Endcut holds spacing for one endcut, optionally overridden per profile type.
type Endcut struct {
	Spacing  *float64           `yaml:"spacing,omitempty"`
	Profiles map[string]float64 `yaml:"profiles,omitempty"`
}

// Part declares the geometry and optional overrides for one profile type.
type Part struct {
	Family       string       `yaml:"family,omitempty"`
	Shape        domain.Shape `yaml:"shape"`
	StockLength  *float64     `yaml:"stockLength,omitempty"`
	MassPerMeter float64      `yaml:"massPerMeter,omitempty"`
	B            float64      `yaml:"b,omitempty"`
	C            float64      `yaml:"c,omitempty"`
	S            float64      `yaml:"s,omitempty"`
	R            float64      `yaml:"r,omitempty"`
	A            float64      `yaml:"a,omitempty"`
	R1           float64      `yaml:"r1,omitempty"`
	R2           float64      `yaml:"r2,omitempty"`
	T            float64      `yaml:"t,omitempty"`
}

// Config is the validated effective configuration.
type Config struct {
	Version  int               `yaml:"version"`
	Defaults Defaults          `yaml:"defaults"`
	Remnants Remnants          `yaml:"remnants"`
	Families map[string]Family `yaml:"families"`
	Profiles map[string]Part   `yaml:"profiles"`

	source string
}

//go:embed defaults.yaml
var defaultsYAML []byte

// Load reads a v2 configuration file and merges it over the embedded defaults.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %q: %w", path, err)
	}

	cfg, err := Parse(data)
	if err != nil {
		return nil, fmt.Errorf("load config %q: %w", path, err)
	}
	cfg.source = path
	return cfg, nil
}

// Parse parses a v2 configuration document and merges it over the embedded
// defaults. Unknown fields are rejected.
func Parse(data []byte) (*Config, error) {
	merged, err := mergeDocuments(defaultsYAML, data)
	if err != nil {
		return nil, err
	}

	decoder := yaml.NewDecoder(bytes.NewReader(merged))
	decoder.KnownFields(true)

	var cfg Config
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidConfig, err)
	}

	var extra yaml.Node
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("%w: multiple yaml documents are not supported", ErrInvalidConfig)
	}

	if err := cfg.normalize(); err != nil {
		return nil, err
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Source returns the path the configuration was loaded from, if any.
func (c *Config) Source() string {
	return c.source
}

// MinRemnantLength returns the minimum useful remnant length.
func (c *Config) MinRemnantLength() float64 {
	return c.Remnants.MinLength
}

// UseRemnants reports whether remnants may be used.
func (c *Config) UseRemnants() bool {
	return c.Remnants.Enabled
}

// RemnantFile returns the configured remnant CSV path.
func (c *Config) RemnantFile() string {
	return c.Remnants.File
}

// Family returns the family configured for a profile type.
func (c *Config) Family(profileType string) (string, error) {
	parsed, err := domain.ParseProfileDesignation(profileType)
	if err != nil {
		return "", err
	}

	if profile, ok := c.Profiles[parsed.Key()]; ok && profile.Family != "" {
		return profile.Family, nil
	}

	var (
		family  string
		longest int
	)
	for name, familyConfig := range c.Families {
		for _, prefix := range familyConfig.Prefixes {
			if !strings.HasPrefix(parsed.Key(), prefix) {
				continue
			}
			if len(prefix) > longest {
				family = name
				longest = len(prefix)
			}
		}
	}
	if family == "" {
		return "", fmt.Errorf("%w: %s", ErrUnknownFamily, parsed.String())
	}
	return family, nil
}

// Spacing resolves the endcut spacing for a profile type and endcut.
func (c *Config) Spacing(profileType, endcut string) (float64, error) {
	parsed, err := domain.ParseProfileDesignation(profileType)
	if err != nil {
		return 0, err
	}
	family, err := c.Family(parsed.Key())
	if err != nil {
		return 0, err
	}

	key := strings.ToLower(strings.TrimSpace(endcut))
	if familyConfig, ok := c.Families[family]; ok {
		if endcutConfig, ok := familyConfig.Endcuts[key]; ok {
			if spacing, ok := endcutConfig.Profiles[parsed.Key()]; ok {
				return spacing, nil
			}
			if endcutConfig.Spacing != nil {
				return *endcutConfig.Spacing, nil
			}
		}
		if familyConfig.EndcutSpacing != nil {
			return *familyConfig.EndcutSpacing, nil
		}
	}
	if c.Defaults.EndcutSpacing >= 0 {
		return c.Defaults.EndcutSpacing, nil
	}
	return 0, fmt.Errorf("%w: family %q endcut %q profile %s", ErrNoSpacing, family, key, parsed.String())
}

// StockLength resolves the stock bar length for a profile type.
func (c *Config) StockLength(profileType string) (float64, error) {
	parsed, err := domain.ParseProfileDesignation(profileType)
	if err != nil {
		return 0, err
	}
	if profile, ok := c.Profiles[parsed.Key()]; ok && profile.StockLength != nil {
		return *profile.StockLength, nil
	}
	family, err := c.Family(parsed.Key())
	if err != nil {
		return 0, err
	}
	if familyConfig, ok := c.Families[family]; ok && familyConfig.StockLength != nil {
		return *familyConfig.StockLength, nil
	}
	if c.Defaults.StockLength > 0 {
		return c.Defaults.StockLength, nil
	}
	return 0, fmt.Errorf("%w: family %q profile %s", ErrNoStockLength, family, parsed.String())
}

// Spec resolves the cross-section geometry for a profile type.
func (c *Config) Spec(profileType string) (domain.SectionGeometry, error) {
	parsed, err := domain.ParseProfileDesignation(profileType)
	if err != nil {
		return domain.SectionGeometry{}, err
	}

	profile, ok := c.Profiles[parsed.Key()]
	if !ok {
		if parsed.Family() == "FB" {
			return domain.SectionGeometry{
				Shape: domain.ShapeFB,
				B:     parsed.Height(),
				S:     parsed.Thickness(),
			}, nil
		}
		return domain.SectionGeometry{}, fmt.Errorf("%w: %s", ErrNoSectionGeometry, parsed.String())
	}

	spec := domain.SectionGeometry{
		Shape:        profile.Shape,
		B:            profile.B,
		C:            profile.C,
		S:            profile.S,
		R:            profile.R,
		A:            profile.A,
		R1:           profile.R1,
		R2:           profile.R2,
		T:            profile.T,
		MassPerMeter: profile.MassPerMeter,
	}
	if err := spec.Validate(); err != nil {
		return domain.SectionGeometry{}, fmt.Errorf("%w: profile %s: %w", ErrInvalidConfig, parsed.String(), err)
	}
	return spec, nil
}

func (c *Config) normalize() error {
	c.Families = nonNilMap(c.Families)
	c.Profiles = nonNilMap(c.Profiles)

	families := make(map[string]Family, len(c.Families))
	for name, family := range c.Families {
		key := strings.ToLower(strings.TrimSpace(name))
		if key == "" {
			return fmt.Errorf("%w: empty family name", ErrInvalidConfig)
		}
		if _, exists := families[key]; exists {
			return fmt.Errorf("%w: duplicate family %q", ErrInvalidConfig, key)
		}

		prefixes := make([]string, 0, len(family.Prefixes))
		seen := make(map[string]struct{}, len(family.Prefixes))
		for _, prefix := range family.Prefixes {
			normalized := strings.ToLower(strings.TrimSpace(prefix))
			if normalized == "" {
				return fmt.Errorf("%w: family %q has an empty prefix", ErrInvalidConfig, key)
			}
			if _, exists := seen[normalized]; exists {
				continue
			}
			seen[normalized] = struct{}{}
			prefixes = append(prefixes, normalized)
		}
		family.Prefixes = prefixes

		endcuts := make(map[string]Endcut, len(family.Endcuts))
		for endcutName, endcut := range family.Endcuts {
			endcutKey := strings.ToLower(strings.TrimSpace(endcutName))
			if endcutKey == "" {
				return fmt.Errorf("%w: family %q has an empty endcut name", ErrInvalidConfig, key)
			}
			if _, exists := endcuts[endcutKey]; exists {
				return fmt.Errorf("%w: duplicate endcut %q in family %q", ErrInvalidConfig, endcutKey, key)
			}
			if endcut.Profiles != nil {
				profiles := make(map[string]float64, len(endcut.Profiles))
				for profileName, spacing := range endcut.Profiles {
					parsed, err := domain.ParseProfileDesignation(profileName)
					if err != nil {
						return fmt.Errorf("%w: family %q endcut %q: %w", ErrInvalidConfig, key, endcutKey, err)
					}
					profiles[parsed.Key()] = spacing
				}
				endcut.Profiles = profiles
			}
			endcuts[endcutKey] = endcut
		}
		family.Endcuts = endcuts
		families[key] = family
	}
	c.Families = families

	profiles := make(map[string]Part, len(c.Profiles))
	for name, profile := range c.Profiles {
		parsed, err := domain.ParseProfileDesignation(name)
		if err != nil {
			return fmt.Errorf("%w: profile %q: %w", ErrInvalidConfig, name, err)
		}
		if _, exists := profiles[parsed.Key()]; exists {
			return fmt.Errorf("%w: duplicate profile %q", ErrInvalidConfig, parsed.String())
		}
		profile.Family = strings.ToLower(strings.TrimSpace(profile.Family))
		profiles[parsed.Key()] = profile
	}
	c.Profiles = profiles
	return nil
}

func (c *Config) validate() error {
	if c.Version != Version {
		return fmt.Errorf("%w: version %d is not supported; expected %d", ErrInvalidConfig, c.Version, Version)
	}
	if c.Defaults.StockLength <= 0 {
		return fmt.Errorf("%w: defaults.stockLength must be positive", ErrInvalidConfig)
	}
	if c.Defaults.EndcutSpacing < 0 {
		return fmt.Errorf("%w: defaults.endcutSpacing must not be negative", ErrInvalidConfig)
	}
	if c.Remnants.MinLength < 0 {
		return fmt.Errorf("%w: remnants.minLength must not be negative", ErrInvalidConfig)
	}
	if c.Remnants.Enabled && strings.TrimSpace(c.Remnants.File) == "" {
		return fmt.Errorf("%w: remnants.file is required when remnants are enabled", ErrInvalidConfig)
	}
	if len(c.Families) == 0 {
		return fmt.Errorf("%w: at least one family is required", ErrInvalidConfig)
	}

	prefixes := make(map[string]string)
	for name, family := range c.Families {
		if len(family.Prefixes) == 0 {
			return fmt.Errorf("%w: family %q has no prefixes", ErrInvalidConfig, name)
		}
		for _, prefix := range family.Prefixes {
			if owner, exists := prefixes[prefix]; exists && owner != name {
				return fmt.Errorf("%w: prefix %q is shared by families %q and %q", ErrInvalidConfig, prefix, owner, name)
			}
			prefixes[prefix] = name
		}
		if family.StockLength != nil && *family.StockLength <= 0 {
			return fmt.Errorf("%w: family %q stockLength must be positive", ErrInvalidConfig, name)
		}
		if family.EndcutSpacing != nil && *family.EndcutSpacing < 0 {
			return fmt.Errorf("%w: family %q endcutSpacing must not be negative", ErrInvalidConfig, name)
		}
		for endcutName, endcut := range family.Endcuts {
			if endcut.Spacing != nil && *endcut.Spacing < 0 {
				return fmt.Errorf("%w: family %q endcut %q spacing must not be negative", ErrInvalidConfig, name, endcutName)
			}
			for profileName, spacing := range endcut.Profiles {
				if spacing < 0 {
					return fmt.Errorf("%w: family %q endcut %q profile %q spacing must not be negative", ErrInvalidConfig, name, endcutName, profileName)
				}
			}
		}
	}

	for name, profile := range c.Profiles {
		spec := domain.SectionGeometry{
			Shape:        profile.Shape,
			B:            profile.B,
			C:            profile.C,
			S:            profile.S,
			R:            profile.R,
			A:            profile.A,
			R1:           profile.R1,
			R2:           profile.R2,
			T:            profile.T,
			MassPerMeter: profile.MassPerMeter,
		}
		if err := spec.Validate(); err != nil {
			return fmt.Errorf("%w: profile %q: %w", ErrInvalidConfig, name, err)
		}
		if profile.StockLength != nil && *profile.StockLength <= 0 {
			return fmt.Errorf("%w: profile %q stockLength must be positive", ErrInvalidConfig, name)
		}
		if profile.Family != "" {
			if _, ok := c.Families[profile.Family]; !ok {
				return fmt.Errorf("%w: profile %q references unknown family %q", ErrInvalidConfig, name, profile.Family)
			}
		} else if _, err := c.Family(name); err != nil {
			return fmt.Errorf("%w: profile %q: %w", ErrInvalidConfig, name, err)
		}
	}

	return nil
}

func nonNilMap[T any](m map[string]T) map[string]T {
	if m == nil {
		return make(map[string]T)
	}
	return m
}
