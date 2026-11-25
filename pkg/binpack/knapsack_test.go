package binpack

import (
	"context"
	"math/rand"
	"testing"
)

func TestKnapsackMaximizesPackedLength(t *testing.T) {
	items := []Item[string]{
		{Original: "six", Length: 6},
		{Original: "five-a", Length: 5},
		{Original: "five-b", Length: 5},
	}
	bins := []*Bin[string, string]{
		NewBin[string, string]("bin", 10),
	}

	algorithm := NewKnapsackAlgorithm[string, string]()
	result, err := algorithm.Pack(context.Background(), items, bins)
	if err != nil {
		t.Fatalf("pack: %v", err)
	}
	if len(result.PackedBins) != 1 {
		t.Fatalf("packed bins = %d, want 1", len(result.PackedBins))
	}
	if len(result.PackedBins[0].Items) != 2 {
		t.Fatalf("packed items = %d, want the two five-length items", len(result.PackedBins[0].Items))
	}
}

func TestKnapsackCeilWeights(t *testing.T) {
	items := []Item[string]{
		{Original: "large", Length: 5.6},
		{Original: "small", Length: 4.2},
	}
	bins := []*Bin[string, string]{
		NewBin[string, string]("bin", 6),
	}

	algorithm := NewKnapsackAlgorithm[string, string]()
	result, err := algorithm.Pack(context.Background(), items, bins)
	if err != nil {
		t.Fatalf("pack: %v", err)
	}
	if len(result.PackedBins) != 1 || len(result.PackedBins[0].Items) != 1 {
		t.Fatalf("unexpected packing: %+v", result.PackedBins)
	}
}

func TestKnapsackSolveIsExact(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	algorithm := NewKnapsackAlgorithm[int, int]()
	workspace := &knapsackWorkspace{}

	for scenario := range 200 {
		n := 1 + rng.Intn(12)
		capacity := 1 + rng.Intn(30)
		items := make([]Item[int], n)
		for i := range items {
			items[i] = Item[int]{Original: i, Length: float64(1 + rng.Intn(15))}
		}

		selected := algorithm.solveKnapsack(items, capacity, workspace)
		var got int
		for _, index := range selected {
			got += int(items[index].Length)
		}

		want := 0
		for mask := 0; mask < 1<<n; mask++ {
			total := 0
			for i := range n {
				if mask&(1<<i) != 0 {
					total += int(items[i].Length)
				}
			}
			if total <= capacity && total > want {
				want = total
			}
		}
		if got != want {
			t.Fatalf("scenario %d: selected length = %d, want %d (capacity %d, items %v, selected %v)",
				scenario, got, want, capacity, items, selected)
		}
	}
}

func TestKnapsackPerBinOptimalIsNotGlobalOptimal(t *testing.T) {
	lengths := []float64{21, 20, 18, 11, 8, 4}

	newBins := func() []*Bin[int, int] {
		bins := make([]*Bin[int, int], len(lengths))
		for i := range bins {
			bins[i] = NewBin[int, int](i, 30)
		}
		return bins
	}
	items := make([]Item[int], len(lengths))
	for i, length := range lengths {
		items[i] = Item[int]{Original: i, Length: length}
	}

	knapsackBins := newBins()
	knapsackResult, err := NewKnapsackAlgorithm[int, int]().Pack(context.Background(), items, knapsackBins)
	if err != nil {
		t.Fatal(err)
	}
	bestFitResult, err := NewBestFitDecreasingAlgorithm[int, int]().Pack(context.Background(), items, newBins())
	if err != nil {
		t.Fatal(err)
	}

	if len(knapsackResult.PackedBins) != 4 {
		t.Fatalf("knapsack bins = %d, want 4 (per-bin optimum is not global optimum)", len(knapsackResult.PackedBins))
	}
	if len(bestFitResult.PackedBins) != 3 {
		t.Fatalf("best-fit bins = %d, want 3", len(bestFitResult.PackedBins))
	}
}
