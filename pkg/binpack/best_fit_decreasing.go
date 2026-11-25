package binpack

import (
	"context"
	"time"
)

// BestFitDecreasingAlgorithm implements the best-fit decreasing bin packing algorithm.
// It sorts items by length (descending) and places each item in the bin
// that will have the least remaining space after placement.
type BestFitDecreasingAlgorithm[B any, I any] struct {
	name string
}

// NewBestFitDecreasingAlgorithm creates a new best-fit decreasing algorithm
func NewBestFitDecreasingAlgorithm[B any, I any]() *BestFitDecreasingAlgorithm[B, I] {
	return &BestFitDecreasingAlgorithm[B, I]{
		name: "best-fit-decreasing",
	}
}

// Name returns the algorithm name
func (b *BestFitDecreasingAlgorithm[B, I]) Name() string {
	return b.name
}

// Pack performs best-fit decreasing bin packing
func (b *BestFitDecreasingAlgorithm[B, I]) Pack(ctx context.Context, items []Item[I], bins []*Bin[B, I]) (*PackingResult[B, I], error) {
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

	// Sort items by length in descending order
	sortedItems := make([]Item[I], len(items))
	copy(sortedItems, items)
	SortItemsByLength(sortedItems)

	// Make a copy of bins that we can mark as used
	availableBins := make([]*Bin[B, I], len(bins))
	copy(availableBins, bins)

	packedBins := make([]*Bin[B, I], 0)
	unpackedItems := make([]Item[I], 0)

	for _, item := range sortedItems {
		// Check context cancellation periodically
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		bestFitIndex := -1
		bestFitWaste := float64(-1)

		// Find the best fit among existing packed bins
		for i, bin := range packedBins {
			if bin.CanFit(item) {
				waste := bin.AvailableLength() - item.Length
				if bestFitIndex == -1 || waste < bestFitWaste {
					bestFitIndex = i
					bestFitWaste = waste
				}
			}
		}

		// If found a fit in existing bins, place it there
		if bestFitIndex != -1 {
			packedBins[bestFitIndex].AddItem(item)
			continue
		}

		// Otherwise, find the best fitting new bin from available stock
		bestNewBinIndex := -1
		bestNewBinWaste := float64(-1)

		for i, bin := range availableBins {
			if bin.CanFit(item) {
				waste := bin.AvailableLength() - item.Length
				if bestNewBinIndex == -1 || waste < bestNewBinWaste {
					bestNewBinIndex = i
					bestNewBinWaste = waste
				}
			}
		}

		// If found a suitable new bin, use it
		if bestNewBinIndex != -1 {
			bin := availableBins[bestNewBinIndex]
			bin.AddItem(item)
			packedBins = append(packedBins, bin)
			// Remove used bin from available
			availableBins = append(availableBins[:bestNewBinIndex], availableBins[bestNewBinIndex+1:]...)
		} else {
			// Could not pack this item
			unpackedItems = append(unpackedItems, item)
		}
	}

	stats := CalculateStatistics(packedBins, unpackedItems, len(items))
	stats.ExecutionTime = time.Since(start).String()
	stats.Algorithm = b.name

	return &PackingResult[B, I]{
		PackedBins:    packedBins,
		UnpackedItems: unpackedItems,
		Statistics:    stats,
	}, nil
}
