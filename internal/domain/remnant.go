package domain

import "uuid"

// Remnant is an unused offcut recorded in the remnant CSV.
type Remnant struct {
	ID          int
	Used        bool
	Length      float64
	Project     string
	Source      string
	Designation ProfileDesignation
	Grade       string
	Marking     string
	ConsumedBy  string
}

// ToStockBar converts an unused remnant into a stock bar.
func (r *Remnant) ToStockBar() (*StockBar, error) {
	if r.Used {
		return nil, ErrRemnantConsumed
	}
	stock := NewStockBar(StockBarParams{
		ID:          uuid.NewV4(),
		Project:     r.Project,
		Designation: r.Designation,
		Grade:       r.Grade,
		Length:      r.Length,
		Type:        StockTypeRemnant,
		Remnant:     r,
	})
	if err := stock.Validate(); err != nil {
		return nil, err
	}
	return &stock, nil
}

// Consume marks the remnant as used by a nest.
func (r *Remnant) Consume(nestName string) {
	r.Used = true
	r.ConsumedBy = nestName
}

// Validate checks the required remnant fields.
func (r *Remnant) Validate() error {
	switch {
	case r.Project == "":
		return ErrProject
	case r.Source == "":
		return ErrRemnantSource
	case !r.Designation.IsValid():
		return ErrProfileDesignation
	case r.Grade == "":
		return ErrGrade
	case r.Length <= 0:
		return ErrLength
	}
	return nil
}
