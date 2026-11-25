// Package scan provides the line scanner used by the legacy text formats.
package scan

import (
	"strings"
	"text/scanner"
)

const eof = "eof"

// Scan reads lines while preserving the legacy tokenization rules.
type Scan struct {
	scanner.Scanner
}

// New creates a scanner over text.
func New(text string) Scan {
	var s Scan
	s.Init(strings.NewReader(text))
	s.Mode = scanner.ScanChars
	s.Whitespace = 0
	return s
}

// ReadNextLine returns the next line without lowercasing or trimming.
func (s *Scan) ReadNextLine() string {
	var result strings.Builder
	ch := s.Scan()
	for ch != '\n' {
		if ch == scanner.EOF {
			if result.Len() > 0 {
				return result.String()
			}
			return eof
		}
		if ch != '\r' {
			result.WriteString(s.TokenText())
		}
		ch = s.Scan()
	}
	return result.String()
}

// ReadLine returns the next lowercased line. The legacy formats use "0" in
// place of the textual "-0".
func (s *Scan) ReadLine() string {
	var result strings.Builder
	ch := s.Scan()
	for ch != '\n' {
		if ch == scanner.EOF {
			if result.Len() == 0 {
				return ""
			}
			break
		}
		if ch != '\r' {
			result.WriteString(s.TokenText())
		}
		ch = s.Scan()
	}
	line := strings.ToLower(result.String())
	if line == " -0" {
		return "0"
	}
	return line
}
