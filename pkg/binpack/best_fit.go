package binpack

import (
	"context"
	"time"
)

// BestFitAlgorithm implements the best-fit bin packing algorithm.
// It places each item in the bin with the maximum load in which it fits.
type BestFitAlgorithm[B any, I any] struct {
	name string
}

// NewBestFitAlgorithm creates a new best-fit algorithm
func NewBestFitAlgorithm[B any, I any]() *BestFitAlgorithm[B, I] {
	return &BestFitAlgorithm[B, I]{
		name: "best-fit",
	}
}

// Name returns the algorithm name
func (b *BestFitAlgorithm[B, I]) Name() string {
	return b.name
}

// Pack performs best-fit bin packing
func (b *BestFitAlgorithm[B, I]) Pack(ctx context.Context, items []Item[I], bins []*Bin[B, I]) (*PackingResult[B, I], error) {
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

	availableBins := make([]*Bin[B, I], len(bins))
	copy(availableBins, bins)

	packedBins := make([]*Bin[B, I], 0)
	unpackedItems := make([]Item[I], 0)

	for _, item := range items {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		bestFitIndex := -1
		bestFitWaste := float64(-1)

		for i, bin := range packedBins {
			if bin.CanFit(item) {
				waste := bin.AvailableLength() - item.Length
				if bestFitIndex == -1 || waste < bestFitWaste {
					bestFitIndex = i
					bestFitWaste = waste
				}
			}
		}

		if bestFitIndex != -1 {
			packedBins[bestFitIndex].AddItem(item)
			continue
		}

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

		if bestNewBinIndex != -1 {
			bin := availableBins[bestNewBinIndex]
			bin.AddItem(item)
			packedBins = append(packedBins, bin)
			availableBins = append(availableBins[:bestNewBinIndex], availableBins[bestNewBinIndex+1:]...)
		} else {
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
