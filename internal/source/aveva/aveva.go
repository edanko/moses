// Package aveva parses AVEVA "Part sketch & list" CSV files paired with
// their LST metadata files.
package aveva

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"uuid"

	"github.com/edanko/moses/internal/domain"
)

var (
	reProfileDesignation = regexp.MustCompile(`(?m)(?:PP|RP)\/(\d+)x(\d+\.?\d?)`)
	reQuality            = regexp.MustCompile(`(?m)Qual :\s+(\w+)`)
)

// ParseFile reads a CSV/LST pair and returns aggregated parts.
func ParseFile(csvPath, project string) ([]domain.Part, error) {
	if strings.TrimSpace(project) == "" {
		return nil, fmt.Errorf("project is required")
	}

	dimension, quality, err := readMetadata(csvPath)
	if err != nil {
		return nil, err
	}
	profileType, err := domain.ParseProfileDesignation(dimension)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", csvPath, err)
	}

	file, err := os.Open(csvPath)
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", csvPath, err)
	}
	defer func() { _ = file.Close() }()

	reader := csv.NewReader(file)
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("%q: read header: %w", csvPath, err)
	}

	profiles := make(map[string]*domain.Part)
	order := make([]string, 0)

	for row := 2; ; row++ {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("%q row %d: %w", csvPath, row, err)
		}
		if len(record) < 31 {
			return nil, fmt.Errorf("%q row %d: expected at least 31 columns, got %d", csvPath, row, len(record))
		}

		name := strings.Split(record[0], "-")
		if len(name) < 2 || name[0] == "" {
			return nil, fmt.Errorf("%q row %d: invalid part name %q", csvPath, row, record[0])
		}
		position := name[len(name)-1]
		if position == "" {
			return nil, fmt.Errorf("%q row %d: empty position in %q", csvPath, row, record[0])
		}
		position = position[:len(position)-1]

		key := project + name[0] + position
		if existing, ok := profiles[key]; ok {
			existing.SetQuantity(existing.Quantity() + 1)
			continue
		}

		length, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return nil, fmt.Errorf("%q row %d: invalid length %q: %w", csvPath, row, record[3], err)
		}

		left, err := parseEnd(record, 5, 16, 6, 7, 8, 9, 10, 11, 12, 13, 14)
		if err != nil {
			return nil, fmt.Errorf("%q row %d left end: %w", csvPath, row, err)
		}
		right, err := parseEnd(record, 19, 30, 20, 21, 22, 23, 24, 25, 26, 27, 28)
		if err != nil {
			return nil, fmt.Errorf("%q row %d right end: %w", csvPath, row, err)
		}

		profile := domain.NewPart(domain.PartParams{
			ID:          uuid.NewV4(),
			Project:     project,
			Section:     name[0],
			PositionNo:  position,
			Designation: profileType,
			Grade:       quality,
			Length:      length,
			Quantity:    1,
			LeftEndCut:  left,
			RightEndCut: right,
		})
		profiles[key] = &profile
		order = append(order, key)
	}

	result := make([]domain.Part, 0, len(order))
	for _, key := range order {
		result = append(result, *profiles[key])
	}
	return result, nil
}

func parseEnd(
	record []string,
	kindIndex, bevelIndex int,
	aIndex, bIndex, cIndex, snipeIndex, r1Index, r2Index, v1Index, v2Index, v3Index int,
) (domain.EndCut, error) {
	if strings.TrimSpace(record[kindIndex]) == "" {
		return domain.EndCut{}, fmt.Errorf("missing end cut type")
	}

	var builder strings.Builder
	builder.WriteString(record[kindIndex])
	builder.WriteString(bevel(record[bevelIndex]))

	appendParam(&builder, "a", record[aIndex])
	appendParam(&builder, "b", record[bIndex])
	appendParam(&builder, "c", record[cIndex])
	appendParam(&builder, "snipe", record[snipeIndex])
	appendParam(&builder, "r1", record[r1Index])
	appendParam(&builder, "r2", record[r2Index])
	if err := appendAngle(&builder, "v1", record[v1Index], true); err != nil {
		return domain.EndCut{}, err
	}
	if err := appendAngle(&builder, "v2", record[v2Index], false); err != nil {
		return domain.EndCut{}, err
	}
	appendParam(&builder, "v3", record[v3Index])

	return domain.ParseEndCut(builder.String())
}

func appendParam(builder *strings.Builder, name, value string) {
	if strings.TrimSpace(value) == "" {
		return
	}
	builder.WriteString(" ")
	builder.WriteString(name)
	builder.WriteString("=")
	builder.WriteString(value)
}

func appendAngle(builder *strings.Builder, name, value string, subtract bool) error {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("invalid %s %q: %w", name, value, err)
	}
	if subtract {
		parsed -= 90
	} else {
		parsed = 90 - parsed
	}
	builder.WriteString(" ")
	builder.WriteString(name)
	builder.WriteString("=")
	builder.WriteString(strconv.FormatFloat(parsed, 'f', -1, 64))
	return nil
}

func bevel(value string) string {
	if value == "" || value == "0" {
		return ""
	}
	switch value {
	case "177", "154", "180", "227", "143", "-245":
		return ""
	case "-131":
		return "-131"
	default:
		return "-" + value
	}
}

func readMetadata(csvPath string) (string, string, error) {
	lstPath := strings.TrimSuffix(csvPath, ".csv") + ".lst"
	file, err := os.Open(lstPath)
	if err != nil {
		return "", "", fmt.Errorf("open LST %q: %w", lstPath, err)
	}
	defer func() { _ = file.Close() }()

	var dimension, quality string
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return "", "", fmt.Errorf("read LST %q: %w", lstPath, err)
		}

		if strings.Contains(line, "Type/Dim.") {
			match := reProfileDesignation.FindStringSubmatch(line)
			if len(match) == 3 {
				dimension = "RP" + match[1] + "X" + match[2]
			}
		}
		if strings.Contains(line, "Qual") {
			match := reQuality.FindStringSubmatch(line)
			if len(match) == 2 {
				quality = match[1]
			}
		}

		if errors.Is(err, io.EOF) {
			break
		}
		if dimension != "" && quality != "" {
			break
		}
	}

	if dimension == "" || quality == "" {
		return "", "", fmt.Errorf("LST %q: missing Type/Dim. or Qual", lstPath)
	}
	dimension = strings.TrimSuffix(dimension, ".0")
	if quality == "D40_B" {
		quality = "D40-B"
	}
	return dimension, quality, nil
}
