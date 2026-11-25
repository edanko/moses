package binpack

import (
	"context"
	"time"
)

// FirstFitAlgorithm implements the first-fit bin packing algorithm.
// It places each item in the first bin that can accommodate it.
type FirstFitAlgorithm[B any, I any] struct {
	name string
}

// NewFirstFitAlgorithm creates a new first-fit algorithm
func NewFirstFitAlgorithm[B any, I any]() *FirstFitAlgorithm[B, I] {
	return &FirstFitAlgorithm[B, I]{
		name: "first-fit",
	}
}

// Name returns the algorithm name
func (f *FirstFitAlgorithm[B, I]) Name() string {
	return f.name
}

// Pack performs first-fit bin packing
func (f *FirstFitAlgorithm[B, I]) Pack(ctx context.Context, items []Item[I], bins []*Bin[B, I]) (*PackingResult[B, I], error) {
	start := time.Now()

	// Check context cancellation
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

	// Make a copy of bins that we can mark as used
	availableBins := make([]*Bin[B, I], len(bins))
	copy(availableBins, bins)

	packedBins := make([]*Bin[B, I], 0)
	unpackedItems := make([]Item[I], 0)

	for _, item := range items {
		// Check context cancellation periodically
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		placed := false

		// Try to place in existing packed bins first
		for _, bin := range packedBins {
			if bin.AddItem(item) {
				placed = true
				break
			}
		}

		// If not placed, try to create a new bin from available stock
		if !placed {
			for i, bin := range availableBins {
				if bin.AddItem(item) {
					packedBins = append(packedBins, bin)
					// Remove used bin from available
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
