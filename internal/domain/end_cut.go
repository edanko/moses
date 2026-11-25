package domain

import (
	"fmt"
	"maps"
	"sort"
	"strconv"
	"strings"
)

// EndCut describes a profile end cut: a named cut type plus its numeric
// parameters. Raw preserves the source text so converter output can be
// reproduced exactly.
type EndCut struct {
	kind   string
	params map[string]float64
	raw    string
}

// ParseEndCut parses an end cut line such as "21 v1=0 v2=50 snipe=10".
func ParseEndCut(line string) (EndCut, error) {
	fields := strings.Fields(strings.TrimSpace(line))
	if len(fields) == 0 {
		return EndCut{}, fmt.Errorf("end cut is empty")
	}

	params := make(map[string]float64, len(fields)-1)
	for _, field := range fields[1:] {
		key, value, ok := strings.Cut(field, "=")
		if !ok || key == "" {
			return EndCut{}, fmt.Errorf("end cut %q: invalid parameter %q", line, field)
		}
		parsed, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return EndCut{}, fmt.Errorf("end cut %q: invalid value for %q: %w", line, key, err)
		}
		params[strings.ToLower(key)] = parsed
	}

	return EndCut{
		kind:   fields[0],
		params: maps.Clone(params),
		raw:    strings.Join(fields, " "),
	}, nil
}

// NewEndCut creates an end cut from its kind and numeric parameters.
func NewEndCut(kind string, params map[string]float64) EndCut {
	return EndCut{
		kind:   kind,
		params: maps.Clone(params),
	}
}

// Kind returns the end cut type.
func (e EndCut) Kind() string {
	return e.kind
}

// Params returns a copy of the numeric parameters.
func (e EndCut) Params() map[string]float64 {
	return maps.Clone(e.params)
}

// String returns the source text when available, otherwise a deterministic
// rendering of the kind and parameters.
func (e EndCut) String() string {
	if e.raw != "" {
		return e.raw
	}
	return formatCut(e.kind, e.params)
}

// Param returns one parameter value, or zero when absent.
func (e EndCut) Param(key string) float64 {
	return e.params[key]
}

func formatCut(kind string, params map[string]float64) string {
	if kind == "" {
		return ""
	}
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var builder strings.Builder
	builder.WriteString(kind)
	for _, key := range keys {
		builder.WriteByte(' ')
		builder.WriteString(key)
		builder.WriteByte('=')
		builder.WriteString(strconv.FormatFloat(params[key], 'f', -1, 64))
	}
	return builder.String()
}
