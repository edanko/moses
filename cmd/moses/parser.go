package main

import (
	"flag"
	"fmt"

	"github.com/edanko/moses/internal/app"
	"github.com/edanko/moses/internal/config"
)

func (c *cli) runParser(args []string) error {
	fs := flag.NewFlagSet("parser", flag.ContinueOnError)
	fs.SetOutput(c.errOut)

	var (
		inDir      string
		outDir     string
		dontRename bool
	)
	fs.StringVar(&inDir, "in", "in", "input directory")
	fs.StringVar(&outDir, "out", "out", "output directory")
	fs.BoolVar(&dontRename, "dont-rename", false, "do not prefix processed input files with _")
	fs.Usage = func() {
		c.printCommandUsage(fs, "parser", "Convert nested TXT files to MGF machine output")
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

	summary, err := app.RunParser(c.ctx, app.ParserOptions{
		Config:     cfg,
		InDir:      inDir,
		OutDir:     outDir,
		DontRename: dontRename,
		Log:        c.errOut,
	})
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(
		c.out,
		"parsed %d files, %d nests, wrote %d MGF files (%d straight cuts), renamed %d inputs\n",
		summary.Files,
		summary.Nests,
		summary.MGFs,
		summary.StraightCuts,
		summary.Renamed,
	)
	return err
}
