package nesting

import (
	"context"
	"math"
	"time"

	"github.com/edanko/moses/internal/domain"
	"github.com/edanko/moses/pkg/binpack"
)

// BranchAndBoundAlgorithm implements branch and bound for optimal nesting
type BranchAndBoundAlgorithm struct {
	timeLimit    time.Duration
	maxSolutions int
}

// NewBranchAndBoundAlgorithm creates a new branch and bound algorithm
func NewBranchAndBoundAlgorithm() *BranchAndBoundAlgorithm {
	return &BranchAndBoundAlgorithm{
		timeLimit:    30 * time.Second,
		maxSolutions: 10000,
	}
}

// Name returns the algorithm name
func (b *BranchAndBoundAlgorithm) Name() string {
	return "branch-and-bound"
}

// SetParameters sets the algorithm parameters
func (b *BranchAndBoundAlgorithm) SetParameters(timeLimit time.Duration, maxSolutions int) {
	b.timeLimit = timeLimit
	b.maxSolutions = maxSolutions
}

// Nest performs the nesting operation using branch and bound
func (b *BranchAndBoundAlgorithm) Nest(ctx context.Context, req *Request) (*Result, error) {
	items := toBinpackItems(req.Profiles)
	bins := toBinpackBins(req.StockBar)

	if len(req.Remnants) > 0 {
		remnantBins := toBinpackBins(req.Remnants)
		bins = append(remnantBins, bins...)
	}

	if len(items) == 0 || len(bins) == 0 {
		return &Result{
			Nests:      []domain.Nest{},
			Unnested:   req.Profiles,
			Statistics: Statistics{TotalProfiles: len(req.Profiles), Algorithm: b.Name()},
		}, nil
	}

	startTime := time.Now()

	bestSolution := make([]int, len(items))
	for i := range bestSolution {
		bestSolution[i] = i
	}

	bestFitness := math.Inf(-1)
	solutionsFound := 0

	// Simple branch and bound with pruning
	b.branchAndBound(
		ctx,
		make([]int, 0, len(items)),
		make([]bool, len(items)),
		items,
		bins,
		startTime,
		&bestSolution,
		&bestFitness,
		&solutionsFound,
	)

	// Apply best solution
	bestItems := make([]binpack.Item[domain.Part], len(bestSolution))
	for i, idx := range bestSolution {
		bestItems[i] = items[idx]
	}

	firstFit := binpack.NewFirstFitAlgorithm[domain.StockBar, domain.Part]()
	packingResult, err := firstFit.Pack(ctx, bestItems, bins)
	if err != nil {
		return nil, err
	}

	nests, unnested := fromBinpackResult(packingResult)

	return &Result{
		Nests:      nests,
		Unnested:   unnested,
		Statistics: convertStatistics(packingResult.Statistics, len(req.Profiles)),
	}, nil
}

func (b *BranchAndBoundAlgorithm) branchAndBound(
	ctx context.Context,
	current []int,
	used []bool,
	items []binpack.Item[domain.Part],
	bins []*binpack.Bin[domain.StockBar, domain.Part],
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

	// Check time limit
	if time.Since(startTime) > b.timeLimit {
		return
	}

	// Check solution limit
	if *solutionsFound >= b.maxSolutions {
		return
	}

	// Complete solution
	if len(current) == len(items) {
		*solutionsFound++

		// Evaluate solution
		orderedItems := make([]binpack.Item[domain.Part], len(current))
		for i, idx := range current {
			orderedItems[i] = items[idx]
		}

		binsCopy := make([]*binpack.Bin[domain.StockBar, domain.Part], len(bins))
		for i, bin := range bins {
			binsCopy[i] = bin.Clone()
		}

		fitness := EvaluateProfileOrder(orderedItems, binsCopy)

		if fitness > *bestFitness {
			*bestFitness = fitness
			*bestSolution = make([]int, len(current))
			copy(*bestSolution, current)
		}

		return
	}

	// Branch: try each unused item
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
