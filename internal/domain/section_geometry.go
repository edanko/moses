package domain

import "fmt"

// Shape is the MGF cross-section shape of a profile.
type Shape string

const (
	ShapeUnknown Shape = ""
	ShapeHP      Shape = "HP"
	ShapeFB      Shape = "FB"
	ShapeL       Shape = "L"
)

// SectionGeometry carries the cross-section geometry used to write machine output.
type SectionGeometry struct {
	Shape        Shape
	B            float64
	C            float64
	S            float64
	R            float64
	A            float64
	R1           float64
	R2           float64
	T            float64
	MassPerMeter float64
}

// Validate checks that the parameters required by the shape are present.
func (s SectionGeometry) Validate() error {
	switch s.Shape {
	case ShapeHP:
		if s.B <= 0 || s.C <= 0 || s.S <= 0 || s.R <= 0 {
			return fmt.Errorf("HP section geometry requires positive B, C, S and R")
		}
	case ShapeFB:
		if s.B <= 0 || s.S <= 0 {
			return fmt.Errorf("FB section geometry requires positive B and S")
		}
	case ShapeL:
		if s.A <= 0 || s.B <= 0 || s.R1 <= 0 || s.R2 <= 0 || s.S <= 0 || s.T <= 0 {
			return fmt.Errorf("section geometry for shape L requires positive A, B, R1, R2, S and T")
		}
	default:
		return fmt.Errorf("unknown profile shape %q", s.Shape)
	}
	return nil
}

// IsValid reports whether the spec has a known shape.
func (s SectionGeometry) IsValid() bool {
	return s.Shape != ShapeUnknown
}
