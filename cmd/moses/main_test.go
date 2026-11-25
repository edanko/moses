package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func runCLI(t *testing.T, args ...string) (string, error) {
	t.Helper()

	var output bytes.Buffer
	err := run(args, &output, &output)
	return output.String(), err
}

func TestVersionCommand(t *testing.T) {
	output, err := runCLI(t, "version")
	if err != nil {
		t.Fatalf("version: %v", err)
	}
	if !strings.Contains(output, "moses dev") {
		t.Fatalf("output = %q, want version line", output)
	}
	if !strings.Contains(output, "commit: none") {
		t.Fatalf("output = %q, want commit line", output)
	}
}

func TestGlobalConfigAfterCommand(t *testing.T) {
	output, err := runCLI(t, "version", "--config", "missing.yaml")
	if err != nil {
		t.Fatalf("version: %v", err)
	}
	if !strings.Contains(output, "moses dev") {
		t.Fatalf("output = %q, want version line", output)
	}
}

func TestUnknownFlagIsUsageError(t *testing.T) {
	_, err := runCLI(t, "--not-a-flag")
	if !errors.Is(err, errUsage) {
		t.Fatalf("error = %v, want usage error", err)
	}
}

func TestUnknownCommandIsUsageError(t *testing.T) {
	_, err := runCLI(t, "unknown")
	if !errors.Is(err, errUsage) {
		t.Fatalf("error = %v, want usage error", err)
	}
}

func TestConfigFlagRequiresValue(t *testing.T) {
	_, err := runCLI(t, "version", "--config")
	if !errors.Is(err, errUsage) {
		t.Fatalf("error = %v, want usage error", err)
	}
}

func TestRootUsageWithoutArguments(t *testing.T) {
	output, err := runCLI(t)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if !strings.Contains(output, "Commands:") {
		t.Fatalf("output = %q, want root usage", output)
	}
}

func TestConvUsageWithoutConverter(t *testing.T) {
	output, err := runCLI(t, "conv")
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if !strings.Contains(output, "Converters:") {
		t.Fatalf("output = %q, want converter usage", output)
	}
}

func TestConvRejectsUnknownConverter(t *testing.T) {
	_, err := runCLI(t, "conv", "unknown")
	if !errors.Is(err, errUsage) {
		t.Fatalf("error = %v, want usage error", err)
	}
}

func TestParserRejectsUnexpectedArguments(t *testing.T) {
	_, err := runCLI(t, "parser", "extra")
	if !errors.Is(err, errUsage) {
		t.Fatalf("error = %v, want usage error", err)
	}
}

func TestConvHelp(t *testing.T) {
	output, err := runCLI(t, "conv", "10022", "--help")
	if err != nil {
		t.Fatalf("help: %v", err)
	}
	for _, expected := range []string{"--in", "--out", "--algorithm", "--remnants"} {
		if !strings.Contains(output, expected) {
			t.Fatalf("help output missing %q:\n%s", expected, output)
		}
	}
}

func TestRootHelp(t *testing.T) {
	output, err := runCLI(t, "--help")
	if err != nil {
		t.Fatalf("help: %v", err)
	}
	for _, expected := range []string{"parser", "conv", "version", "--config"} {
		if !strings.Contains(output, expected) {
			t.Fatalf("root help missing %q:\n%s", expected, output)
		}
	}
}

func TestParserHelp(t *testing.T) {
	output, err := runCLI(t, "parser", "--help")
	if err != nil {
		t.Fatalf("help: %v", err)
	}
	for _, expected := range []string{"--in", "--out", "--dont-rename"} {
		if !strings.Contains(output, expected) {
			t.Fatalf("parser help missing %q:\n%s", expected, output)
		}
	}
}
