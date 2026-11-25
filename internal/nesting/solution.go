package nesting

import (
	"github.com/edanko/moses/internal/domain"
	"github.com/edanko/moses/pkg/binpack"
)

// EvaluateProfileOrder evaluates a profile ordering by simulating packing
// Returns a fitness score (higher is better)
func EvaluateProfileOrder(
	items []binpack.Item[domain.Part],
	bins []*binpack.Bin[domain.StockBar, domain.Part],
) float64 {
	if len(items) == 0 || len(bins) == 0 {
		return 0
	}

	// Clone bins to avoid modifying originals
	availableBins := make([]*binpack.Bin[domain.StockBar, domain.Part], len(bins))
	for i, bin := range bins {
		availableBins[i] = bin.Clone()
	}

	packedBins := make([]*binpack.Bin[domain.StockBar, domain.Part], 0)
	packedCount := 0
	var totalWaste float64

	// Try to pack each item in order
	for _, item := range items {
		placed := false

		// Try to fit in existing packed bins first
		for _, bin := range packedBins {
			if bin.AddItem(item) {
				placed = true
				packedCount++
				break
			}
		}

		// If not placed, try to use a new bin
		if !placed {
			for i, bin := range availableBins {
				if bin.AddItem(item) {
					packedBins = append(packedBins, bin)
					// Remove used bin from available
					availableBins = append(availableBins[:i], availableBins[i+1:]...)
					placed = true
					packedCount++
					break
				}
			}
		}
	}

	// Calculate total waste from used bins
	for _, bin := range packedBins {
		totalWaste += bin.Waste()
	}

	// Calculate fitness: balance packing efficiency and items packed
	utilization := float64(0)
	totalCapacity := float64(0)
	for _, bin := range packedBins {
		totalCapacity += bin.Capacity
	}

	if totalCapacity > 0 {
		utilization = ((totalCapacity - totalWaste) / totalCapacity) * 100
	}

	// Fitness considers: number of items packed, utilization, and bins used
	// Higher fitness is better
	fitness := float64(packedCount)*100 + utilization*10 - float64(len(packedBins))*5

	return fitness
}

// CalculateFitness calculates the fitness of a packing solution
// This is a more detailed fitness function that can be used by metaheuristic algorithms
func CalculateFitness(
	packedBins []*binpack.Bin[domain.StockBar, domain.Part],
	unpackedItems []binpack.Item[domain.Part],
	totalItems int,
) float64 {
	if len(packedBins) == 0 {
		return 0
	}

	packedCount := 0
	var totalWaste, totalCapacity float64

	for _, bin := range packedBins {
		if len(bin.Items) > 0 {
			packedCount += len(bin.Items)
			totalWaste += bin.Waste()
			totalCapacity += bin.Capacity
		}
	}

	utilization := float64(0)
	if totalCapacity > 0 {
		utilization = ((totalCapacity - totalWaste) / totalCapacity) * 100
	}

	// Fitness components:
	// 1. Packed items ratio (0-100)
	packedRatio := float64(packedCount) / float64(totalItems) * 100

	// 2. Material utilization (0-100)
	// 3. Bin efficiency (fewer bins is better)
	binEfficiency := float64(packedCount) / float64(len(packedBins)*10) * 100

	// Combined fitness (weighted sum)
	fitness := packedRatio*0.4 + utilization*0.4 + binEfficiency*0.2

	return fitness
}
