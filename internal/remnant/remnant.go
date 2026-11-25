// Package remnant stores unused offcuts in the production CSV file.
package remnant

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/edanko/moses/internal/domain"
)

const (
	delimiter = ';'
	header    = "id;prj;sec;dim;quality;rem_len;marking_text;used_in"
)

// Store is the in-memory view of the remnant CSV.
type Store struct {
	path   string
	items  []*domain.Remnant
	nextID int
}

// Load reads an existing remnant CSV. A missing file is an empty store.
func Load(path string) (*Store, error) {
	store := &Store{path: path, items: make([]*domain.Remnant, 0)}
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return store, nil
		}
		return nil, fmt.Errorf("open remnants %q: %w", path, err)
	}
	defer func() { _ = file.Close() }()

	reader := csv.NewReader(file)
	reader.Comma = delimiter
	reader.FieldsPerRecord = -1

	if _, err := reader.Read(); err != nil {
		if errors.Is(err, io.EOF) {
			return store, nil
		}
		return nil, fmt.Errorf("read remnants header %q: %w", path, err)
	}

	for row := 2; ; row++ {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read remnants %q row %d: %w", path, row, err)
		}
		if len(record) < 8 {
			return nil, fmt.Errorf("read remnants %q row %d: expected 8 columns", path, row)
		}

		id, err := strconv.Atoi(strings.TrimSpace(record[0]))
		if err != nil {
			return nil, fmt.Errorf("remnants row %d: invalid id %q: %w", row, record[0], err)
		}
		profileType, err := domain.ParseProfileDesignation(record[3])
		if err != nil {
			return nil, fmt.Errorf("remnants row %d: %w", row, err)
		}
		length, err := strconv.ParseFloat(strings.TrimSpace(record[5]), 64)
		if err != nil {
			return nil, fmt.Errorf("remnants row %d: invalid length %q: %w", row, record[5], err)
		}

		item := &domain.Remnant{
			ID:          id,
			Length:      length,
			Project:     strings.TrimSpace(record[1]),
			Source:      strings.TrimSpace(record[2]),
			Designation: profileType,
			Grade:       strings.TrimSpace(record[4]),
			Marking:     strings.TrimSpace(record[6]),
			ConsumedBy:  strings.TrimSpace(record[7]),
		}
		item.Used = item.ConsumedBy != ""
		store.items = append(store.items, item)
		if id > store.nextID {
			store.nextID = id
		}
	}
	return store, nil
}

// Available returns unused remnants for a project, designation and grade.
func (s *Store) Available(project string, designation domain.ProfileDesignation, grade string) []domain.Remnant {
	result := make([]domain.Remnant, 0)
	for _, item := range s.items {
		if item.Used {
			continue
		}
		if item.Project != project {
			continue
		}
		if !item.Designation.Equal(designation) {
			continue
		}
		if !strings.EqualFold(item.Grade, grade) {
			continue
		}
		result = append(result, *item)
	}
	return result
}

// MarkUsed records the nest that consumed a remnant.
func (s *Store) MarkUsed(id int, nestName string) {
	for _, item := range s.items {
		if item.ID == id {
			item.Consume(nestName)
			return
		}
	}
}

// Add assigns an id and stores a new remnant.
func (s *Store) Add(item domain.Remnant) *domain.Remnant {
	s.nextID++
	item.ID = s.nextID
	item.Used = false
	item.ConsumedBy = ""
	stored := item
	s.items = append(s.items, &stored)
	return &stored
}

// Items returns a copy of the stored remnants.
func (s *Store) Items() []domain.Remnant {
	result := make([]domain.Remnant, len(s.items))
	for i, item := range s.items {
		result[i] = *item
	}
	return result
}

// Save writes the store atomically.
func (s *Store) Save() error {
	if s.path == "" {
		return fmt.Errorf("remnants: path is empty")
	}

	directory := filepath.Dir(s.path)
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return fmt.Errorf("create remnant directory: %w", err)
	}

	temporary, err := os.CreateTemp(directory, ".remnants-*.csv")
	if err != nil {
		return fmt.Errorf("create remnant temp file: %w", err)
	}
	temporaryPath := temporary.Name()
	defer func() { _ = os.Remove(temporaryPath) }()
	if err := temporary.Chmod(0o644); err != nil {
		_ = temporary.Close()
		return fmt.Errorf("set remnant permissions: %w", err)
	}

	writer := csv.NewWriter(temporary)
	writer.Comma = delimiter
	if err := writer.Write(strings.Split(header, ";")); err != nil {
		_ = temporary.Close()
		return fmt.Errorf("write remnants header: %w", err)
	}

	for _, item := range s.items {
		usedIn := item.ConsumedBy
		if item.Used && usedIn == "" {
			usedIn = "used"
		}
		record := []string{
			strconv.Itoa(item.ID),
			item.Project,
			item.Source,
			item.Designation.String(),
			item.Grade,
			strconv.FormatFloat(item.Length, 'f', -1, 64),
			item.Marking,
			usedIn,
		}
		if err := writer.Write(record); err != nil {
			_ = temporary.Close()
			return fmt.Errorf("write remnant row: %w", err)
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		_ = temporary.Close()
		return fmt.Errorf("flush remnants: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close remnants: %w", err)
	}
	if err := os.Rename(temporaryPath, s.path); err != nil {
		return fmt.Errorf("replace remnants %q: %w", s.path, err)
	}
	return nil
}
