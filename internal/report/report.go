// Package report renders operator-facing bar and nesting reports.
package report

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/edanko/moses/internal/domain"
)

// Bar couples a nest with the one-based number and name used in a converter
// run.
type Bar struct {
	Nest     domain.Nest
	Index    int
	Name     string
	Consumed *domain.Remnant
	Created  *domain.Remnant
}

// SpacingFunc resolves endcut spacing.
type SpacingFunc func(profileType, endcut string) (float64, error)

// SpecFunc resolves profile geometry.
type SpecFunc func(profileType string) (domain.SectionGeometry, error)

// Reporter renders reports for one converter run.
type Reporter struct {
	Spacing          SpacingFunc
	Spec             SpecFunc
	MinRemnantLength float64
}

// BarList renders the material consumption summary.
func (r *Reporter) BarList(bars []Bar) string {
	type group struct {
		dim     string
		quality string
		stock   []Bar
		remnant []Bar
	}

	groups := make(map[string]*group)
	for _, bar := range bars {
		stock := bar.Nest.StockBar()
		dim := stock.DesignationName()
		quality := strings.ToUpper(stock.Grade())
		key := dim + "\x00" + quality
		entry, ok := groups[key]
		if !ok {
			entry = &group{dim: dim, quality: quality}
			groups[key] = entry
		}
		if stock.Type() == domain.StockTypeRemnant || bar.Consumed != nil {
			entry.remnant = append(entry.remnant, bar)
		} else {
			entry.stock = append(entry.stock, bar)
		}
	}

	keys := make([]string, 0, len(groups))
	for key := range groups {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		left, right := groups[keys[i]], groups[keys[j]]
		if left.dim != right.dim {
			return left.dim < right.dim
		}
		return left.quality < right.quality
	})

	var builder strings.Builder
	for _, key := range keys {
		entry := groups[key]
		fmt.Fprintf(&builder, "%10s", strings.ToUpper(entry.dim))
		fmt.Fprintf(&builder, " |%6s |", entry.quality)

		if len(entry.stock) > 0 {
			stock := entry.stock[0].Nest.StockBar()
			fmt.Fprintf(&builder, "%3d bars (%5f)", len(entry.stock), stock.Length())
		}
		if len(entry.stock) > 0 && len(entry.remnant) > 0 {
			builder.WriteString(" and ")
		}
		if len(entry.remnant) > 0 {
			stock := entry.remnant[0].Nest.StockBar()
			fmt.Fprintf(&builder, "%3d remnants (%s", len(entry.remnant), formatFloat(stock.Length()))
			for _, bar := range entry.remnant[1:] {
				stock := bar.Nest.StockBar()
				builder.WriteString(" + ")
				builder.WriteString(formatFloat(stock.Length()))
			}
			builder.WriteString(")")
		}
		builder.WriteString("\r\n")
	}
	return builder.String()
}

// NestingList renders the per-bar cutting cards.
func (r *Reporter) NestingList(bars []Bar) (string, error) {
	var builder strings.Builder
	for _, bar := range bars {
		scrap, err := r.scrap(bar)
		if err != nil {
			return "", fmt.Errorf("nest %q: %w", bar.Name, err)
		}

		stock := bar.Nest.StockBar()
		builder.WriteString(strings.ToUpper(stock.DesignationName()))
		builder.WriteString(" ")
		builder.WriteString(strings.ToUpper(stock.Grade()))
		builder.WriteString("\r\n")
		builder.WriteString("Имя программы      ")
		builder.WriteString(bar.Name)
		builder.WriteString("\r\n")

		if bar.Consumed != nil {
			builder.WriteString("Использовать отход ")
			builder.WriteString(bar.Consumed.Marking)
			builder.WriteString("\r\n")
		}

		builder.WriteString("Длина: ")
		builder.WriteString(formatFloat(stock.Length()))
		builder.WriteString(" / ")
		builder.WriteString("Исп.: ")
		builder.WriteString(formatFloat(bar.Nest.TotalPartLength()))
		usable := stock.Length() - bar.Nest.RemnantLength(r.MinRemnantLength)
		fmt.Fprintf(&builder, " (%.2f %%)", percent(bar.Nest.TotalPartLength(), usable))

		if bar.Created != nil {
			builder.WriteString(" / Отход: ")
			builder.WriteString(formatFloat(bar.Created.Length))
		}

		builder.WriteString(" / Лом: ")
		builder.WriteString(formatFloat(scrap))
		builder.WriteString("\r\n")

		if bar.Created != nil {
			builder.WriteString("Маркировать отход  ")
			builder.WriteString(bar.Created.Marking)
			builder.WriteString("\r\n")
		}

		builder.WriteString("----------------------------\r\n")
		builder.WriteString("| Секция | Позиция | Длина |\r\n")
		builder.WriteString("----------------------------\r\n")
		for _, profile := range bar.Nest.Parts() {
			fmt.Fprintf(
				&builder,
				"| %6s | %7s | %5s |\r\n",
				profile.Section(),
				profile.PositionNo(),
				formatFloat(profile.Length()),
			)
		}
		builder.WriteString("----------------------------\r\n")
		builder.WriteString("\r\n")
	}
	return builder.String(), nil
}

// Txt renders the nested TXT format consumed by moses parser.
func Txt(bars []Bar, withSection bool) string {
	var builder strings.Builder
	for i, bar := range bars {
		if i != 0 {
			builder.WriteString(".\r\n")
		}
		stock := bar.Nest.StockBar()
		builder.WriteString(formatFloat(stock.Length()))
		builder.WriteString("\r\n")

		for _, part := range bar.Nest.Parts() {
			builder.WriteString("\r\n")
			if withSection {
				builder.WriteString(part.Section())
				builder.WriteString("-")
			}
			builder.WriteString(part.PositionNo())
			builder.WriteString("\r\n")
			builder.WriteString(formatFloat(part.Length()))
			builder.WriteString("\r\n")
			builder.WriteString(part.LeftEndCut().String())
			builder.WriteString("\r\n")
			builder.WriteString(part.RightEndCut().String())
			builder.WriteString("\r\n")
			for _, cutout := range part.Cutouts() {
				builder.WriteString(cutout.String())
				builder.WriteString("\r\n")
			}
		}
	}
	return builder.String()
}

// TxtFileName returns the production TXT name for a bar.
func TxtFileName(bar Bar) (string, error) {
	parts := bar.Nest.Parts()
	if len(parts) == 0 {
		return "", fmt.Errorf("bar %q has no parts", bar.Name)
	}
	part := parts[0]
	stock := bar.Nest.StockBar()
	var builder strings.Builder
	if part.Project() != "" {
		builder.WriteString(part.Project())
		builder.WriteString("-")
	}
	builder.WriteString(part.Section())
	builder.WriteString("_")
	builder.WriteString(strings.ToUpper(stock.DesignationName()))
	if grade := stock.Grade(); grade != "" {
		builder.WriteString("_")
		builder.WriteString(strings.ToUpper(grade))
	}
	return builder.String(), nil
}

func (r *Reporter) usedLength(bar Bar) (float64, error) {
	nest := bar.Nest
	stock := nest.StockBar()
	designationKey := stock.DesignationKey()
	var total float64
	for _, part := range nest.Parts() {
		left, err := r.Spacing(designationKey, part.LeftEndKind())
		if err != nil {
			return 0, err
		}
		right, err := r.Spacing(designationKey, part.RightEndKind())
		if err != nil {
			return 0, err
		}
		total += left + part.Length() + right
	}
	return total, nil
}

func (r *Reporter) scrap(bar Bar) (float64, error) {
	nest := bar.Nest
	stock := nest.StockBar()
	designationKey := stock.DesignationKey()
	var total float64
	for _, part := range nest.Parts() {
		left, err := r.Spacing(designationKey, part.LeftEndKind())
		if err != nil {
			return 0, err
		}
		right, err := r.Spacing(designationKey, part.RightEndKind())
		if err != nil {
			return 0, err
		}
		total += left + right
	}
	if !nest.HasUsefulRemnant(r.MinRemnantLength) {
		total += stock.Length() - nest.TotalPartLength()
	}
	return total, nil
}

func percent(part, used float64) float64 {
	if used == 0 {
		return 0
	}
	return part / used * 100
}

func formatFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func mass(profileType string, length float64, spec SpecFunc) (float64, error) {
	geometry, err := spec(profileType)
	if err != nil {
		return 0, err
	}
	return geometry.MassPerMeter * length / 1000, nil
}
