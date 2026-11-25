package dxf

import (
	"strings"
	"testing"
)

func TestProcessEndCut(t *testing.T) {
	web := []string{"ba", "30", "0", "0", "0", "0"}
	flange := []string{"", "0", "0", "0", "0", "0"}
	got, err := processEndCut(web, flange)
	if err != nil {
		t.Fatalf("process: %v", err)
	}
	if got != "ba angle=30" {
		t.Fatalf("end = %q, want ba angle=30", got)
	}

	web = []string{"bk", "15", "30", "0", "0", "0"}
	flange = []string{"bz", "10", "45", "0", "0", "0"}
	got, err = processEndCut(web, flange)
	if err != nil {
		t.Fatalf("process bz: %v", err)
	}
	if got != "bkbz h=15 angle=30 h2=10 angle2=45" {
		t.Fatalf("end = %q", got)
	}

	if _, err := processEndCut([]string{"unknown"}, nil); err == nil {
		t.Fatal("expected an error for an unknown end cut")
	}
}

func TestParseHeader(t *testing.T) {
	fragment := strings.Join([]string{
		"056001", "05102", "date", "4835", "1", "panel", "A40", "RP160X8",
		"ba", "30", "0", "0", "0", "0",
		"ba", "30", "0", "0", "0", "0",
		"", "0", "0", "0", "0", "0",
		"", "0", "0", "0", "0", "0",
		"direction", "LENGTH=1000",
	}, "\n")

	part, err := parseHeader(fragment)
	if err != nil {
		t.Fatalf("parse header: %v", err)
	}
	if part == nil {
		t.Fatal("profile was skipped")
	}
	if part.posNo != "4835" || part.length != 1000 || part.quantity != 1 {
		t.Fatalf("part = %+v", part)
	}

	profile, err := part.toDomain()
	if err != nil {
		t.Fatalf("to domain: %v", err)
	}
	if profile.DesignationName() != "RP160X8" {
		t.Fatalf("profile type = %s", profile.DesignationName())
	}
	if profile.LeftEndKind() != "ba" || profile.LeftEndCut().Param("angle") != 30 {
		t.Fatalf("left end = %+v", profile.LeftEndCut())
	}
}

func TestParseHeaderSkipsUnsupportedFamily(t *testing.T) {
	fragment := strings.Join([]string{
		"056001", "05102", "date", "4835", "1", "panel", "A40", "AE250X90X10X15",
	}, "\n")
	part, err := parseHeader(fragment)
	if err != nil {
		t.Fatalf("parse header: %v", err)
	}
	if part != nil {
		t.Fatalf("part = %+v, want nil", part)
	}
}

func TestParseRejectsEmptyInput(t *testing.T) {
	if _, err := Parse(strings.NewReader("")); err == nil {
		t.Fatal("expected an error for empty DXF input")
	}
}

func TestInvertFragment(t *testing.T) {
	got := invertFragment("top\nmiddle\nbottom\n")
	if want := "bottom\nmiddle\ntop\n"; got != want {
		t.Fatalf("inverted = %q, want %q", got, want)
	}
}

func TestShiftCutX(t *testing.T) {
	cuts := []string{"i1 x=100 r=10", "i1 x=200 r=5", "no-coordinate"}
	if err := shiftCutX(cuts, 25); err != nil {
		t.Fatalf("shift: %v", err)
	}
	want := []string{"i1 x=125 r=10", "i1 x=225 r=5", "no-coordinate"}
	for i := range want {
		if cuts[i] != want[i] {
			t.Fatalf("cut %d = %q, want %q", i, cuts[i], want[i])
		}
	}
}

func TestRemoveHCCuts(t *testing.T) {
	got := removeHCCuts([]string{"hc x=1", "i1 x=2", "hin x=3"})
	want := []string{"i1 x=2", "hin x=3"}
	if len(got) != len(want) {
		t.Fatalf("cuts = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("cut %d = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestRawPartToDomainErrors(t *testing.T) {
	tests := []struct {
		name string
		part rawPart
	}{
		{name: "bad designation", part: rawPart{dim: "bad"}},
		{name: "bad left end", part: rawPart{dim: "RP160X8", lEnd: "21 broken"}},
		{name: "bad right end", part: rawPart{dim: "RP160X8", lEnd: "21 v1=0", rEnd: "21 broken"}},
		{name: "bad cut", part: rawPart{dim: "RP160X8", lEnd: "21 v1=0", rEnd: "21 v1=0", icuts: []string{"i1 bad"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.part.toDomain(); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}
