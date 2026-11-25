package domain

import "testing"

func TestParseProfileDesignation(t *testing.T) {
	tests := []struct {
		input       string
		designation string
		family      string
		height      float64
		thick       float64
	}{
		{input: "RP160X10", designation: "RP160X10", family: "RP", height: 160, thick: 10},
		{input: "PP160*8.0", designation: "RP160X8", family: "RP", height: 160, thick: 8},
		{input: "16a", designation: "RP160X8", family: "RP", height: 160, thick: 8},
		{input: "fb80x8", designation: "FB80X8", family: "FB", height: 80, thick: 8},
		{input: "AE250x90x10x15", designation: "AE250X90X10X15", family: "AE", height: 250, thick: 10},
		{input: "ae40x4", designation: "AE40X4", family: "AE", height: 40, thick: 4},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got, err := ParseProfileDesignation(test.input)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			if got.String() != test.designation || got.Family() != test.family {
				t.Fatalf("designation/family = %q/%q, want %q/%q", got.String(), got.Family(), test.designation, test.family)
			}
			if got.Height() != test.height || got.Thickness() != test.thick {
				t.Fatalf("height/thickness = %v/%v", got.Height(), got.Thickness())
			}
		})
	}
}

func TestParseProfileDesignationErrors(t *testing.T) {
	for _, input := range []string{"", "RP", "XX160X10", "RP160", "AE250"} {
		if _, err := ParseProfileDesignation(input); err == nil {
			t.Fatalf("expected an error for %q", input)
		}
	}
}

func TestParseEndCut(t *testing.T) {
	end, err := ParseEndCut("21 V1=0 v2=50 snipe=10")
	if err != nil {
		t.Fatalf("parse end: %v", err)
	}
	if end.Kind() != "21" || end.Param("v1") != 0 || end.Param("v2") != 50 || end.Param("snipe") != 10 {
		t.Fatalf("unexpected end: %+v", end)
	}
	if end.String() != "21 V1=0 v2=50 snipe=10" {
		t.Fatalf("raw = %q", end.String())
	}
	if _, err := ParseEndCut("21 bad"); err == nil {
		t.Fatal("expected an error for a malformed parameter")
	}
}

func TestParseCutout(t *testing.T) {
	hole, err := ParseCutout("i1 x=180 r=30")
	if err != nil {
		t.Fatalf("parse hole: %v", err)
	}
	if hole.Kind() != "i1" || hole.Param("x") != 180 || hole.Param("r") != 30 {
		t.Fatalf("unexpected hole: %+v", hole)
	}
	if hole.String() != "i1 x=180 r=30" {
		t.Fatalf("raw = %q", hole.String())
	}
}

func TestSectionGeometryValidate(t *testing.T) {
	valid := SectionGeometry{Shape: ShapeHP, B: 160, C: 28, S: 8, R: 7}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid spec: %v", err)
	}
	invalid := SectionGeometry{Shape: ShapeHP, B: 160}
	if err := invalid.Validate(); err == nil {
		t.Fatal("expected an error for an incomplete HP spec")
	}
	if err := (SectionGeometry{Shape: ShapeFB, B: 80, S: 8}).Validate(); err != nil {
		t.Fatalf("valid FB spec: %v", err)
	}
	if err := (SectionGeometry{Shape: ShapeL, A: 250, B: 90, R1: 19, R2: 9.5, S: 10, T: 15}).Validate(); err != nil {
		t.Fatalf("valid L spec: %v", err)
	}
}
