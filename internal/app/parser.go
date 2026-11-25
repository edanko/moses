// Package app contains the file-pipeline use cases.
package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/edanko/moses/internal/config"
	"github.com/edanko/moses/internal/domain"
	"github.com/edanko/moses/internal/mgf"
	"github.com/edanko/moses/internal/source/txt"
)

// ParserOptions configures the TXT to MGF pipeline.
type ParserOptions struct {
	Config     *config.Config
	InDir      string
	OutDir     string
	DontRename bool
	Log        io.Writer
}

// ParserSummary reports what the parser produced.
type ParserSummary struct {
	Files        int
	Nests        int
	MGFs         int
	StraightCuts int
	Renamed      int
}

// RunParser parses every TXT file under InDir and writes MGF output.
func RunParser(ctx context.Context, options ParserOptions) (ParserSummary, error) {
	var summary ParserSummary
	if options.Config == nil {
		return summary, fmt.Errorf("parser: config is required")
	}

	files, err := txtFiles(options.InDir)
	if err != nil {
		return summary, err
	}

	formatter := mgf.New(options.Config, options.Config)
	for _, path := range files {
		select {
		case <-ctx.Done():
			return summary, ctx.Err()
		default:
		}

		nests, err := txt.ParseFile(path)
		if err != nil {
			return summary, fmt.Errorf("parse %q: %w", path, err)
		}
		summary.Files++
		summary.Nests += len(nests)

		for i, nest := range nests {
			index := i + 1
			content, err := formatter.Format(nest, index)
			if err != nil {
				return summary, fmt.Errorf("format %q nest %d: %w", path, index, err)
			}
			name, err := outputName(nest, index)
			if err != nil {
				return summary, fmt.Errorf("%q nest %d: %w", path, index, err)
			}

			stock := nest.StockBar()
			parts := nest.Parts()
			section := ""
			if len(parts) > 0 {
				section = parts[0].Section()
			}
			directory := filepath.Join(options.OutDir, stock.Project(), section)
			target := filepath.Join(directory, name+".mgf")
			if err := writeFile(target, content); err != nil {
				return summary, err
			}
			summary.MGFs++
			logf(options.Log, "wrote %s\n", target)

			if straight := formatter.StraightCut(nest); straight != "" {
				target := filepath.Join(directory, name+"_strightcut.txt")
				if err := writeFile(target, straight); err != nil {
					return summary, err
				}
				summary.StraightCuts++
			}
		}

		if !options.DontRename {
			renamed := filepath.Join(filepath.Dir(path), "_"+filepath.Base(path))
			if err := os.Rename(path, renamed); err != nil {
				return summary, fmt.Errorf("rename %q: %w", path, err)
			}
			summary.Renamed++
		}
	}

	return summary, nil
}

func txtFiles(directory string) ([]string, error) {
	files := make([]string, 0)
	err := filepath.WalkDir(directory, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		if strings.EqualFold(filepath.Ext(path), ".txt") && !strings.HasPrefix(filepath.Base(path), "_") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("scan %q: %w", directory, err)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no TXT files found in %q", directory)
	}
	return files, nil
}

func outputName(nest domain.Nest, index int) (string, error) {
	if len(nest.Parts()) == 0 {
		return "", fmt.Errorf("nest has no profiles")
	}
	if nest.Parts()[0].Section() == "" {
		return "", fmt.Errorf("profile has no section")
	}
	return nest.Name(index), nil
}

func writeFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create %q: %w", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write %q: %w", path, err)
	}
	return nil
}

func logf(writer io.Writer, format string, args ...any) {
	if writer == nil {
		return
	}
	_, _ = fmt.Fprintf(writer, format, args...)
}
