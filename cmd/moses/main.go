package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// errUsage marks command-line usage errors so main can exit with code 2.
var errUsage = errors.New("usage")

type cli struct {
	ctx        context.Context
	out        io.Writer
	errOut     io.Writer
	configPath string
}

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "moses:", err)
		if errors.Is(err, errUsage) {
			os.Exit(2)
		}
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	configPath, rest, err := extractConfig(args)
	if err != nil {
		return err
	}

	c := &cli{
		ctx:        context.Background(),
		out:        stdout,
		errOut:     stderr,
		configPath: configPath,
	}
	return c.run(rest)
}

func (c *cli) run(args []string) error {
	if len(args) == 0 {
		c.printRootUsage()
		return nil
	}

	switch args[0] {
	case "help", "-h", "--help":
		c.printRootUsage()
		return nil
	case "version":
		return c.runVersion(args[1:])
	case "parser":
		return c.runParser(args[1:])
	case "conv":
		return c.runConv(args[1:])
	default:
		if strings.HasPrefix(args[0], "-") {
			return fmt.Errorf("%w: unknown option %q", errUsage, args[0])
		}
		return fmt.Errorf("%w: unknown command %q", errUsage, args[0])
	}
}

func extractConfig(args []string) (string, []string, error) {
	configPath := "config.yaml"
	rest := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		switch arg := args[i]; {
		case arg == "--config" || arg == "-config":
			if i+1 >= len(args) {
				return "", nil, fmt.Errorf("%w: --config requires a value", errUsage)
			}
			configPath = args[i+1]
			i++
		case strings.HasPrefix(arg, "--config="):
			configPath = strings.TrimPrefix(arg, "--config=")
		case strings.HasPrefix(arg, "-config="):
			configPath = strings.TrimPrefix(arg, "-config=")
		default:
			rest = append(rest, arg)
		}
	}
	return configPath, rest, nil
}

func (c *cli) parseFlags(fs *flag.FlagSet, args []string) (bool, error) {
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return true, nil
		}
		return false, fmt.Errorf("%w: %w", errUsage, err)
	}
	if fs.NArg() != 0 {
		return false, fmt.Errorf("%w: unexpected arguments: %s", errUsage, strings.Join(fs.Args(), " "))
	}
	return false, nil
}

func (c *cli) printCommandUsage(fs *flag.FlagSet, name, description string) {
	_, _ = fmt.Fprintf(c.out, "%s\n\nUsage:\n  moses %s [options]\n\nOptions:\n", description, name)
	fs.VisitAll(func(f *flag.Flag) {
		_, _ = fmt.Fprintf(c.out, "  --%s\n        %s", f.Name, f.Usage)
		if f.DefValue != "" && f.DefValue != "false" {
			_, _ = fmt.Fprintf(c.out, " (default %s)", f.DefValue)
		}
		_, _ = fmt.Fprintln(c.out)
	})
}

func (c *cli) printRootUsage() {
	_, _ = fmt.Fprint(c.out, `Moses converts profile data between production formats and writes machine output.

Usage:
  moses [--config FILE] <command> [options]

Commands:
  parser    Convert nested TXT files to MGF machine output
  conv      Convert production source files into nested TXT and reports
  version   Print version and build information

Global options:
  --config FILE   path to the configuration file (default "config.yaml")
`)
}
