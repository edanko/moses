package nesting

import (
	"context"
	"testing"
	"uuid"

	"github.com/edanko/moses/internal/domain"
)

func testProfile(position string, length float64, quantity int) domain.Part {
	profileType, err := domain.ParseProfileDesignation("RP160X8")
	if err != nil {
		panic(err)
	}
	end, err := domain.ParseEndCut("21 v1=0 v2=0")
	if err != nil {
		panic(err)
	}
	return domain.NewPart(domain.PartParams{
		ID:          uuid.NewV4(),
		Project:     "project",
		Section:     "section",
		PositionNo:  position,
		Designation: profileType,
		Grade:       "A40",
		Length:      length,
		Quantity:    quantity,
		LeftEndCut:  end,
		RightEndCut: end,
	})
}

func testOptions() PackOptions {
	return PackOptions{
		Spacing:     func(_, _ string) (float64, error) { return 0, nil },
		StockLength: func(string) (float64, error) { return 10000, nil },
		Algorithm:   "knapsack",
	}
}

func TestPackPlacesEveryPart(t *testing.T) {
	profiles := []domain.Part{
		testProfile("1", 6000, 1),
		testProfile("2", 5000, 1),
		testProfile("3", 2000, 1),
		testProfile("4", 1000, 1),
	}

	result, err := Pack(context.Background(), profiles, nil, testOptions())
	if err != nil {
		t.Fatalf("pack: %v", err)
	}
	if len(result.Nests) != 2 {
		t.Fatalf("nests = %d, want 2", len(result.Nests))
	}

	placed := 0
	for _, nest := range result.Nests {
		placed += len(nest.Parts())
		var used float64
		for _, profile := range nest.Parts() {
			used += profile.Length()
		}
		if used > nest.StockBar().Length() {
			t.Fatalf("nest used %v > capacity %v", used, nest.StockBar().Length())
		}
	}
	if placed != 4 {
		t.Fatalf("placed = %d, want 4", placed)
	}
	if result.Utilization <= 0 {
		t.Fatal("expected positive utilization")
	}
}

func TestPackExpandsQuantity(t *testing.T) {
	profiles := []domain.Part{
		testProfile("1", 6000, 2),
		testProfile("2", 4000, 1),
	}
	result, err := Pack(context.Background(), profiles, nil, testOptions())
	if err != nil {
		t.Fatalf("pack: %v", err)
	}
	placed := 0
	for _, nest := range result.Nests {
		placed += len(nest.Parts())
	}
	if placed != 3 {
		t.Fatalf("placed = %d, want 3", placed)
	}
}

func TestPackUsesRemnant(t *testing.T) {
	profileType, err := domain.ParseProfileDesignation("RP160X8")
	if err != nil {
		t.Fatal(err)
	}
	remnants := []domain.Remnant{{
		ID:          7,
		Length:      4000,
		Project:     "project",
		Source:      "section",
		Designation: profileType,
		Grade:       "A40",
		Marking:     "rem-1",
	}}

	profiles := []domain.Part{
		testProfile("1", 5000, 1),
		testProfile("2", 5000, 1),
		testProfile("3", 4000, 1),
	}
	result, err := Pack(context.Background(), profiles, remnants, testOptions())
	if err != nil {
		t.Fatalf("pack: %v", err)
	}

	consumed := 0
	for _, nest := range result.Nests {
		if nest.ConsumedRemnant() != nil {
			consumed++
			if nest.ConsumedRemnant().ID != 7 {
				t.Fatalf("consumed id = %d", nest.ConsumedRemnant().ID)
			}
		}
	}
	if consumed != 1 {
		t.Fatalf("consumed remnants = %d, want 1", consumed)
	}
}

func TestPackOverlength(t *testing.T) {
	result, err := Pack(context.Background(), []domain.Part{testProfile("1", 12000, 1)}, nil, testOptions())
	if err != nil {
		t.Fatalf("pack: %v", err)
	}
	if len(result.Nests) != 1 {
		t.Fatalf("nests = %d, want 1", len(result.Nests))
	}
	if result.Nests[0].StockBar().Type() != domain.StockTypeOverlength {
		t.Fatalf("kind = %q, want overlength", result.Nests[0].StockBar().Type())
	}
}

func TestPackDeterministic(t *testing.T) {
	profiles := []domain.Part{
		testProfile("1", 3000, 1),
		testProfile("2", 3000, 1),
		testProfile("3", 2000, 1),
		testProfile("4", 2000, 1),
	}

	first, err := Pack(context.Background(), profiles, nil, testOptions())
	if err != nil {
		t.Fatal(err)
	}
	second, err := Pack(context.Background(), profiles, nil, testOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(first.Nests) != len(second.Nests) {
		t.Fatalf("nest counts differ: %d vs %d", len(first.Nests), len(second.Nests))
	}
	for i := range first.Nests {
		left := first.Nests[i].Parts()
		right := second.Nests[i].Parts()
		if len(left) != len(right) {
			t.Fatalf("nest %d profile counts differ", i)
		}
		for j := range left {
			if left[j].PositionNo() != right[j].PositionNo() {
				t.Fatalf("nest %d profile %d differs: %s vs %s", i, j, left[j].PositionNo(), right[j].PositionNo())
			}
		}
	}
}

func TestPackBestFitPlacesEveryPart(t *testing.T) {
	options := testOptions()
	options.Algorithm = "best-fit-decreasing"
	profiles := []domain.Part{
		testProfile("1", 6000, 1),
		testProfile("2", 5000, 1),
		testProfile("3", 2000, 1),
		testProfile("4", 1000, 1),
	}

	result, err := Pack(context.Background(), profiles, nil, options)
	if err != nil {
		t.Fatalf("pack: %v", err)
	}
	placed := 0
	for _, nest := range result.Nests {
		placed += len(nest.Parts())
	}
	if placed != len(profiles) {
		t.Fatalf("placed = %d, want %d", placed, len(profiles))
	}
}
