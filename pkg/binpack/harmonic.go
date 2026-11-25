package binpack

import (
	"context"
	"sort"
	"time"
)

type HarmonicKConfig struct {
	K int
}

type HarmonicKAlgorithm[B any, I any] struct {
	name   string
	config HarmonicKConfig
}

func NewHarmonicKAlgorithm[B any, I any](config HarmonicKConfig) *HarmonicKAlgorithm[B, I] {
	if config.K == 0 {
		config.K = 10
	}
	return &HarmonicKAlgorithm[B, I]{
		name:   "harmonic-k",
		config: config,
	}
}

func (h *HarmonicKAlgorithm[B, I]) Name() string {
	return h.name
}

func (h *HarmonicKAlgorithm[B, I]) Pack(ctx context.Context, items []Item[I], bins []*Bin[B, I]) (*PackingResult[B, I], error) {
	start := time.Now()

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

	availableBins := make([]*Bin[B, I], len(bins))
	copy(availableBins, bins)

	k := h.config.K
	binsByClass := make([][]*Bin[B, I], k)
	for i := range binsByClass {
		binsByClass[i] = make([]*Bin[B, I], 0)
	}

	itemsByClass := make([][]Item[I], k)
	for i := range itemsByClass {
		itemsByClass[i] = make([]Item[I], 0)
	}

	maxItemLength := 0.0
	for _, item := range items {
		if item.Length > maxItemLength {
			maxItemLength = item.Length
		}
	}

	for _, item := range items {
		class := h.getItemClass(item.Length, k, maxItemLength)
		if class >= 0 && class < k {
			itemsByClass[class] = append(itemsByClass[class], item)
		}
	}

	packedBins := make([]*Bin[B, I], 0)
	unpackedItems := make([]Item[I], 0)

	for class := range k {
		classItems := itemsByClass[class]
		if len(classItems) == 0 {
			continue
		}

		sort.Slice(classItems, func(i, j int) bool {
			return classItems[i].Length > classItems[j].Length
		})

		for _, item := range classItems {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			placed := false

			for _, bin := range binsByClass[class] {
				if bin.AddItem(item) {
					placed = true
					break
				}
			}

			if !placed {
				if len(availableBins) > 0 {
					newBin := availableBins[0]
					availableBins = availableBins[1:]
					if newBin.AddItem(item) {
						binsByClass[class] = append(binsByClass[class], newBin)
						placed = true
					}
				}
			}

			if !placed {
				unpackedItems = append(unpackedItems, item)
			}
		}
	}

	for class := range k {
		packedBins = append(packedBins, binsByClass[class]...)
	}

	stats := CalculateStatistics(packedBins, unpackedItems, len(items))
	stats.ExecutionTime = time.Since(start).String()
	stats.Algorithm = h.name

	return &PackingResult[B, I]{
		PackedBins:    packedBins,
		UnpackedItems: unpackedItems,
		Statistics:    stats,
	}, nil
}

func (h *HarmonicKAlgorithm[B, I]) getItemClass(length float64, k int, maxLength float64) int {
	if length <= 0 {
		return -1
	}

	normalizedLength := length
	if maxLength > 0 {
		normalizedLength = length / maxLength
	}

	if normalizedLength > 0.5 {
		return 0
	}

	for j := 1; j < k; j++ {
		lower := 1.0 / float64(j+1)
		upper := 1.0 / float64(j)
		if normalizedLength > lower && normalizedLength <= upper {
			return j
		}
		j++
	}

	return k - 1
}

type HarmonicMAlgorithm[B any, I any] struct {
	name string
}

func NewHarmonicMAlgorithm[B any, I any]() *HarmonicMAlgorithm[B, I] {
	return &HarmonicMAlgorithm[B, I]{
		name: "harmonic-m",
	}
}

func (h *HarmonicMAlgorithm[B, I]) Name() string {
	return h.name
}

func (h *HarmonicMAlgorithm[B, I]) Pack(ctx context.Context, items []Item[I], bins []*Bin[B, I]) (*PackingResult[B, I], error) {
	start := time.Now()

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

	largeItems := make([]Item[I], 0)
	smallItems := make([]Item[I], 0)

	for _, item := range items {
		if item.Length > 1.0/3.0 {
			largeItems = append(largeItems, item)
		} else {
			smallItems = append(smallItems, item)
		}
	}

	sort.Slice(largeItems, func(i, j int) bool {
		return largeItems[i].Length > largeItems[j].Length
	})

	availableBins := make([]*Bin[B, I], len(bins))
	copy(availableBins, bins)

	packedBins := make([]*Bin[B, I], 0)
	unpackedItems := make([]Item[I], 0)

	for _, item := range largeItems {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		placed := false
		for _, bin := range packedBins {
			if bin.AddItem(item) {
				placed = true
				break
			}
		}

		if !placed && len(availableBins) > 0 {
			bin := availableBins[0]
			availableBins = availableBins[1:]
			if bin.AddItem(item) {
				packedBins = append(packedBins, bin)
				placed = true
			}
		}

		if !placed {
			unpackedItems = append(unpackedItems, item)
		}
	}

	harmonicK := NewHarmonicKAlgorithm[B, I](HarmonicKConfig{K: 20})
	result, err := harmonicK.Pack(ctx, smallItems, availableBins)
	if err != nil {
		return nil, err
	}

	packedBins = append(packedBins, result.PackedBins...)
	unpackedItems = append(unpackedItems, result.UnpackedItems...)

	stats := CalculateStatistics(packedBins, unpackedItems, len(items))
	stats.ExecutionTime = time.Since(start).String()
	stats.Algorithm = h.name

	return &PackingResult[B, I]{
		PackedBins:    packedBins,
		UnpackedItems: unpackedItems,
		Statistics:    stats,
	}, nil
}

type RefinedFirstFitAlgorithm[B any, I any] struct {
	name string
}

func NewRefinedFirstFitAlgorithm[B any, I any]() *RefinedFirstFitAlgorithm[B, I] {
	return &RefinedFirstFitAlgorithm[B, I]{
		name: "refined-first-fit",
	}
}

func (r *RefinedFirstFitAlgorithm[B, I]) Name() string {
	return r.name
}

func (r *RefinedFirstFitAlgorithm[B, I]) Pack(ctx context.Context, items []Item[I], bins []*Bin[B, I]) (*PackingResult[B, I], error) {
	start := time.Now()

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

	class0 := make([]Item[I], 0)
	class1 := make([]Item[I], 0)
	class2 := make([]Item[I], 0)
	class3 := make([]Item[I], 0)

	for _, item := range items {
		normalizedLen := item.Length
		switch {
		case normalizedLen > 0.5:
			class0 = append(class0, item)
		case normalizedLen > 0.4:
			class1 = append(class1, item)
		case normalizedLen > 1.0/3.0:
			class2 = append(class2, item)
		default:
			class3 = append(class3, item)
		}
	}

	availableBins := make([]*Bin[B, I], len(bins))
	copy(availableBins, bins)

	binsByClass := make([][]*Bin[B, I], 4)
	for i := range binsByClass {
		binsByClass[i] = make([]*Bin[B, I], 0)
	}

	packedBins := make([]*Bin[B, I], 0)
	unpackedItems := make([]Item[I], 0)

	classes := [][]Item[I]{class0, class1, class2, class3}

	for classIdx, classItems := range classes {
		if len(classItems) == 0 {
			continue
		}

		sort.Slice(classItems, func(i, j int) bool {
			return classItems[i].Length > classItems[j].Length
		})

		for _, item := range classItems {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			placed := false

			for _, bin := range binsByClass[classIdx] {
				if bin.AddItem(item) {
					placed = true
					break
				}
			}

			if !placed {
				for _, bin := range binsByClass[classIdx] {
					if bin.AddItem(item) {
						placed = true
						break
					}
				}
			}

			if !placed && len(availableBins) > 0 {
				newBin := availableBins[0]
				availableBins = availableBins[1:]
				if newBin.AddItem(item) {
					binsByClass[classIdx] = append(binsByClass[classIdx], newBin)
					placed = true
				}
			}

			if !placed {
				unpackedItems = append(unpackedItems, item)
			}
		}
	}

	for _, classBins := range binsByClass {
		packedBins = append(packedBins, classBins...)
	}

	stats := CalculateStatistics(packedBins, unpackedItems, len(items))
	stats.ExecutionTime = time.Since(start).String()
	stats.Algorithm = r.name

	return &PackingResult[B, I]{
		PackedBins:    packedBins,
		UnpackedItems: unpackedItems,
		Statistics:    stats,
	}, nil
}
