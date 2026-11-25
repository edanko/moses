package binpack

import (
	"context"
	"math"
	"time"
)

// KnapsackAlgorithm implements the 0/1 knapsack dynamic programming algorithm
// for optimal packing of items into each bin.
type KnapsackAlgorithm[B any, I any] struct {
	name string
}

// NewKnapsackAlgorithm creates a new knapsack algorithm
func NewKnapsackAlgorithm[B any, I any]() *KnapsackAlgorithm[B, I] {
	return &KnapsackAlgorithm[B, I]{
		name: "knapsack",
	}
}

// Name returns the algorithm name
func (k *KnapsackAlgorithm[B, I]) Name() string {
	return k.name
}

// Pack performs knapsack-based bin packing
func (k *KnapsackAlgorithm[B, I]) Pack(ctx context.Context, items []Item[I], bins []*Bin[B, I]) (*PackingResult[B, I], error) {
	start := time.Now()

	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if len(bins) == 0 || len(items) == 0 {
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

	// Sort bins by capacity (ascending) to use smallest suitable bin first
	sortedBins := make([]*Bin[B, I], len(bins))
	copy(sortedBins, bins)
	SortBinsByCapacity(sortedBins)

	// Sort items by length (descending) for better knapsack results
	sortedItems := make([]Item[I], len(items))
	copy(sortedItems, items)
	SortItemsByLength(sortedItems)

	// Track which items have been packed (by index in sortedItems)
	usedIndices := make([]bool, len(sortedItems))

	maxCapacity := 0
	for _, bin := range sortedBins {
		if capacity := int(math.Ceil(bin.Capacity)); capacity > maxCapacity {
			maxCapacity = capacity
		}
	}
	workspace := &knapsackWorkspace{
		lengths:  make([]int, 0, len(sortedItems)),
		dp:       make([]int, maxCapacity+1),
		take:     make([]uint64, ((len(sortedItems)+1)*(maxCapacity+1)+63)/64),
		selected: make([]int, 0, len(sortedItems)),
	}

	packedBins := make([]*Bin[B, I], 0, len(sortedBins))
	remainingItems := make([]Item[I], 0, len(sortedItems))
	remainingIndices := make([]int, 0, len(sortedItems))

	for _, bin := range sortedBins {
		// Check context cancellation periodically
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// Build list of remaining items (not yet packed)
		remainingItems = remainingItems[:0]
		remainingIndices = remainingIndices[:0]
		for i, item := range sortedItems {
			if !usedIndices[i] {
				remainingItems = append(remainingItems, item)
				remainingIndices = append(remainingIndices, i)
			}
		}

		if len(remainingItems) == 0 {
			break
		}

		// Find the subset that packs the most length into this bin.
		selectedIndices := k.solveKnapsack(remainingItems, int(math.Ceil(bin.Capacity)), workspace)

		added := 0
		for _, idx := range selectedIndices {
			item := remainingItems[idx]
			if !bin.AddItem(item) {
				continue
			}
			usedIndices[remainingIndices[idx]] = true
			added++
		}
		if added > 0 {
			packedBins = append(packedBins, bin)
		}
	}

	// Build remaining items list
	unpackedItems := make([]Item[I], 0, len(sortedItems))
	for i, used := range usedIndices {
		if !used {
			unpackedItems = append(unpackedItems, sortedItems[i])
		}
	}

	stats := CalculateStatistics(packedBins, unpackedItems, len(items))
	stats.ExecutionTime = time.Since(start).String()
	stats.Algorithm = k.name

	return &PackingResult[B, I]{
		PackedBins:    packedBins,
		UnpackedItems: unpackedItems,
		Statistics:    stats,
	}, nil
}

type knapsackWorkspace struct {
	lengths  []int
	dp       []int
	take     []uint64
	selected []int
}

// solveKnapsack solves the 0/1 knapsack problem for items fitting into capacity.
// It returns indices of selected items in the input slice.
//
// The DP keeps one row per item as bits in a decision bitset, while the value
// table itself is a single rolling row. The caller owns the workspace so the
// buffers are reused across bins instead of being reallocated per bin.
func (k *KnapsackAlgorithm[B, I]) solveKnapsack(items []Item[I], capacity int, workspace *knapsackWorkspace) []int {
	n := len(items)
	if n == 0 || capacity <= 0 {
		return nil
	}

	width := capacity + 1
	if width > len(workspace.dp) {
		workspace.dp = make([]int, width)
	}
	dp := workspace.dp[:width]
	clear(dp)

	words := ((n+1)*width + 63) / 64
	if words > len(workspace.take) {
		workspace.take = make([]uint64, words)
	}
	take := workspace.take[:words]
	clear(take)

	if n > cap(workspace.lengths) {
		workspace.lengths = make([]int, 0, n)
	}
	lengths := workspace.lengths[:n]
	for i, item := range items {
		lengths[i] = int(math.Ceil(item.Length))
	}

	for i := 1; i <= n; i++ {
		length := lengths[i-1]
		if length > capacity {
			continue
		}
		row := i * width
		for w := capacity; w >= length; w-- {
			candidate := dp[w-length] + length
			if candidate > dp[w] {
				dp[w] = candidate
				bit := row + w
				take[bit>>6] |= 1 << uint(bit&63)
			}
		}
	}

	selected := workspace.selected[:0]
	w := capacity
	for i := n; i > 0 && w > 0; i-- {
		bit := i*width + w
		if (take[bit>>6]>>uint(bit&63))&1 == 0 {
			continue
		}
		selected = append(selected, i-1)
		w -= lengths[i-1]
	}
	workspace.selected = selected

	return selected
}
