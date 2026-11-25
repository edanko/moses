package domain

import (
	"fmt"
	"strconv"
	"strings"
)

// ProfileDesignation is a normalized steel profile designation such as RP160X10,
// FB80x8 or AE250x90x10x15.
type ProfileDesignation struct {
	code      string
	family    string
	height    float64
	thickness float64
}

// ParseProfileDesignation parses the profile designations used by the source formats.
// It accepts RP/PP, FB and AE families, "*" or "x" separators and optional
// trailing ".0" components.
func ParseProfileDesignation(raw string) (ProfileDesignation, error) {
	s := strings.ToUpper(strings.TrimSpace(raw))
	s = strings.ReplaceAll(s, "*", "X")
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, ",", ".")
	if s == "" {
		return ProfileDesignation{}, fmt.Errorf("profile designation is empty")
	}
	if alias, ok := profileAliases[strings.ToLower(s)]; ok {
		s = alias
	}

	prefixEnd := 0
	for prefixEnd < len(s) && (s[prefixEnd] < '0' || s[prefixEnd] > '9') {
		prefixEnd++
	}
	if prefixEnd == 0 || prefixEnd == len(s) {
		return ProfileDesignation{}, fmt.Errorf("profile designation %q has no dimensions", raw)
	}

	family := s[:prefixEnd]
	parts := strings.Split(s[prefixEnd:], "X")
	values := make([]float64, len(parts))
	for i, part := range parts {
		value, err := strconv.ParseFloat(part, 64)
		if err != nil {
			return ProfileDesignation{}, fmt.Errorf("profile designation %q: invalid dimension %q", raw, part)
		}
		values[i] = value
	}

	switch family {
	case "PP":
		family = "RP"
	case "RP":
	case "FB":
	case "AE":
	default:
		return ProfileDesignation{}, fmt.Errorf("profile designation %q: unsupported family %q", raw, family)
	}

	var height, thickness float64
	switch family {
	case "RP", "FB":
		if len(values) != 2 {
			return ProfileDesignation{}, fmt.Errorf("profile designation %q: expected 2 dimensions", raw)
		}
		height, thickness = values[0], values[1]
	case "AE":
		switch len(values) {
		case 2:
			height, thickness = values[0], values[1]
		case 4:
			height, thickness = values[0], values[2]
		default:
			return ProfileDesignation{}, fmt.Errorf("profile designation %q: expected 2 or 4 dimensions", raw)
		}
	}

	code := family + joinDimensions(values)
	return ProfileDesignation{
		code:      code,
		family:    family,
		height:    height,
		thickness: thickness,
	}, nil
}

func joinDimensions(values []float64) string {
	parts := make([]string, len(values))
	for i, value := range values {
		parts[i] = strconv.FormatFloat(value, 'f', -1, 64)
	}
	return strings.Join(parts, "X")
}

// Key returns the lowercase designation used as a configuration key.
func (p ProfileDesignation) Key() string {
	return strings.ToLower(p.code)
}

// Family returns the profile family prefix: RP, FB or AE.
func (p ProfileDesignation) Family() string {
	return p.family
}

// Height returns the profile height (the first dimension).
func (p ProfileDesignation) Height() float64 {
	return p.height
}

// Thickness returns the profile wall thickness.
func (p ProfileDesignation) Thickness() float64 {
	return p.thickness
}

// String implements fmt.Stringer.
func (p ProfileDesignation) String() string {
	return p.code
}

// IsValid reports whether the type was populated.
func (p ProfileDesignation) IsValid() bool {
	return p.code != ""
}

// Equal reports whether two designations describe the same profile.
func (p ProfileDesignation) Equal(other ProfileDesignation) bool {
	return p.code == other.code
}

// profileAliases maps the short dimension codes used in legacy filenames and
// short configuration keys to their full profile designations.
var profileAliases = map[string]string{
	"10":  "RP100X6",
	"12":  "RP120X6.5",
	"14a": "RP140X7",
	"16a": "RP160X8",
	"16b": "RP160X10",
	"18a": "RP180X9",
	"18b": "RP180X11",
	"20a": "RP200X10",
	"20b": "RP200X11",
	"22a": "RP220X11",
	"22b": "RP220X13",
	"24a": "RP240X12",
	"24b": "RP240X14",
}
