package binpack

import (
	"context"
	"testing"
)

func testInput() ([]Item[string], []*Bin[string, string]) {
	items := []Item[string]{
		NewItem("a", 6),
		NewItem("b", 4),
		NewItem("c", 5),
		NewItem("d", 3),
	}
	bins := []*Bin[string, string]{
		NewBin[string, string]("bar-1", 10),
		NewBin[string, string]("bar-2", 10),
	}
	return items, bins
}

func TestAlgorithmsPackEveryItemWithinCapacity(t *testing.T) {
	types := []AlgorithmType{
		AlgorithmNextFit,
		AlgorithmFirstFit,
		AlgorithmBestFit,
		AlgorithmWorstFit,
		AlgorithmAlmostWorstFit,
		AlgorithmFirstFitDecreasing,
		AlgorithmBestFitDecreasing,
		AlgorithmNextKFit,
		AlgorithmHarmonicK,
		AlgorithmHarmonicM,
		AlgorithmRefinedFirstFit,
		AlgorithmKnapsack,
	}

	for _, algorithmType := range types {
		t.Run(string(algorithmType), func(t *testing.T) {
			items, bins := testInput()
			algorithm := NewAlgorithm[string, string](algorithmType)
			if algorithm == nil || algorithm.Name() == "" {
				t.Fatal("algorithm has no name")
			}

			result, err := algorithm.Pack(context.Background(), items, bins)
			if err != nil {
				t.Fatalf("pack: %v", err)
			}
			if result == nil {
				t.Fatal("nil result")
			}

			packed := make(map[string]int)
			for _, bin := range result.PackedBins {
				var used float64
				for _, item := range bin.Items {
					packed[item.Original]++
					used += item.Length
				}
				if used > bin.Capacity {
					t.Fatalf("bin %q uses %v > capacity %v", bin.Original, used, bin.Capacity)
				}
			}
			for _, item := range result.UnpackedItems {
				packed[item.Original]++
			}

			if len(packed) != len(items) {
				t.Fatalf("packed %d distinct items, want %d", len(packed), len(items))
			}
			for _, item := range items {
				if packed[item.Original] != 1 {
					t.Fatalf("item %q packed %d times, want once", item.Original, packed[item.Original])
				}
			}
		})
	}
}

func TestBin(t *testing.T) {
	bin := NewBin[string, string]("bar", 10)
	if !bin.CanFit(NewItem("a", 10)) {
		t.Fatal("exact fit should be possible")
	}
	if bin.CanFit(NewItem("b", 11)) {
		t.Fatal("over-capacity item should not fit")
	}
	if !bin.AddItem(NewItem("a", 4)) {
		t.Fatal("item should fit")
	}
	if bin.UsedLength() != 4 || bin.AvailableLength() != 6 || bin.Waste() != 6 {
		t.Fatalf("bin lengths = used %v available %v waste %v", bin.UsedLength(), bin.AvailableLength(), bin.Waste())
	}
	if got := bin.Utilization(); got != 40 {
		t.Fatalf("utilization = %v, want 40", got)
	}
	if bin.AddItem(NewItem("b", 7)) {
		t.Fatal("over-capacity item should be rejected")
	}

	clone := bin.Clone()
	clone.AddItem(NewItem("c", 1))
	if len(bin.Items) != 1 {
		t.Fatal("clone shares the items slice")
	}

	bin.Reset()
	if len(bin.Items) != 0 || bin.UsedLength() != 0 {
		t.Fatalf("reset bin = %+v", bin)
	}
}
