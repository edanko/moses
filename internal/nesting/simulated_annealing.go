package nesting

import (
	"context"
	"math"
	"math/rand"

	"github.com/edanko/moses/internal/domain"
	"github.com/edanko/moses/pkg/binpack"
)

// SimulatedAnnealingAlgorithm implements simulated annealing for nesting
type SimulatedAnnealingAlgorithm struct {
	initialTemp float64
	coolingRate float64
	minTemp     float64
	iterations  int
}

// NewSimulatedAnnealingAlgorithm creates a new simulated annealing algorithm
func NewSimulatedAnnealingAlgorithm() *SimulatedAnnealingAlgorithm {
	return &SimulatedAnnealingAlgorithm{
		initialTemp: 100.0,
		coolingRate: 0.95,
		minTemp:     0.001,
		iterations:  1000,
	}
}

// Name returns the algorithm name
func (s *SimulatedAnnealingAlgorithm) Name() string {
	return "simulated-annealing"
}

// SetParameters sets the algorithm parameters
func (s *SimulatedAnnealingAlgorithm) SetParameters(initialTemp, coolingRate, minTemp float64, iterations int) {
	s.initialTemp = initialTemp
	s.coolingRate = coolingRate
	s.minTemp = minTemp
	s.iterations = iterations
}

// Nest performs the nesting operation using simulated annealing
func (s *SimulatedAnnealingAlgorithm) Nest(ctx context.Context, req *Request) (*Result, error) {
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
			Statistics: Statistics{TotalProfiles: len(req.Profiles), Algorithm: s.Name()},
		}, nil
	}

	rand := rand.New(rand.NewSource(1))

	// Initialize current solution with random permutation
	currentSolution := make([]int, len(items))
	for i := range currentSolution {
		currentSolution[i] = i
	}
	for i := len(currentSolution) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		currentSolution[i], currentSolution[j] = currentSolution[j], currentSolution[i]
	}

	bestSolution := make([]int, len(currentSolution))
	copy(bestSolution, currentSolution)

	currentItems := make([]binpack.Item[domain.Part], len(items))
	for i, idx := range currentSolution {
		currentItems[i] = items[idx]
	}

	binsCopy := make([]*binpack.Bin[domain.StockBar, domain.Part], len(bins))
	for i, bin := range bins {
		binsCopy[i] = bin.Clone()
	}

	currentFitness := EvaluateProfileOrder(currentItems, binsCopy)
	bestFitness := currentFitness

	temperature := s.initialTemp

	for iter := 0; iter < s.iterations && temperature > s.minTemp; iter++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// Create neighbor by swapping two random positions
		newSolution := make([]int, len(currentSolution))
		copy(newSolution, currentSolution)
		i := rand.Intn(len(newSolution))
		j := rand.Intn(len(newSolution))
		newSolution[i], newSolution[j] = newSolution[j], newSolution[i]

		// Evaluate neighbor
		newItems := make([]binpack.Item[domain.Part], len(items))
		for k, idx := range newSolution {
			newItems[k] = items[idx]
		}

		binsCopy = make([]*binpack.Bin[domain.StockBar, domain.Part], len(bins))
		for k, bin := range bins {
			binsCopy[k] = bin.Clone()
		}

		newFitness := EvaluateProfileOrder(newItems, binsCopy)
		delta := newFitness - currentFitness

		// Accept if better, or with probability based on temperature
		if delta > 0 || rand.Float64() < math.Exp(delta/temperature) {
			currentSolution = newSolution
			currentFitness = newFitness

			if currentFitness > bestFitness {
				bestFitness = currentFitness
				copy(bestSolution, currentSolution)
			}
		}

		temperature *= s.coolingRate
	}

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
