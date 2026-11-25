package binpack

import (
	"context"
	"time"
)

// NextFitAlgorithm implements the next-fit bin packing algorithm.
// It keeps a single open bin and closes it when the current item doesn't fit,
// then opens a new bin.
type NextFitAlgorithm[B any, I any] struct {
	name string
}

// NewNextFitAlgorithm creates a new next-fit algorithm
func NewNextFitAlgorithm[B any, I any]() *NextFitAlgorithm[B, I] {
	return &NextFitAlgorithm[B, I]{
		name: "next-fit",
	}
}

// Name returns the algorithm name
func (n *NextFitAlgorithm[B, I]) Name() string {
	return n.name
}

// Pack performs next-fit bin packing
func (n *NextFitAlgorithm[B, I]) Pack(ctx context.Context, items []Item[I], bins []*Bin[B, I]) (*PackingResult[B, I], error) {
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

	var currentBin *Bin[B, I]

	for _, item := range items {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if currentBin == nil {
			if len(availableBins) == 0 {
				unpackedItems = append(unpackedItems, item)
				continue
			}
			currentBin = availableBins[0]
			availableBins = availableBins[1:]
			packedBins = append(packedBins, currentBin)
		}

		if currentBin.AddItem(item) {
			continue
		}

		if len(availableBins) == 0 {
			packedBins = packedBins[:len(packedBins)-1]
			unpackedItems = append(unpackedItems, item)
			currentBin = nil
			continue
		}

		currentBin = availableBins[0]
		availableBins = availableBins[1:]
		packedBins = append(packedBins, currentBin)
		if !currentBin.AddItem(item) {
			packedBins = packedBins[:len(packedBins)-1]
			unpackedItems = append(unpackedItems, item)
			currentBin = nil
		}
	}

	stats := CalculateStatistics(packedBins, unpackedItems, len(items))
	stats.ExecutionTime = time.Since(start).String()
	stats.Algorithm = n.name

	return &PackingResult[B, I]{
		PackedBins:    packedBins,
		UnpackedItems: unpackedItems,
		Statistics:    stats,
	}, nil
}
