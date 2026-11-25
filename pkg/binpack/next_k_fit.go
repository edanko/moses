package binpack

import (
	"context"
	"slices"
	"time"
)

type NextKFitConfig struct {
	K int
}

type NextKFitAlgorithm[B any, I any] struct {
	name   string
	config NextKFitConfig
}

func NewNextKFitAlgorithm[B any, I any](config NextKFitConfig) *NextKFitAlgorithm[B, I] {
	if config.K <= 0 {
		config.K = 2
	}
	return &NextKFitAlgorithm[B, I]{
		name:   "next-k-fit",
		config: config,
	}
}

func (n *NextKFitAlgorithm[B, I]) Name() string {
	return n.name
}

func (n *NextKFitAlgorithm[B, I]) Pack(ctx context.Context, items []Item[I], bins []*Bin[B, I]) (*PackingResult[B, I], error) {
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

	k := n.config.K
	openBins := make([]*Bin[B, I], 0, k)
	packedBins := make([]*Bin[B, I], 0)
	unpackedItems := make([]Item[I], 0)

	for _, item := range items {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		placed := false

		for _, bin := range openBins {
			if bin.AddItem(item) {
				placed = true
				break
			}
		}

		if !placed {
			closed := false
			for _, bin := range packedBins {
				isOpen := slices.Contains(openBins, bin)
				if !isOpen && bin.AddItem(item) {
					openBins = append(openBins, bin)
					placed = true
					closed = true
					break
				}
			}

			if !closed && len(openBins) >= k && !placed {
				closedBin := openBins[0]
				packedBins = append(packedBins, closedBin)
				openBins = openBins[1:]

				for _, bin := range openBins {
					if bin.AddItem(item) {
						placed = true
						break
					}
				}

				if !placed && len(availableBins) > 0 {
					newBin := availableBins[0]
					availableBins = availableBins[1:]
					if newBin.AddItem(item) {
						openBins = append(openBins, newBin)
						placed = true
					}
				}
			}

			if !placed && len(availableBins) > 0 {
				newBin := availableBins[0]
				availableBins = availableBins[1:]
				if newBin.AddItem(item) {
					openBins = append(openBins, newBin)
					placed = true
				}
			}
		}

		if !placed {
			unpackedItems = append(unpackedItems, item)
		}
	}

	packedBins = append(packedBins, openBins...)

	stats := CalculateStatistics(packedBins, unpackedItems, len(items))
	stats.ExecutionTime = time.Since(start).String()
	stats.Algorithm = n.name

	return &PackingResult[B, I]{
		PackedBins:    packedBins,
		UnpackedItems: unpackedItems,
		Statistics:    stats,
	}, nil
}
