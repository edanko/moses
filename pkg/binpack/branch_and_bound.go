package binpack

import (
	"context"
	"math"
	"time"
)

type BranchAndBoundConfig struct {
	TimeLimit    time.Duration
	MaxSolutions int
}

type BranchAndBound[B any, I any] struct {
	name    string
	config  BranchAndBoundConfig
	base    Algorithm[B, I]
	fitness FitnessFunc[B, I]
}

func NewBranchAndBound[B any, I any](config BranchAndBoundConfig, base Algorithm[B, I], fitness FitnessFunc[B, I]) *BranchAndBound[B, I] {
	if config.TimeLimit == 0 {
		config.TimeLimit = 30 * time.Second
	}
	if config.MaxSolutions == 0 {
		config.MaxSolutions = 10000
	}
	return &BranchAndBound[B, I]{
		name:    "branch-and-bound",
		config:  config,
		base:    base,
		fitness: fitness,
	}
}

func (b *BranchAndBound[B, I]) Name() string {
	return b.name
}

func (b *BranchAndBound[B, I]) Pack(ctx context.Context, items []Item[I], bins []*Bin[B, I]) (*PackingResult[B, I], error) {
	start := time.Now()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if len(items) == 0 || len(bins) == 0 {
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

	bestSolution := make([]int, len(items))
	for i := range bestSolution {
		bestSolution[i] = i
	}

	bestFitness := math.Inf(-1)
	solutionsFound := 0

	b.branchAndBound(
		ctx,
		make([]int, 0, len(items)),
		make([]bool, len(items)),
		items,
		bins,
		start,
		&bestSolution,
		&bestFitness,
		&solutionsFound,
	)

	bestItems := make([]Item[I], len(bestSolution))
	for i, idx := range bestSolution {
		bestItems[i] = items[idx]
	}

	result, err := b.base.Pack(ctx, bestItems, bins)
	if err != nil {
		return nil, err
	}

	stats := result.Statistics
	stats.ExecutionTime = time.Since(start).String()
	stats.Algorithm = b.name

	return &PackingResult[B, I]{
		PackedBins:    result.PackedBins,
		UnpackedItems: result.UnpackedItems,
		Statistics:    stats,
	}, nil
}

func (b *BranchAndBound[B, I]) branchAndBound(
	ctx context.Context,
	current []int,
	used []bool,
	items []Item[I],
	bins []*Bin[B, I],
	startTime time.Time,
	bestSolution *[]int,
	bestFitness *float64,
	solutionsFound *int,
) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	if time.Since(startTime) > b.config.TimeLimit {
		return
	}

	if *solutionsFound >= b.config.MaxSolutions {
		return
	}

	if len(current) == len(items) {
		*solutionsFound++

		orderedItems := make([]Item[I], len(current))
		for i, idx := range current {
			orderedItems[i] = items[idx]
		}

		binsCopy := make([]*Bin[B, I], len(bins))
		for i, bin := range bins {
			binsCopy[i] = bin.Clone()
		}

		fitness := b.fitness(orderedItems, binsCopy)

		if fitness > *bestFitness {
			*bestFitness = fitness
			*bestSolution = make([]int, len(current))
			copy(*bestSolution, current)
		}

		return
	}

	for i := range items {
		if used[i] {
			continue
		}

		used[i] = true
		current = append(current, i)

		b.branchAndBound(
			ctx,
			current,
			used,
			items,
			bins,
			startTime,
			bestSolution,
			bestFitness,
			solutionsFound,
		)

		current = current[:len(current)-1]
		used[i] = false
	}
}
