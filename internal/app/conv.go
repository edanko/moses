package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/edanko/moses/internal/config"
	"github.com/edanko/moses/internal/domain"
	"github.com/edanko/moses/internal/nesting"
	"github.com/edanko/moses/internal/remnant"
	"github.com/edanko/moses/internal/report"
	"github.com/edanko/moses/internal/source/aveva"
	"github.com/edanko/moses/internal/source/dxf"
	"github.com/edanko/moses/internal/source/gen"
)

// ConvertOptions configures a converter run.
type ConvertOptions struct {
	Config       *config.Config
	InDir        string
	OutDir       string
	RemnantsPath string
	Algorithm    string
	Project      string
	WriteExcel   bool
	Log          io.Writer
}

// ConvertSummary reports a converter run.
type ConvertSummary struct {
	Files           int
	Parts           int
	Nests           int
	RemnantsUsed    int
	RemnantsCreated int
	Utilization     float64
}

// Convert10022 converts Cadmatic DXF reports.
func Convert10022(ctx context.Context, options ConvertOptions) (ConvertSummary, error) {
	files, err := filepath.Glob(filepath.Join(options.InDir, "*", "*.dxf"))
	if err != nil {
		return ConvertSummary{}, err
	}
	sort.Strings(files)
	if len(files) == 0 {
		return ConvertSummary{}, fmt.Errorf("no DXF files found in %q", options.InDir)
	}

	parts := make([]domain.Part, 0)
	summary := ConvertSummary{}
	for _, path := range files {
		select {
		case <-ctx.Done():
			return summary, ctx.Err()
		default:
		}

		profile, err := dxf.ParseFile(path)
		if err != nil {
			return summary, err
		}
		summary.Files++
		if profile == nil {
			continue
		}
		parts = append(parts, *profile)
		logf(options.Log, "parsed %s\n", path)
	}
	summary.Parts = len(parts)
	if len(parts) == 0 {
		return summary, fmt.Errorf("no supported profiles in %q", options.InDir)
	}

	converted, err := convertParts(ctx, parts, options)
	if err != nil {
		return summary, err
	}
	summary.Nests = converted.Nests
	summary.RemnantsUsed = converted.RemnantsUsed
	summary.RemnantsCreated = converted.RemnantsCreated
	summary.Utilization = converted.Utilization
	return summary, nil
}

// Convert10510 converts AVEVA profile sketch/list CSV files.
func Convert10510(ctx context.Context, options ConvertOptions) (ConvertSummary, error) {
	if strings.TrimSpace(options.Project) == "" {
		return ConvertSummary{}, fmt.Errorf("project is required for 10510")
	}
	files, err := filepath.Glob(filepath.Join(options.InDir, "*.csv"))
	if err != nil {
		return ConvertSummary{}, err
	}
	sort.Strings(files)
	if len(files) == 0 {
		return ConvertSummary{}, fmt.Errorf("no CSV files found in %q", options.InDir)
	}

	parts := make([]domain.Part, 0)
	summary := ConvertSummary{}
	for _, path := range files {
		select {
		case <-ctx.Done():
			return summary, ctx.Err()
		default:
		}

		parsed, err := aveva.ParseFile(path, options.Project)
		if err != nil {
			return summary, err
		}
		summary.Files++
		parts = append(parts, parsed...)
		logf(options.Log, "parsed %s\n", path)
	}
	summary.Parts = len(parts)

	converted, err := convertParts(ctx, parts, options)
	if err != nil {
		return summary, err
	}
	summary.Nests = converted.Nests
	summary.RemnantsUsed = converted.RemnantsUsed
	summary.RemnantsCreated = converted.RemnantsCreated
	summary.Utilization = converted.Utilization
	return summary, nil
}

// Convert10510Gen converts AVEVA marine GEN files.
func Convert10510Gen(ctx context.Context, options ConvertOptions) (ConvertSummary, error) {
	if strings.TrimSpace(options.Project) == "" {
		return ConvertSummary{}, fmt.Errorf("project is required for 10510gen")
	}

	var files []string
	err := filepath.WalkDir(options.InDir, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		if strings.EqualFold(filepath.Ext(path), ".gen") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return ConvertSummary{}, fmt.Errorf("scan %q: %w", options.InDir, err)
	}
	sort.Strings(files)
	if len(files) == 0 {
		return ConvertSummary{}, fmt.Errorf("no GEN files found in %q", options.InDir)
	}

	parts := make([]domain.Part, 0)
	summary := ConvertSummary{}
	for _, path := range files {
		select {
		case <-ctx.Done():
			return summary, ctx.Err()
		default:
		}

		parsed, err := gen.ParseFile(path, options.Project)
		if err != nil {
			return summary, err
		}
		summary.Files++
		parts = append(parts, parsed...)
		logf(options.Log, "parsed %s\n", path)
	}
	summary.Parts = len(parts)

	converted, err := convertParts(ctx, parts, options)
	if err != nil {
		return summary, err
	}
	summary.Nests = converted.Nests
	summary.RemnantsUsed = converted.RemnantsUsed
	summary.RemnantsCreated = converted.RemnantsCreated
	summary.Utilization = converted.Utilization
	return summary, nil
}

func convertParts(ctx context.Context, parts []domain.Part, options ConvertOptions) (ConvertSummary, error) {
	var summary ConvertSummary
	if options.Config == nil {
		return summary, fmt.Errorf("conv: config is required")
	}

	store, err := loadRemnants(options)
	if err != nil {
		return summary, err
	}
	remnants := []domain.Remnant{}
	if store != nil {
		remnants = store.Items()
	}

	packed, err := nesting.Pack(ctx, parts, remnants, nesting.PackOptions{
		Spacing:     options.Config.Spacing,
		StockLength: options.Config.StockLength,
		Algorithm:   options.Algorithm,
	})
	if err != nil {
		return summary, err
	}
	summary.Utilization = packed.Utilization

	bars, err := buildBars(packed.Nests, store, options.Config.MinRemnantLength())
	if err != nil {
		return summary, err
	}
	summary.Nests = len(bars)
	for _, bar := range bars {
		if bar.Consumed != nil {
			summary.RemnantsUsed++
		}
		if bar.Created != nil {
			summary.RemnantsCreated++
		}
	}

	if err := writeConversionOutput(options, parts, bars); err != nil {
		return summary, err
	}
	if store != nil {
		if err := store.Save(); err != nil {
			return summary, err
		}
	}
	return summary, nil
}

func buildBars(nests []domain.Nest, store *remnant.Store, minRemnant float64) ([]report.Bar, error) {
	type group struct {
		bars []report.Bar
	}
	groups := make(map[string]*group)
	for _, nest := range nests {
		stock := nest.StockBar()
		key := stock.DesignationName() + "\x00" + stock.Grade()
		entry, ok := groups[key]
		if !ok {
			entry = &group{}
			groups[key] = entry
		}
		index := len(entry.bars) + 1
		name := nest.ShortName(index)
		bar := report.Bar{
			Nest:     nest,
			Index:    index,
			Name:     name,
			Consumed: nest.ConsumedRemnant(),
		}
		if bar.Consumed != nil && store != nil {
			store.MarkUsed(bar.Consumed.ID, name)
		}
		created := nest.CreateRemnant(name, minRemnant)
		if created != nil {
			if profiles := nest.Parts(); len(profiles) > 0 {
				created.Source = profiles[0].Section()
			}
			if store != nil {
				stored := store.Add(*created)
				bar.Created = stored
			} else {
				bar.Created = created
			}
		}
		entry.bars = append(entry.bars, bar)
	}

	keys := make([]string, 0, len(groups))
	for key := range groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	bars := make([]report.Bar, 0, len(nests))
	for _, key := range keys {
		bars = append(bars, groups[key].bars...)
	}
	return bars, nil
}

func loadRemnants(options ConvertOptions) (*remnant.Store, error) {
	if !options.Config.UseRemnants() {
		return nil, nil
	}
	path := options.RemnantsPath
	if path == "" {
		path = options.Config.RemnantFile()
	}
	return remnant.Load(path)
}

func writeConversionOutput(options ConvertOptions, parts []domain.Part, bars []report.Bar) error {
	withSection := sections(parts) > 1

	textGroups := make(map[string][]report.Bar)
	textOrder := make([]string, 0)
	for _, bar := range bars {
		stock := bar.Nest.StockBar()
		key := stock.DesignationName() + "\x00" + stock.Grade()
		if _, ok := textGroups[key]; !ok {
			textOrder = append(textOrder, key)
		}
		textGroups[key] = append(textGroups[key], bar)
	}
	for _, key := range textOrder {
		group := textGroups[key]
		name, err := report.TxtFileName(group[0])
		if err != nil {
			return err
		}
		path := filepath.Join(options.OutDir, "txt", name+".txt")
		if err := writeFile(path, report.Txt(group, withSection)); err != nil {
			return err
		}
		logf(options.Log, "wrote %s\n", path)
	}

	base := reportBase(parts)
	barPath := filepath.Join(options.OutDir, base+"_bar.txt")
	if err := writeFile(barPath, reportBarList(options, bars)); err != nil {
		return err
	}
	nstPath := filepath.Join(options.OutDir, base+"_nst.txt")
	nestingList, err := reportNestingList(options, bars)
	if err != nil {
		return err
	}
	if err := writeFile(nstPath, nestingList); err != nil {
		return err
	}

	if options.WriteExcel {
		reporter := newReporter(options)
		excelPath := filepath.Join(options.OutDir, base+".xlsx")
		if err := reporter.Excel(excelPath, bars); err != nil {
			return err
		}
		logf(options.Log, "wrote %s\n", excelPath)
	}
	return nil
}

func reportBase(parts []domain.Part) string {
	project := ""
	section := ""
	if len(parts) > 0 {
		project = parts[0].Project()
		sections := make(map[string]struct{})
		for _, part := range parts {
			sections[part.Section()] = struct{}{}
		}
		if len(sections) == 1 {
			section = parts[0].Section()
		} else {
			section = "xxx"
		}
	}
	if project == "" {
		return section
	}
	if section == "" {
		return project
	}
	return project + "-" + section
}

func sections(parts []domain.Part) int {
	unique := make(map[string]struct{})
	for _, part := range parts {
		unique[part.Section()] = struct{}{}
	}
	return len(unique)
}

func newReporter(options ConvertOptions) *report.Reporter {
	return &report.Reporter{
		Spacing:          options.Config.Spacing,
		Spec:             options.Config.Spec,
		MinRemnantLength: options.Config.MinRemnantLength(),
	}
}

func reportBarList(options ConvertOptions, bars []report.Bar) string {
	return newReporter(options).BarList(bars)
}

func reportNestingList(options ConvertOptions, bars []report.Bar) (string, error) {
	return newReporter(options).NestingList(bars)
}
