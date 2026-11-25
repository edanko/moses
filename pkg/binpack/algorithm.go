package binpack

import (
	"context"
	"slices"
)

// Algorithm defines the interface for bin packing algorithms
type Algorithm[B any, I any] interface {
	// Name returns the algorithm name
	Name() string

	// Pack performs the bin packing operation
	// items: items to be packed
	// bins: available bins (stock)
	// Returns a PackingResult with packed bins and any unpacked items
	Pack(ctx context.Context, items []Item[I], bins []*Bin[B, I]) (*PackingResult[B, I], error)
}

// AlgorithmType represents the type of bin packing algorithm
type AlgorithmType string

const (
	AlgorithmNextFit            AlgorithmType = "next-fit"
	AlgorithmFirstFit           AlgorithmType = "first-fit"
	AlgorithmBestFit            AlgorithmType = "best-fit"
	AlgorithmWorstFit           AlgorithmType = "worst-fit"
	AlgorithmAlmostWorstFit     AlgorithmType = "almost-worst-fit"
	AlgorithmFirstFitDecreasing AlgorithmType = "first-fit-decreasing"
	AlgorithmBestFitDecreasing  AlgorithmType = "best-fit-decreasing"
	AlgorithmNextKFit           AlgorithmType = "next-k-fit"
	AlgorithmHarmonicK          AlgorithmType = "harmonic-k"
	AlgorithmHarmonicM          AlgorithmType = "harmonic-m"
	AlgorithmRefinedFirstFit    AlgorithmType = "refined-first-fit"
	AlgorithmKnapsack           AlgorithmType = "knapsack"
)

// NewAlgorithm creates a new algorithm instance by type
func NewAlgorithm[B any, I any](algoType AlgorithmType, opts ...any) Algorithm[B, I] {
	switch algoType {
	case AlgorithmNextFit:
		return NewNextFitAlgorithm[B, I]()
	case AlgorithmFirstFit:
		return NewFirstFitAlgorithm[B, I]()
	case AlgorithmBestFit:
		return NewBestFitAlgorithm[B, I]()
	case AlgorithmWorstFit:
		return NewWorstFitAlgorithm[B, I]()
	case AlgorithmAlmostWorstFit:
		return NewAlmostWorstFitAlgorithm[B, I]()
	case AlgorithmFirstFitDecreasing:
		return NewFirstFitDecreasingAlgorithm[B, I]()
	case AlgorithmBestFitDecreasing:
		return NewBestFitDecreasingAlgorithm[B, I]()
	case AlgorithmNextKFit:
		config := NextKFitConfig{K: 2}
		if len(opts) > 0 {
			if k, ok := opts[0].(int); ok {
				config.K = k
			}
		}
		return NewNextKFitAlgorithm[B, I](config)
	case AlgorithmHarmonicK:
		config := HarmonicKConfig{K: 10}
		if len(opts) > 0 {
			if k, ok := opts[0].(int); ok {
				config.K = k
			}
		}
		return NewHarmonicKAlgorithm[B, I](config)
	case AlgorithmHarmonicM:
		return NewHarmonicMAlgorithm[B, I]()
	case AlgorithmRefinedFirstFit:
		return NewRefinedFirstFitAlgorithm[B, I]()
	case AlgorithmKnapsack:
		return NewKnapsackAlgorithm[B, I]()
	default:
		return NewFirstFitDecreasingAlgorithm[B, I]()
	}
}

// LengthFunc is a function that returns the length for a given item
type LengthFunc[T any] func(item T) float64

// BinCapacityFunc is a function that returns the capacity for a given bin
type BinCapacityFunc[B any] func(bin B) float64

// PackRequest is a generic request for packing operations
type PackRequest[B any, I any] struct {
	Items           []I
	Bins            []B
	ItemLengthFunc  LengthFunc[I]
	BinCapacityFunc BinCapacityFunc[B]
}

// ToItems converts typed items to binpack items
func ToItems[T any](
	items []T,
	lengthFunc LengthFunc[T],
) []Item[T] {
	result := make([]Item[T], len(items))
	for i, item := range items {
		result[i] = Item[T]{
			Original: item,
			Length:   lengthFunc(item),
		}
	}
	return result
}

// ToBins converts typed bins to binpack bins
func ToBins[B any, I any](
	bins []B,
	capacityFunc BinCapacityFunc[B],
) []*Bin[B, I] {
	result := make([]*Bin[B, I], len(bins))
	for i, bin := range bins {
		result[i] = NewBin[B, I](bin, capacityFunc(bin))
	}
	return result
}

// SortItemsByLength sorts items by length in descending order
func SortItemsByLength[T any](items []Item[T]) {
	slices.SortFunc(items, func(a, b Item[T]) int {
		if a.Length > b.Length {
			return -1
		} else if a.Length < b.Length {
			return 1
		}
		return 0
	})
}

// SortItemsByLengthAsc sorts items by length in ascending order
func SortItemsByLengthAsc[T any](items []Item[T]) {
	slices.SortFunc(items, func(a, b Item[T]) int {
		if a.Length < b.Length {
			return -1
		} else if a.Length > b.Length {
			return 1
		}
		return 0
	})
}

// SortBinsByCapacity sorts bins by capacity in ascending order
func SortBinsByCapacity[B any, I any](bins []*Bin[B, I]) {
	slices.SortFunc(bins, func(a, b *Bin[B, I]) int {
		if a.Capacity < b.Capacity {
			return -1
		} else if a.Capacity > b.Capacity {
			return 1
		}
		return 0
	})
}

// CalculateStatistics computes statistics for a packing result
func CalculateStatistics[B any, I any](packedBins []*Bin[B, I], unpackedItems []Item[I], totalItems int) Statistics {
	var totalBinsUsed int
	var totalWaste float64
	var totalCapacity float64
	var packedItems int

	for _, bin := range packedBins {
		if len(bin.Items) > 0 {
			totalBinsUsed++
			totalWaste += bin.Waste()
			totalCapacity += bin.Capacity
			packedItems += len(bin.Items)
		}
	}

	utilization := float64(0)
	if totalCapacity > 0 {
		utilization = ((totalCapacity - totalWaste) / totalCapacity) * 100
	}

	return Statistics{
		TotalItems:          totalItems,
		PackedItems:         packedItems,
		UnpackedItems:       len(unpackedItems),
		TotalBinsUsed:       totalBinsUsed,
		MaterialUtilization: utilization,
		TotalWaste:          totalWaste,
	}
}

// DefaultFitnessFunc returns a default fitness function that maximizes bin utilization
// and minimizes the number of bins used
func DefaultFitnessFunc[B any, I any](items []Item[I], bins []*Bin[B, I]) float64 {
	base := NewFirstFitDecreasingAlgorithm[B, I]()
	ctx := context.Background()
	result, err := base.Pack(ctx, items, bins)
	if err != nil {
		return 0
	}

	fitness := 0.0
	for _, bin := range result.PackedBins {
		fitness += bin.Utilization()
	}

	fitness -= float64(len(result.PackedBins)) * 10

	return fitness
}

// AvailableAlgorithms returns a list of all available algorithm types
func AvailableAlgorithms() []AlgorithmType {
	return []AlgorithmType{
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
}
