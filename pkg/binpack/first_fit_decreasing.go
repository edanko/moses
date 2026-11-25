package binpack

import (
	"context"
	"time"
)

// FirstFitDecreasingAlgorithm implements the first-fit decreasing bin packing algorithm.
// It sorts items by length (descending) and places each item in the first bin
// that can accommodate it.
type FirstFitDecreasingAlgorithm[B any, I any] struct {
	name string
}

// NewFirstFitDecreasingAlgorithm creates a new first-fit decreasing algorithm
func NewFirstFitDecreasingAlgorithm[B any, I any]() *FirstFitDecreasingAlgorithm[B, I] {
	return &FirstFitDecreasingAlgorithm[B, I]{
		name: "first-fit-decreasing",
	}
}

// Name returns the algorithm name
func (f *FirstFitDecreasingAlgorithm[B, I]) Name() string {
	return f.name
}

// Pack performs first-fit decreasing bin packing
func (f *FirstFitDecreasingAlgorithm[B, I]) Pack(ctx context.Context, items []Item[I], bins []*Bin[B, I]) (*PackingResult[B, I], error) {
	start := time.Now()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if len(bins) == 0 {
		return &PackingResult[B, I]{
			PackedBins:    make([]*Bin[B, I], 0),
			UnpackedItems: items,
			Statistics: Statistics{
				TotalItems:    len(items),
				UnpackedItems: len(items),
				ExecutionTime: time.Since(start).String(),
			},
		}, nil
	}

	sortedItems := make([]Item[I], len(items))
	copy(sortedItems, items)
	SortItemsByLength(sortedItems)

	availableBins := make([]*Bin[B, I], len(bins))
	copy(availableBins, bins)

	packedBins := make([]*Bin[B, I], 0)
	unpackedItems := make([]Item[I], 0)

	for _, item := range sortedItems {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		placed := false

		for _, bin := range packedBins {
			if bin.AddItem(item) {
				placed = true
				break
			}
		}

		if !placed {
			for i, bin := range availableBins {
				if bin.AddItem(item) {
					packedBins = append(packedBins, bin)
					availableBins = append(availableBins[:i], availableBins[i+1:]...)
					placed = true
					break
				}
			}
		}

		if !placed {
			unpackedItems = append(unpackedItems, item)
		}
	}

	stats := CalculateStatistics(packedBins, unpackedItems, len(items))
	stats.ExecutionTime = time.Since(start).String()
	stats.Algorithm = f.name

	return &PackingResult[B, I]{
		PackedBins:    packedBins,
		UnpackedItems: unpackedItems,
		Statistics:    stats,
	}, nil
}
