package domain

import "uuid"

// StockBar is a bar or remnant that parts can be nested onto.
type StockBar struct {
	id uuid.UUID

	project     string
	designation ProfileDesignation
	grade       string
	length      float64
	stockType   StockType
	remnant     *Remnant
}

// StockBarParams describes a stock bar.
type StockBarParams struct {
	ID          uuid.UUID
	Project     string
	Designation ProfileDesignation
	Grade       string
	Length      float64
	Type        StockType
	Remnant     *Remnant
}

// NewStockBar creates a stock bar.
func NewStockBar(params StockBarParams) StockBar {
	return StockBar{
		id:          params.ID,
		project:     params.Project,
		designation: params.Designation,
		grade:       params.Grade,
		length:      params.Length,
		stockType:   params.Type,
		remnant:     params.Remnant,
	}
}

// ID returns the stock identifier.
func (s StockBar) ID() uuid.UUID {
	return s.id
}

// Project returns the project number.
func (s StockBar) Project() string {
	return s.project
}

// Designation returns the normalized profile designation.
func (s StockBar) Designation() ProfileDesignation {
	return s.designation
}

// DesignationKey returns the lowercase configuration key.
func (s StockBar) DesignationKey() string {
	return s.designation.Key()
}

// DesignationName returns the canonical uppercase designation.
func (s StockBar) DesignationName() string {
	return s.designation.String()
}

// Grade returns the steel grade.
func (s StockBar) Grade() string {
	return s.grade
}

// Length returns the stock length in millimetres.
func (s StockBar) Length() float64 {
	return s.length
}

// Type returns the stock type.
func (s StockBar) Type() StockType {
	return s.stockType
}

// IsRemnant reports whether this bar is a consumed remnant.
func (s StockBar) IsRemnant() bool {
	return s.stockType == StockTypeRemnant
}

// IsOverlength reports whether this bar is an explicit over-length bar.
func (s StockBar) IsOverlength() bool {
	return s.stockType == StockTypeOverlength
}

// ConsumedRemnant returns the remnant this stock bar represents, if any.
func (s StockBar) ConsumedRemnant() *Remnant {
	return s.remnant
}

// CanHold reports whether used length fits within the bar.
func (s StockBar) CanHold(usedLength float64) bool {
	return usedLength <= s.length
}

// Validate checks the required stock fields.
func (s StockBar) Validate() error {
	switch {
	case !s.designation.IsValid():
		return ErrProfileDesignation
	case s.grade == "":
		return ErrGrade
	case s.length <= 0:
		return ErrLength
	case !s.stockType.IsValid():
		return ErrStockType
	}
	return nil
}
