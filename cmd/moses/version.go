package main

import "fmt"

// These values are injected with -ldflags at build time.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func (c *cli) runVersion(args []string) error {
	if len(args) == 1 && (args[0] == "-h" || args[0] == "--help") {
		_, _ = fmt.Fprintln(c.out, "Print version and build information.")
		return nil
	}
	if len(args) != 0 {
		return fmt.Errorf("%w: version takes no arguments", errUsage)
	}

	_, err := fmt.Fprintf(
		c.out,
		"moses %s\ncommit: %s\nbuilt: %s\n",
		version,
		commit,
		date,
	)
	return err
}
