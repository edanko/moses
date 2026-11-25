package main

import (
	"flag"
	"fmt"

	"github.com/edanko/moses/internal/app"
	"github.com/edanko/moses/internal/config"
)

type convFlags struct {
	inDir     string
	outDir    string
	remnants  string
	algorithm string
	project   string
}

func (c *cli) runConv(args []string) error {
	if len(args) == 0 {
		c.printConvUsage()
		return nil
	}

	switch args[0] {
	case "-h", "--help":
		c.printConvUsage()
		return nil
	case "10022":
		return c.runConvert("10022", "Cadmatic DXF profile reports", false, false, args[1:])
	case "10510":
		return c.runConvert("10510", "AVEVA profile sketch and list CSV/LST", true, true, args[1:])
	case "10510gen":
		return c.runConvert("10510gen", "AVEVA marine GEN profiles", true, true, args[1:])
	default:
		return fmt.Errorf("%w: unknown converter %q", errUsage, args[0])
	}
}

func (c *cli) printConvUsage() {
	_, _ = fmt.Fprint(c.out, `Convert production source files into nested TXT and reports.

Usage:
  moses conv <converter> [options]

Converters:
  10022      Cadmatic DXF profile reports
  10510      AVEVA profile sketch and list CSV/LST
  10510gen   AVEVA marine GEN profiles
`)
}

func (c *cli) runConvert(use, short string, projectFlag, excel bool, args []string) error {
	fs := flag.NewFlagSet(use, flag.ContinueOnError)
	fs.SetOutput(c.errOut)

	flags := &convFlags{}
	fs.StringVar(&flags.inDir, "in", "in", "input directory")
	fs.StringVar(&flags.outDir, "out", "out", "output directory")
	fs.StringVar(&flags.remnants, "remnants", "", "remnant CSV file (default from config)")
	fs.StringVar(&flags.algorithm, "algorithm", "best-fit-decreasing", "nesting algorithm")
	if projectFlag {
		fs.StringVar(&flags.project, "project", "", "project number (required)")
	}
	fs.Usage = func() {
		c.printCommandUsage(fs, "conv "+use, short)
	}

	help, err := c.parseFlags(fs, args)
	if err != nil {
		return err
	}
	if help {
		return nil
	}

	cfg, err := config.Load(c.configPath)
	if err != nil {
		return err
	}

	options := app.ConvertOptions{
		Config:       cfg,
		InDir:        flags.inDir,
		OutDir:       flags.outDir,
		RemnantsPath: flags.remnants,
		Algorithm:    flags.algorithm,
		Project:      flags.project,
		WriteExcel:   excel,
		Log:          c.errOut,
	}

	var summary app.ConvertSummary
	switch use {
	case "10022":
		summary, err = app.Convert10022(c.ctx, options)
	case "10510":
		summary, err = app.Convert10510(c.ctx, options)
	case "10510gen":
		summary, err = app.Convert10510Gen(c.ctx, options)
	default:
		return fmt.Errorf("unknown converter %q", use)
	}
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(
		c.out,
		"conv %s: %d files, %d parts, %d nests, %d remnants used, %d created, %.2f%% utilization\n",
		use,
		summary.Files,
		summary.Parts,
		summary.Nests,
		summary.RemnantsUsed,
		summary.RemnantsCreated,
		summary.Utilization,
	)
	return err
}
