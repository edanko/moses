package nesting

import (
	"github.com/edanko/moses/internal/domain"
	"github.com/edanko/moses/pkg/binpack"
)

// These functions translate between the domain model and the generic binpack
// package. They are deliberately free functions: translation has no state.

func toBinpackItems(parts []domain.Part) []binpack.Item[domain.Part] {
	items := make([]binpack.Item[domain.Part], len(parts))
	for i, part := range parts {
		items[i] = binpack.Item[domain.Part]{
			Original: part,
			Length:   part.Length(),
		}
	}
	return items
}

func toBinpackBins(stock []domain.StockBar) []*binpack.Bin[domain.StockBar, domain.Part] {
	bins := make([]*binpack.Bin[domain.StockBar, domain.Part], len(stock))
	for i, bar := range stock {
		bins[i] = binpack.NewBin[domain.StockBar, domain.Part](bar, bar.Length())
	}
	return bins
}

func fromBinpackResult(
	result *binpack.PackingResult[domain.StockBar, domain.Part],
) ([]domain.Nest, []domain.Part) {
	nests := make([]domain.Nest, 0, len(result.PackedBins))

	for _, bin := range result.PackedBins {
		if len(bin.Items) == 0 {
			continue
		}

		parts := make([]domain.Part, len(bin.Items))
		for i, item := range bin.Items {
			parts[i] = item.Original
		}

		nest := domain.NewNest(domain.NestParams{
			ID:       bin.Original.ID(),
			StockBar: bin.Original,
			Parts:    parts,
		})
		nests = append(nests, *nest)
	}

	unpacked := make([]domain.Part, len(result.UnpackedItems))
	for i, item := range result.UnpackedItems {
		unpacked[i] = item.Original
	}
	return nests, unpacked
}

func convertStatistics(stats binpack.Statistics, totalParts int) Statistics {
	return Statistics{
		TotalProfiles:       totalParts,
		NestedProfiles:      stats.PackedItems,
		UnnestedProfiles:    stats.UnpackedItems,
		MaterialUtilization: stats.MaterialUtilization,
		TotalWaste:          stats.TotalWaste,
		Algorithm:           stats.Algorithm,
		ExecutionTime:       stats.ExecutionTime,
	}
}
