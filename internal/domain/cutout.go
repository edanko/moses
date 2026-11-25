package domain

import (
	"fmt"
	"maps"
	"strconv"
	"strings"
)

// Cutout describes an internal cut (hole, notch or other icut). Raw preserves
// the source text so converter output can be reproduced exactly.
type Cutout struct {
	kind   string
	params map[string]float64
	raw    string
}

// ParseCutout parses a cut line such as "i1 x=180 r=30".
func ParseCutout(line string) (Cutout, error) {
	fields := strings.Fields(strings.TrimSpace(line))
	if len(fields) == 0 {
		return Cutout{}, fmt.Errorf("cut is empty")
	}

	params := make(map[string]float64, len(fields)-1)
	for _, field := range fields[1:] {
		key, value, ok := strings.Cut(field, "=")
		if !ok || key == "" {
			return Cutout{}, fmt.Errorf("cut %q: invalid parameter %q", line, field)
		}
		parsed, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return Cutout{}, fmt.Errorf("cut %q: invalid value for %q: %w", line, key, err)
		}
		params[strings.ToLower(key)] = parsed
	}

	return Cutout{
		kind:   fields[0],
		params: maps.Clone(params),
		raw:    strings.Join(fields, " "),
	}, nil
}

// NewCutout creates a cut from its kind and numeric parameters.
func NewCutout(kind string, params map[string]float64) *Cutout {
	return &Cutout{
		kind:   kind,
		params: maps.Clone(params),
	}
}

// Kind returns the cut type.
func (h Cutout) Kind() string {
	return h.kind
}

func (h Cutout) clone() Cutout {
	h.params = maps.Clone(h.params)
	return h
}

// Params returns a copy of the numeric parameters.
func (h Cutout) Params() map[string]float64 {
	return maps.Clone(h.params)
}

// Param returns one parameter value, or zero when absent.
func (h Cutout) Param(key string) float64 {
	return h.params[key]
}

// SetParam sets one numeric parameter.
func (h *Cutout) SetParam(key string, value float64) {
	if h.params == nil {
		h.params = make(map[string]float64)
	}
	h.params[strings.ToLower(key)] = value
	h.raw = ""
}

// String returns the source text when available, otherwise a rendering of the
// kind and parameters.
func (h Cutout) String() string {
	if h.raw != "" {
		return h.raw
	}
	return formatCut(h.kind, h.params)
}
