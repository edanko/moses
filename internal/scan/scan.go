package scan

import (
	"strings"
	"text/scanner"
)

type Scan struct {
	scanner.Scanner
}

func NewScanner(str string) Scan {
	var s Scan

	reader := strings.NewReader(str)

	s.Init(reader)
	s.Mode = scanner.ScanChars
	s.Whitespace = 0

	return s
}

func (s *Scan) ReadNextLine() string {
	var result string

	ch := s.Scan()

	for ch != '\n' {
		if ch == scanner.EOF {
			if len(result) > 0 {
				return result
			}
			return "eof"
		}
		if ch != '\r' {
			result += s.TokenText()
		}
		ch = s.Scan()
	}
	return result
}

func (s *Scan) ReadLine() string {
	var result string

	ch := s.Scan()
	for ch != '\n' {
		if ch == scanner.EOF {
			if len(result) > 0 {
				break
			}
			return ""
		}
		if ch != '\r' {
			result += s.TokenText()
		}
		ch = s.Scan()
	}
	result = strings.ToLower(result)
	if result == " -0" {
		result = "0"
	}
	return result
}
