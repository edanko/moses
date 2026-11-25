package binpack

import (
	"context"
	"time"
)

// WorstFitAlgorithm implements the worst-fit bin packing algorithm.
// It places each item in the bin with the maximum remaining space (least filled).
type WorstFitAlgorithm[B any, I any] struct {
	name string
}

// NewWorstFitAlgorithm creates a new worst-fit algorithm
func NewWorstFitAlgorithm[B any, I any]() *WorstFitAlgorithm[B, I] {
	return &WorstFitAlgorithm[B, I]{
		name: "worst-fit",
	}
}

// Name returns the algorithm name
func (w *WorstFitAlgorithm[B, I]) Name() string {
	return w.name
}

// Pack performs worst-fit bin packing
func (w *WorstFitAlgorithm[B, I]) Pack(ctx context.Context, items []Item[I], bins []*Bin[B, I]) (*PackingResult[B, I], error) {
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

		worstFitIndex := -1
		worstFitSpace := float64(-1)

		for i, bin := range packedBins {
			if bin.CanFit(item) {
				space := bin.AvailableLength()
				if worstFitIndex == -1 || space > worstFitSpace {
					worstFitIndex = i
					worstFitSpace = space
				}
			}
		}

		if worstFitIndex != -1 {
			packedBins[worstFitIndex].AddItem(item)
			continue
		}

		worstNewBinIndex := -1
		worstNewBinSpace := float64(-1)

		for i, bin := range availableBins {
			if bin.CanFit(item) {
				space := bin.AvailableLength()
				if worstNewBinIndex == -1 || space > worstNewBinSpace {
					worstNewBinIndex = i
					worstNewBinSpace = space
				}
			}
		}

		if worstNewBinIndex != -1 {
			bin := availableBins[worstNewBinIndex]
			bin.AddItem(item)
			packedBins = append(packedBins, bin)
			availableBins = append(availableBins[:worstNewBinIndex], availableBins[worstNewBinIndex+1:]...)
		} else {
			unpackedItems = append(unpackedItems, item)
		}
	}

	stats := CalculateStatistics(packedBins, unpackedItems, len(items))
	stats.ExecutionTime = time.Since(start).String()
	stats.Algorithm = w.name

	return &PackingResult[B, I]{
		PackedBins:    packedBins,
		UnpackedItems: unpackedItems,
		Statistics:    stats,
	}, nil
}
