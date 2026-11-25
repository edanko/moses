package binpack

import (
	"context"
	"sort"
	"time"
)

// AlmostWorstFitAlgorithm implements the almost-worst-fit bin packing algorithm.
// It attempts to place each item in the second most empty bin that can fit it.
// If it doesn't fit in the second emptiest, it tries the emptiest bin.
type AlmostWorstFitAlgorithm[B any, I any] struct {
	name string
}

// NewAlmostWorstFitAlgorithm creates a new almost-worst-fit algorithm
func NewAlmostWorstFitAlgorithm[B any, I any]() *AlmostWorstFitAlgorithm[B, I] {
	return &AlmostWorstFitAlgorithm[B, I]{
		name: "almost-worst-fit",
	}
}

// Name returns the algorithm name
func (a *AlmostWorstFitAlgorithm[B, I]) Name() string {
	return a.name
}

// Pack performs almost-worst-fit bin packing
func (a *AlmostWorstFitAlgorithm[B, I]) Pack(ctx context.Context, items []Item[I], bins []*Bin[B, I]) (*PackingResult[B, I], error) {
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

		placed := false

		if len(packedBins) >= 2 {
			sortedBins := make([]*Bin[B, I], len(packedBins))
			copy(sortedBins, packedBins)
			sort.Slice(sortedBins, func(i, j int) bool {
				return sortedBins[i].UsedLength() < sortedBins[j].UsedLength()
			})

			if sortedBins[1].CanFit(item) {
				sortedBins[1].AddItem(item)
				placed = true
			} else if sortedBins[0].CanFit(item) {
				sortedBins[0].AddItem(item)
				placed = true
			}
		} else if len(packedBins) == 1 {
			if packedBins[0].CanFit(item) {
				packedBins[0].AddItem(item)
				placed = true
			}
		}

		if !placed {
			minSpace := float64(-1)
			bestNewBinIdx := -1
			for i, bin := range availableBins {
				if bin.CanFit(item) {
					space := bin.AvailableLength()
					if bestNewBinIdx == -1 || space < minSpace {
						minSpace = space
						bestNewBinIdx = i
					}
				}
			}

			if bestNewBinIdx != -1 {
				bin := availableBins[bestNewBinIdx]
				bin.AddItem(item)
				packedBins = append(packedBins, bin)
				availableBins = append(availableBins[:bestNewBinIdx], availableBins[bestNewBinIdx+1:]...)
			} else {
				unpackedItems = append(unpackedItems, item)
			}
		}
	}

	stats := CalculateStatistics(packedBins, unpackedItems, len(items))
	stats.ExecutionTime = time.Since(start).String()
	stats.Algorithm = a.name

	return &PackingResult[B, I]{
		PackedBins:    packedBins,
		UnpackedItems: unpackedItems,
		Statistics:    stats,
	}, nil
}
