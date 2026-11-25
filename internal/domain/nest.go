package domain

import (
	"strconv"
	"strings"
	"uuid"
)

// Nest is one stock bar with the parts assigned to it.
type Nest struct {
	id       uuid.UUID
	stock    StockBar
	parts    []Part
	consumed *Remnant
}

// NestParams describes a nest.
type NestParams struct {
	ID       uuid.UUID
	StockBar StockBar
	Parts    []Part
	Consumed *Remnant
}

// NewNest creates a nest.
func NewNest(params NestParams) *Nest {
	parts := make([]Part, len(params.Parts))
	for i := range params.Parts {
		parts[i] = params.Parts[i].clone()
	}
	return &Nest{
		id:       params.ID,
		stock:    params.StockBar,
		parts:    parts,
		consumed: params.Consumed,
	}
}

// ID returns the nest identifier.
func (n Nest) ID() uuid.UUID {
	return n.id
}

// StockBar returns the stock bar used by the nest.
func (n Nest) StockBar() StockBar {
	return n.stock
}

// Parts returns a deep copy of the nested parts.
func (n Nest) Parts() []Part {
	parts := make([]Part, len(n.parts))
	for i := range n.parts {
		parts[i] = n.parts[i].clone()
	}
	return parts
}

// ConsumedRemnant returns the remnant bar this nest consumed, if any.
func (n Nest) ConsumedRemnant() *Remnant {
	return n.consumed
}

// StockLength returns the length of the nest's stock bar.
func (n Nest) StockLength() float64 {
	return n.stock.Length()
}

// TotalPartLength returns the sum of raw part lengths.
func (n Nest) TotalPartLength() float64 {
	var total float64
	for _, part := range n.parts {
		total += part.Length()
	}
	return total
}

// ShortName returns section-DIM-GRADE-index, matching converter reports.
func (n Nest) ShortName(index int) string {
	parts := make([]string, 0, 4)
	if len(n.parts) > 0 && n.parts[0].Section() != "" {
		parts = append(parts, n.parts[0].Section())
	}
	parts = append(parts, strings.ToUpper(n.stock.DesignationName()))
	if n.stock.Grade() != "" {
		parts = append(parts, strings.ToUpper(n.stock.Grade()))
	}
	parts = append(parts, strconv.Itoa(index))
	return strings.Join(parts, "-")
}

// Name returns project-section-DIM-GRADE-index, matching parser MGF names.
func (n Nest) Name(index int) string {
	short := n.ShortName(index)
	if project := n.stock.Project(); project != "" {
		return project + "-" + short
	}
	return short
}

// HasUsefulRemnant reports whether the leftover is at least minLength.
func (n Nest) HasUsefulRemnant(minLength float64) bool {
	return n.stock.Length()-n.TotalPartLength() >= minLength
}

// RemnantLength returns the useful leftover length.
func (n Nest) RemnantLength(minLength float64) float64 {
	if n.HasUsefulRemnant(minLength) {
		return n.stock.Length() - n.TotalPartLength()
	}
	return 0
}

// CreateRemnant creates the remnant item produced by this nest, or nil when
// the leftover is below minLength.
func (n Nest) CreateRemnant(name string, minLength float64) *Remnant {
	if !n.HasUsefulRemnant(minLength) {
		return nil
	}
	length := n.RemnantLength(minLength)
	return &Remnant{
		Length:      length,
		Project:     n.stock.Project(),
		Source:      name,
		Designation: n.stock.Designation(),
		Grade:       n.stock.Grade(),
		Marking:     name + "-L" + strconv.FormatFloat(length, 'f', -1, 64),
	}
}

// Validate checks the nest aggregate invariants.
func (n Nest) Validate() error {
	if err := n.stock.Validate(); err != nil {
		return err
	}
	if len(n.parts) == 0 {
		return ErrNoParts
	}
	for _, part := range n.parts {
		if err := part.Validate(); err != nil {
			return err
		}
	}
	if used := n.TotalPartLength(); !n.stock.CanHold(used) {
		return ErrLength
	}
	return nil
}
