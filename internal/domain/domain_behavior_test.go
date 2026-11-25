package domain

import (
	"errors"
	"testing"

	"uuid"
)

func testDesignation(t *testing.T, raw string) ProfileDesignation {
	t.Helper()
	designation, err := ParseProfileDesignation(raw)
	if err != nil {
		t.Fatalf("parse designation %q: %v", raw, err)
	}
	return designation
}

func testPart(t *testing.T) Part {
	t.Helper()
	end := NewEndCut("E11", map[string]float64{"X": 0, "Y": 0})
	return NewPart(PartParams{
		ID:          uuid.NewV4(),
		Project:     "056001",
		Section:     "05102",
		PositionNo:  "4835",
		Designation: testDesignation(t, "RP160X8"),
		Grade:       "A40",
		Length:      1000,
		Quantity:    1,
		LeftEndCut:  end,
		RightEndCut: end,
	})
}

func TestPartAccessors(t *testing.T) {
	part := testPart(t)
	if part.DesignationKey() != "rp160x8" || part.DesignationName() != "RP160X8" {
		t.Fatalf("designation accessors = %q/%q", part.DesignationKey(), part.DesignationName())
	}
	if part.LeftEndKind() != "E11" || part.RightEndKind() != "E11" {
		t.Fatalf("end kinds = %q/%q", part.LeftEndKind(), part.RightEndKind())
	}
}

func TestPartValidate(t *testing.T) {
	if err := testPart(t).Validate(); err != nil {
		t.Fatalf("valid part: %v", err)
	}

	invalid := testPart(t)
	invalid.grade = ""
	if err := invalid.Validate(); !errors.Is(err, ErrGrade) {
		t.Fatalf("empty grade = %v, want ErrGrade", err)
	}

	invalid = testPart(t)
	invalid.quantity = 0
	if err := invalid.Validate(); !errors.Is(err, ErrQuantity) {
		t.Fatalf("zero quantity = %v, want ErrQuantity", err)
	}

	invalid = testPart(t)
	invalid.leftEndCut = EndCut{}
	if err := invalid.Validate(); !errors.Is(err, ErrEndCut) {
		t.Fatalf("missing end cut = %v, want ErrEndCut", err)
	}
}

func TestPartMirrorCutouts(t *testing.T) {
	part := testPart(t)
	part.cutouts = []Cutout{*NewCutout("i1", map[string]float64{"x": 100})}
	part.MirrorCutouts()
	if got := part.Cutouts()[0].Param("x"); got != 900 {
		t.Fatalf("mirrored x = %v, want 900", got)
	}
}

func TestPartCutoutsAreCopies(t *testing.T) {
	cutout := NewCutout("i1", map[string]float64{"x": 100, "r": 10})
	part := NewPart(PartParams{
		ID:          uuid.NewV4(),
		Project:     "project",
		Section:     "section",
		PositionNo:  "1",
		Designation: testDesignation(t, "RP160X8"),
		Grade:       "A40",
		Length:      1000,
		Quantity:    1,
		LeftEndCut:  NewEndCut("E11", nil),
		RightEndCut: NewEndCut("E11", nil),
		Cutouts:     []Cutout{*cutout},
	})

	cutout.SetParam("x", 999)
	if got := part.Cutouts()[0].Param("x"); got != 100 {
		t.Fatalf("cutout was aliased: x = %v", got)
	}

	got := part.Cutouts()
	got[0].SetParam("x", 500)
	if got := part.Cutouts()[0].Param("x"); got != 100 {
		t.Fatalf("returned cutout was aliased: x = %v", got)
	}
}

func TestCutValueObjectsCopyParams(t *testing.T) {
	params := map[string]float64{"x": 100}
	end := NewEndCut("21", params)
	cutout := NewCutout("i1", params)

	params["x"] = 999
	if got := end.Param("x"); got != 100 {
		t.Fatalf("end params aliased: %v", got)
	}
	if got := cutout.Param("x"); got != 100 {
		t.Fatalf("cutout params aliased: %v", got)
	}

	endCopy := end.Params()
	endCopy["x"] = 500
	if got := end.Param("x"); got != 100 {
		t.Fatalf("returned end params aliased: %v", got)
	}
}

func TestNestRemnantAndValidation(t *testing.T) {
	part := testPart(t)
	stock := NewStockBar(StockBarParams{
		ID:          uuid.NewV4(),
		Project:     "056001",
		Designation: part.Designation(),
		Grade:       part.Grade(),
		Length:      6000,
		Type:        StockTypeBar,
	})
	nest := NewNest(NestParams{Parts: []Part{part}, StockBar: stock})

	if err := nest.Validate(); err != nil {
		t.Fatalf("valid nest: %v", err)
	}
	if got := nest.TotalPartLength(); got != 1000 {
		t.Fatalf("total part length = %v", got)
	}

	remnant := nest.CreateRemnant("nest-1", 500)
	if remnant == nil {
		t.Fatal("expected a remnant")
	}
	if remnant.Length != 5000 || remnant.Source != "nest-1" || remnant.Marking != "nest-1-L5000" {
		t.Fatalf("unexpected remnant: %+v", remnant)
	}

	empty := NewNest(NestParams{StockBar: stock})
	if err := empty.Validate(); !errors.Is(err, ErrNoParts) {
		t.Fatalf("empty nest error = %v, want ErrNoParts", err)
	}

	tooLong := NewStockBar(StockBarParams{
		ID:          uuid.NewV4(),
		Project:     "056001",
		Designation: part.Designation(),
		Grade:       part.Grade(),
		Length:      500,
		Type:        StockTypeBar,
	})
	if err := NewNest(NestParams{Parts: []Part{part}, StockBar: tooLong}).Validate(); !errors.Is(err, ErrLength) {
		t.Fatalf("over-capacity nest error = %v, want ErrLength", err)
	}
}

func TestStockBarAndRemnantBehaviour(t *testing.T) {
	stock := NewStockBar(StockBarParams{
		Designation: testDesignation(t, "RP160X8"),
		Grade:       "A40",
		Length:      6000,
		Type:        StockTypeRemnant,
	})
	if !stock.IsRemnant() || stock.IsOverlength() || !stock.CanHold(6000) || stock.CanHold(6001) {
		t.Fatalf("unexpected stock behaviour: %+v", stock)
	}

	remnant := &Remnant{Project: "p", Source: "s", Designation: stock.Designation(), Grade: "A40", Length: 1000}
	if err := remnant.Validate(); err != nil {
		t.Fatalf("valid remnant: %v", err)
	}
	remnant.Consume("nest-1")
	if !remnant.Used || remnant.ConsumedBy != "nest-1" {
		t.Fatalf("remnant not consumed: %+v", remnant)
	}
	if _, err := remnant.ToStockBar(); !errors.Is(err, ErrRemnantConsumed) {
		t.Fatalf("consumed remnant error = %v, want ErrRemnantConsumed", err)
	}
}

func TestStockType(t *testing.T) {
	values := StockType("").Values()
	if len(values) != 3 {
		t.Fatalf("values = %v", values)
	}
	for _, value := range []StockType{StockTypeBar, StockTypeRemnant, StockTypeOverlength} {
		if !value.IsValid() {
			t.Fatalf("%q should be valid", value)
		}
	}
	if StockType("invalid").IsValid() {
		t.Fatal("invalid stock type accepted")
	}
}

func TestSectionGeometryIsValid(t *testing.T) {
	if !(SectionGeometry{Shape: ShapeHP}).IsValid() || (SectionGeometry{}).IsValid() {
		t.Fatal("unexpected section geometry validity")
	}
}

func TestProfileDesignationKeyAndEqual(t *testing.T) {
	designation := testDesignation(t, "RP160X8")
	if designation.Key() != "rp160x8" {
		t.Fatalf("key = %q", designation.Key())
	}
	if !designation.Equal(testDesignation(t, "PP160*8")) {
		t.Fatal("aliases should compare equal")
	}
	if designation.Equal(testDesignation(t, "RP160X10")) {
		t.Fatal("different designations compare equal")
	}
}

func TestNestNames(t *testing.T) {
	part := testPart(t)
	stock := NewStockBar(StockBarParams{
		Project:     part.Project(),
		Designation: part.Designation(),
		Grade:       part.Grade(),
		Length:      6000,
		Type:        StockTypeBar,
	})
	nest := NewNest(NestParams{StockBar: stock, Parts: []Part{part}})

	if got, want := nest.ShortName(2), "05102-RP160X8-A40-2"; got != want {
		t.Fatalf("short name = %q, want %q", got, want)
	}
	if got, want := nest.Name(2), "056001-05102-RP160X8-A40-2"; got != want {
		t.Fatalf("name = %q, want %q", got, want)
	}
}

func TestRemnantToStockBar(t *testing.T) {
	remnant := &Remnant{
		Project:     "056001",
		Source:      "nest-1",
		Designation: testDesignation(t, "RP160X8"),
		Grade:       "A40",
		Length:      1000,
	}
	stock, err := remnant.ToStockBar()
	if err != nil {
		t.Fatalf("to stock: %v", err)
	}
	if !stock.IsRemnant() || stock.ConsumedRemnant() != remnant || stock.Length() != 1000 {
		t.Fatalf("unexpected stock: %+v", stock)
	}
	if stock.Project() != "056001" || stock.DesignationName() != "RP160X8" || stock.Grade() != "A40" {
		t.Fatalf("unexpected stock fields: %+v", stock)
	}
}

func TestStockBarValidate(t *testing.T) {
	valid := StockBarParams{
		Designation: testDesignation(t, "RP160X8"),
		Grade:       "A40",
		Length:      6000,
		Type:        StockTypeBar,
	}
	if err := NewStockBar(valid).Validate(); err != nil {
		t.Fatalf("valid stock: %v", err)
	}

	tests := []struct {
		name string
		edit func(*StockBarParams)
		want error
	}{
		{name: "designation", edit: func(p *StockBarParams) { p.Designation = ProfileDesignation{} }, want: ErrProfileDesignation},
		{name: "grade", edit: func(p *StockBarParams) { p.Grade = "" }, want: ErrGrade},
		{name: "length", edit: func(p *StockBarParams) { p.Length = 0 }, want: ErrLength},
		{name: "type", edit: func(p *StockBarParams) { p.Type = StockType("bad") }, want: ErrStockType},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := valid
			tt.edit(&params)
			if err := NewStockBar(params).Validate(); !errors.Is(err, tt.want) {
				t.Fatalf("error = %v, want %v", err, tt.want)
			}
		})
	}
}
