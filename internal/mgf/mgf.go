package mgf

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/edanko/moses/internal/domain"
)

var newLine = "\r\n"

// SpacingResolver resolves the endcut spacing for a profile type.
type SpacingResolver interface {
	Spacing(profileType, endcut string) (float64, error)
}

// SpecResolver resolves the cross-section geometry for a profile type.
type SpecResolver interface {
	Spec(profileType string) (domain.SectionGeometry, error)
}

// Formatter writes IMG MGF output.
type Formatter struct {
	spacing SpacingResolver
	spec    SpecResolver
}

// New creates an MGF formatter.
func New(spacing SpacingResolver, spec SpecResolver) *Formatter {
	return &Formatter{spacing: spacing, spec: spec}
}

// Params renders an end cut block. It is exported for tests.
func Params(profileType string, end domain.EndCut) (string, string) {
	return params(strings.ToLower(profileType), canonical(end.Kind(), end.Params()))
}

// IcutString renders one internal cut. It is exported for tests.
func IcutString(hole domain.Cutout) string {
	return icutString(canonical(hole.Kind(), hole.Params()))
}

func canonical(kind string, params map[string]float64) string {
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
		builder.WriteString(" ")
		builder.WriteString(key)
		builder.WriteString("=")
		builder.WriteString(strconv.FormatFloat(params[key], 'f', -1, 64))
	}
	return builder.String()
}

// StraightCut returns the part numbers that need a straight cut after them.
func (f *Formatter) StraightCut(nest domain.Nest) string {
	parts := nest.Parts()
	var builder strings.Builder

	for i, part := range parts {
		switch part.RightEndKind() {
		case "fv", "fv.1", "b-102", "b-103", "b-104", "b-106",
			"sb1", "sb2", "sb6", "sb9", "sb9-01":
			if len(parts) != i+1 {
				builder.WriteString(strconv.Itoa(i + 1))
			}
		}
	}
	return builder.String()
}

// Format renders one nest. index is the one-based bar number used in the nest
// name.
func (f *Formatter) Format(nest domain.Nest, index int) (output string, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("mgf: %v", recovered)
		}
	}()

	stock := nest.StockBar()
	profileType := stock.Designation()
	spec, err := f.spec.Spec(profileType.Key())
	if err != nil {
		return "", fmt.Errorf("profile %s: %w", profileType.String(), err)
	}
	norm, err := normText(profileType, spec)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	builder.WriteString("TYPE_OF_GENERIC_FILE=MOSES_NESTED_PROFILE")
	builder.WriteString(newLine)
	builder.WriteString("VERSION=1.0")
	builder.WriteString(newLine)
	builder.WriteString("USAGE=PLASMA")
	builder.WriteString(newLine)
	builder.WriteString("COMMON_DATA")
	builder.WriteString(newLine)
	builder.WriteString("NEST_NAME=")
	builder.WriteString(nest.Name(index))
	builder.WriteString(newLine)
	builder.WriteString(norm)
	builder.WriteString(newLine)

	builder.WriteString("RAW_LENGTH=")
	if stock.Length() > 1800 {
		builder.WriteString(strconv.FormatFloat(stock.Length(), 'f', -1, 64))
	} else {
		builder.WriteString("1800")
	}
	builder.WriteString(newLine)

	used, err := f.usedLength(nest)
	if err != nil {
		return "", err
	}
	builder.WriteString("USED_LENGTH=")
	builder.WriteString(strconv.FormatFloat(used, 'f', -1, 64))
	builder.WriteString(newLine)
	builder.WriteString("REST_LENGTH=")
	builder.WriteString(strconv.FormatFloat(stock.Length()-used, 'f', -1, 64))
	builder.WriteString(newLine)

	builder.WriteString("TEXT_HEIGHT=20")
	builder.WriteString(newLine)
	builder.WriteString("TEXT_WIDTH=16")
	builder.WriteString(newLine)
	builder.WriteString("TEXT_PLANE=0")
	builder.WriteString(newLine)
	builder.WriteString("TEXT_PLACING=3")
	builder.WriteString(newLine)
	builder.WriteString("TEXT_U=100")
	builder.WriteString(newLine)
	builder.WriteString("TEXT_V=80")
	builder.WriteString(newLine)
	builder.WriteString("MATERIAL=St37-2")
	builder.WriteString(newLine)
	builder.WriteString("NO_OF_PROFS=")
	builder.WriteString(strconv.Itoa(len(nest.Parts())))
	builder.WriteString(newLine)
	builder.WriteString("END_OF_COMMON_DATA")
	builder.WriteString(newLine)

	var leftPoint float64
	for i, profile := range nest.Parts() {
		builder.WriteString("PROFILE_DATA")
		builder.WriteString(newLine)

		builder.WriteString("TLENGTH=")
		builder.WriteString(strconv.FormatFloat(profile.Length(), 'f', -1, 64))
		builder.WriteString(newLine)

		builder.WriteString("IDENT_STRING=")
		builder.WriteString(stock.Project())
		builder.WriteString("-")
		builder.WriteString(profile.Section())
		builder.WriteString("-")
		builder.WriteString(profile.PositionNo())
		builder.WriteString(newLine)

		builder.WriteString("NO_OF_MARKS=0")
		builder.WriteString(newLine)
		builder.WriteString("NO_OF_ICUTS=")
		builder.WriteString(strconv.Itoa(len(profile.Cutouts())))
		builder.WriteString(newLine)
		builder.WriteString("NO_OF_PARTS=")
		builder.WriteString(strconv.Itoa(i + 1))
		builder.WriteString(newLine)
		builder.WriteString("END_OF_PROFILE_DATA")
		builder.WriteString(newLine)

		builder.WriteString("LEFT_END")
		builder.WriteString(newLine)

		leftSpacing, err := f.spacing.Spacing(profileType.Key(), profile.LeftEndKind())
		if err != nil {
			return "", fmt.Errorf("left end %q: %w", profile.LeftEndKind(), err)
		}
		leftPoint += leftSpacing
		builder.WriteString("LEFT_CLOSEST_POINT=")
		builder.WriteString(strconv.FormatFloat(leftPoint, 'f', -1, 64))
		builder.WriteString(newLine)

		left, leftScallop := Params(profileType.Key(), profile.LeftEndCut())
		builder.WriteString(left)
		builder.WriteString("END_OF_LEFT_END")
		builder.WriteString(newLine)

		leftPoint += profile.Length()
		builder.WriteString("RIGHT_END")
		builder.WriteString(newLine)
		builder.WriteString("LEFT_FARTHEST_POINT=")
		builder.WriteString(strconv.FormatFloat(leftPoint, 'f', -1, 64))
		builder.WriteString(newLine)

		right, rightScallop := Params(profileType.Key(), profile.RightEndCut())
		builder.WriteString(right)
		builder.WriteString("END_OF_RIGHT_END")
		builder.WriteString(newLine)

		rightSpacing, err := f.spacing.Spacing(profileType.Key(), profile.RightEndKind())
		if err != nil {
			return "", fmt.Errorf("right end %q: %w", profile.RightEndKind(), err)
		}
		leftPoint += rightSpacing

		if leftScallop != "" {
			builder.WriteString("LEFT_SCALLOP")
			builder.WriteString(newLine)
			builder.WriteString(leftScallop)
			builder.WriteString("END_OF_LEFT_SCALLOP")
			builder.WriteString(newLine)
		}
		if rightScallop != "" {
			builder.WriteString("RIGHT_SCALLOP")
			builder.WriteString(newLine)
			builder.WriteString(rightScallop)
			builder.WriteString("END_OF_RIGHT_SCALLOP")
			builder.WriteString(newLine)
		}

		for _, hole := range profile.Cutouts() {
			builder.WriteString("START_OF_ICUT")
			builder.WriteString(newLine)
			builder.WriteString(IcutString(hole))
			builder.WriteString("END_OF_ICUT")
			builder.WriteString(newLine)
		}
	}

	return builder.String(), nil
}

func (f *Formatter) usedLength(nest domain.Nest) (float64, error) {
	stock := nest.StockBar()
	designationKey := stock.DesignationKey()
	var used float64
	for _, part := range nest.Parts() {
		left, err := f.spacing.Spacing(designationKey, part.LeftEndKind())
		if err != nil {
			return 0, err
		}
		right, err := f.spacing.Spacing(designationKey, part.RightEndKind())
		if err != nil {
			return 0, err
		}
		used += left + part.Length() + right
	}
	return used, nil
}

func normText(profileType domain.ProfileDesignation, spec domain.SectionGeometry) (string, error) {
	switch spec.Shape {
	case domain.ShapeHP:
		return fmt.Sprintf(
			"SHAPE=HP\r\nSTART_OF_PARAMS\r\nNO_OF_PARAMS=5\r\nNORM=%s\r\nB=%0.1f\r\nC=%0.1f\r\nS=%0.1f\r\nR=%0.1f\r\nEND_OF_PARAMS",
			profileType.String(),
			spec.B,
			spec.C,
			spec.S,
			spec.R,
		), nil
	case domain.ShapeFB:
		return fmt.Sprintf(
			"SHAPE=FB\r\nSTART_OF_PARAMS\r\nNO_OF_PARAMS=3\r\nNORM=%s\r\nB=%s\r\nS=%s\r\nEND_OF_PARAMS",
			profileType.String(),
			strconv.FormatFloat(spec.B, 'f', -1, 64),
			strconv.FormatFloat(spec.S, 'f', -1, 64),
		), nil
	case domain.ShapeL:
		return fmt.Sprintf(
			"SHAPE=L\r\nSTART_OF_PARAMS\r\nNO_OF_PARAMS=6\r\nNORM=%s\r\nA=%0.1f\r\nB=%0.1f\r\nR1=%0.1f\r\nR2=%0.1f\r\nS=%0.1f\r\nT=%0.1f\r\nEND_OF_PARAMS",
			"L"+strings.ToUpper(profileType.String()[2:]),
			spec.A,
			spec.B,
			spec.R1,
			spec.R2,
			spec.S,
			spec.T,
		), nil
	default:
		return "", fmt.Errorf("unsupported profile shape %q", spec.Shape)
	}
}
