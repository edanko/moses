package domain

import "uuid"

// Part is one production part to be cut from a stock bar.
type Part struct {
	id uuid.UUID

	project    string
	section    string
	positionNo string

	designation ProfileDesignation
	grade       string
	length      float64
	quantity    int

	leftEndCut  EndCut
	rightEndCut EndCut
	cutouts     []Cutout
}

// PartParams describes a production part without exposing its layout.
type PartParams struct {
	ID          uuid.UUID
	Project     string
	Section     string
	PositionNo  string
	Designation ProfileDesignation
	Grade       string
	Length      float64
	Quantity    int
	LeftEndCut  EndCut
	RightEndCut EndCut
	Cutouts     []Cutout
}

// NewPart creates a Part. Validation is explicit via Validate.
func NewPart(params PartParams) Part {
	cutouts := make([]Cutout, len(params.Cutouts))
	for i := range params.Cutouts {
		cutouts[i] = params.Cutouts[i].clone()
	}
	return Part{
		id:          params.ID,
		project:     params.Project,
		section:     params.Section,
		positionNo:  params.PositionNo,
		designation: params.Designation,
		grade:       params.Grade,
		length:      params.Length,
		quantity:    params.Quantity,
		leftEndCut:  params.LeftEndCut,
		rightEndCut: params.RightEndCut,
		cutouts:     cutouts,
	}
}

// ID returns the part identifier.
func (p Part) ID() uuid.UUID {
	return p.id
}

// Project returns the project number.
func (p Part) Project() string {
	return p.project
}

// Section returns the assembly section or block.
func (p Part) Section() string {
	return p.section
}

// PositionNo returns the drawing position number.
func (p Part) PositionNo() string {
	return p.positionNo
}

// Designation returns the normalized profile designation.
func (p Part) Designation() ProfileDesignation {
	return p.designation
}

// DesignationKey returns the lowercase configuration key.
func (p Part) DesignationKey() string {
	return p.designation.Key()
}

// DesignationName returns the canonical uppercase designation.
func (p Part) DesignationName() string {
	return p.designation.String()
}

// Grade returns the steel grade.
func (p Part) Grade() string {
	return p.grade
}

// Length returns the raw part length in millimetres.
func (p Part) Length() float64 {
	return p.length
}

// Quantity returns how many copies of this part are required.
func (p Part) Quantity() int {
	return p.quantity
}

// SetQuantity replaces the required copy count.
func (p *Part) SetQuantity(quantity int) {
	p.quantity = quantity
}

// LeftEndCut returns the left end cut.
func (p Part) LeftEndCut() EndCut {
	return p.leftEndCut
}

// RightEndCut returns the right end cut.
func (p Part) RightEndCut() EndCut {
	return p.rightEndCut
}

// LeftEndKind returns the left end-cut kind.
func (p Part) LeftEndKind() string {
	return p.leftEndCut.Kind()
}

// RightEndKind returns the right end-cut kind.
func (p Part) RightEndKind() string {
	return p.rightEndCut.Kind()
}

// Cutouts returns a deep copy of the internal cuts.
func (p Part) Cutouts() []Cutout {
	cutouts := make([]Cutout, len(p.cutouts))
	for i := range p.cutouts {
		cutouts[i] = p.cutouts[i].clone()
	}
	return cutouts
}

func (p Part) clone() Part {
	cutouts := make([]Cutout, len(p.cutouts))
	for i := range p.cutouts {
		cutouts[i] = p.cutouts[i].clone()
	}
	p.cutouts = cutouts
	return p
}

// MirrorCutouts mirrors internal cuts around the profile centre.
func (p *Part) MirrorCutouts() {
	for i := range p.cutouts {
		x := p.Length() - p.cutouts[i].Param("x")
		p.cutouts[i].SetParam("x", x)
	}
}

// Validate checks the required part fields.
func (p Part) Validate() error {
	switch {
	case p.project == "":
		return ErrProject
	case p.section == "":
		return ErrSection
	case p.positionNo == "":
		return ErrPositionNumber
	case !p.designation.IsValid():
		return ErrProfileDesignation
	case p.grade == "":
		return ErrGrade
	case p.length <= 0:
		return ErrLength
	case p.quantity <= 0:
		return ErrQuantity
	case p.leftEndCut.Kind() == "" || p.rightEndCut.Kind() == "":
		return ErrEndCut
	}
	return nil
}
