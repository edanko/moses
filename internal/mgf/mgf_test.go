package mgf

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"uuid"

	"github.com/edanko/moses/internal/config"
	"github.com/edanko/moses/internal/domain"
	"github.com/edanko/moses/internal/source/txt"
)

type specFunc func(string) (domain.SectionGeometry, error)

func (f specFunc) Spec(profileType string) (domain.SectionGeometry, error) {
	return f(profileType)
}

type spacingFunc func(string, string) (float64, error)

func (f spacingFunc) Spacing(profileType, endcut string) (float64, error) {
	return f(profileType, endcut)
}

func testFormatter(t *testing.T) *Formatter {
	t.Helper()
	cfg, err := config.Parse([]byte("version: 2\n"))
	if err != nil {
		t.Fatalf("config: %v", err)
	}
	return New(cfg, cfg)
}

func TestFormat(t *testing.T) {
	nests, err := txt.ParseFile("../../testdata/txt/056001-05101_RP140X7_A40.txt")
	if err != nil {
		t.Fatalf("parse fixture: %v", err)
	}
	if len(nests) == 0 {
		t.Fatal("no nests in fixture")
	}
	if err := nests[0].Validate(); err != nil {
		t.Fatalf("fixture nest: %v", err)
	}

	output, err := testFormatter(t).Format(nests[0], 1)
	if err != nil {
		t.Fatalf("format: %v", err)
	}
	for _, want := range []string{
		"TYPE_OF_GENERIC_FILE=MOSES_NESTED_PROFILE\r\n",
		"NEST_NAME=056001-05101-RP140X7-A40-1\r\n",
		fmt.Sprintf("NO_OF_PROFS=%d\r\n", len(nests[0].Parts())),
		"PROFILE_DATA\r\n",
		"LEFT_END\r\n",
		"END_OF_RIGHT_END\r\n",
		"END_OF_COMMON_DATA\r\n",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q:\n%s", want, output)
		}
	}
}

func TestFormatResolverErrors(t *testing.T) {
	nests, err := txt.ParseFile("../../testdata/txt/056001-05101_RP140X7_A40.txt")
	if err != nil {
		t.Fatalf("parse fixture: %v", err)
	}
	nest := nests[0]

	t.Run("spec", func(t *testing.T) {
		f := New(spacingFunc(func(string, string) (float64, error) { return 5, nil }),
			specFunc(func(string) (domain.SectionGeometry, error) {
				return domain.SectionGeometry{}, errors.New("no geometry")
			}))
		if _, err := f.Format(nest, 1); err == nil || !strings.Contains(err.Error(), "no geometry") {
			t.Fatalf("error = %v, want geometry error", err)
		}
	})

	t.Run("spacing", func(t *testing.T) {
		cfg, err := config.Parse([]byte("version: 2\n"))
		if err != nil {
			t.Fatalf("config: %v", err)
		}
		f := New(spacingFunc(func(string, string) (float64, error) {
			return 0, errors.New("no spacing")
		}), cfg)
		if _, err := f.Format(nest, 1); err == nil || !strings.Contains(err.Error(), "no spacing") {
			t.Fatalf("error = %v, want spacing error", err)
		}
	})
}

func TestNormText(t *testing.T) {
	designation, err := domain.ParseProfileDesignation("RP160X10")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		spec domain.SectionGeometry
		want string
	}{
		{
			name: "HP",
			spec: domain.SectionGeometry{Shape: domain.ShapeHP, B: 160, C: 28, S: 10, R: 8},
			want: "SHAPE=HP",
		},
		{
			name: "FB",
			spec: domain.SectionGeometry{Shape: domain.ShapeFB, B: 80, S: 8},
			want: "SHAPE=FB",
		},
		{
			name: "L",
			spec: domain.SectionGeometry{Shape: domain.ShapeL, A: 250, B: 90, R1: 19, R2: 9, S: 10, T: 15},
			want: "SHAPE=L",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normText(designation, tt.spec)
			if err != nil {
				t.Fatalf("norm: %v", err)
			}
			if !strings.Contains(got, tt.want) {
				t.Fatalf("norm = %q, want %q", got, tt.want)
			}
		})
	}

	if _, err := normText(designation, domain.SectionGeometry{Shape: domain.ShapeUnknown}); err == nil {
		t.Fatal("expected an error for an unknown shape")
	}
}

func TestParams(t *testing.T) {
	block, scallop := Params("RP160X10", domain.NewEndCut("21", map[string]float64{"v1": 0, "v2": 0}))
	for _, want := range []string{"ENDCUT_TYPE=21R", "V1=90", "V3=90", "END_OF_PARAMS"} {
		if !strings.Contains(block, want) {
			t.Fatalf("params missing %q:\n%s", want, block)
		}
	}
	if scallop != "" {
		t.Fatalf("unexpected scallop %q", scallop)
	}

	if got, _ := Params("RP160X10", domain.EndCut{}); got != "" {
		t.Fatalf("empty end cut = %q, want empty", got)
	}
}

func TestIcutString(t *testing.T) {
	got := IcutString(*domain.NewCutout("i1", map[string]float64{"x": 10, "r": 5}))
	for _, want := range []string{"ICUT_TYPE=1", "X=10", "R=5", "GAMMA=0"} {
		if !strings.Contains(got, want) {
			t.Fatalf("icut missing %q:\n%s", want, got)
		}
	}
}

func TestStraightCut(t *testing.T) {
	parts := []domain.Part{
		mgfPart(t, "1", "fv"),
		mgfPart(t, "2", "ba"),
		mgfPart(t, "3", "b-102"),
		mgfPart(t, "4", "ba"),
	}
	stock := domain.NewStockBar(domain.StockBarParams{
		ID:          uuid.NewV4(),
		Project:     "p",
		Designation: parts[0].Designation(),
		Grade:       parts[0].Grade(),
		Length:      12000,
		Type:        domain.StockTypeBar,
	})
	nest := domain.NewNest(domain.NestParams{StockBar: stock, Parts: parts})

	if got, want := testFormatter(t).StraightCut(*nest), "13"; got != want {
		t.Fatalf("straight cut = %q, want %q", got, want)
	}

	parts = []domain.Part{mgfPart(t, "1", "ba"), mgfPart(t, "2", "sb1"), mgfPart(t, "3", "ba")}
	nest = domain.NewNest(domain.NestParams{StockBar: stock, Parts: parts})
	if got, want := testFormatter(t).StraightCut(*nest), "2"; got != want {
		t.Fatalf("straight cut = %q, want %q", got, want)
	}
}

func mgfPart(t *testing.T, position, rightKind string) domain.Part {
	t.Helper()
	designation, err := domain.ParseProfileDesignation("RP160X10")
	if err != nil {
		t.Fatal(err)
	}
	end := domain.NewEndCut("21", map[string]float64{"v1": 0, "v2": 0})
	return domain.NewPart(domain.PartParams{
		ID:          uuid.NewV4(),
		Project:     "p",
		Section:     "s",
		PositionNo:  position,
		Designation: designation,
		Grade:       "A40",
		Length:      1000,
		Quantity:    1,
		LeftEndCut:  end,
		RightEndCut: domain.NewEndCut(rightKind, nil),
	})
}
