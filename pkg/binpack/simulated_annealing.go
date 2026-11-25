package binpack

import (
	"context"
	"math"
	"math/rand"
	"time"
)

type SimulatedAnnealingConfig struct {
	InitialTemp float64
	CoolingRate float64
	Iterations  int
	MinTemp     float64
}

type SimulatedAnnealing[B any, I any] struct {
	name    string
	config  SimulatedAnnealingConfig
	base    Algorithm[B, I]
	fitness FitnessFunc[B, I]
}

func NewSimulatedAnnealing[B any, I any](config SimulatedAnnealingConfig, base Algorithm[B, I], fitness FitnessFunc[B, I]) *SimulatedAnnealing[B, I] {
	if config.InitialTemp == 0 {
		config.InitialTemp = 1000
	}
	if config.CoolingRate == 0 {
		config.CoolingRate = 0.995
	}
	if config.Iterations == 0 {
		config.Iterations = 1000
	}
	if config.MinTemp == 0 {
		config.MinTemp = 1
	}
	return &SimulatedAnnealing[B, I]{
		name:    "simulated-annealing",
		config:  config,
		base:    base,
		fitness: fitness,
	}
}

func (s *SimulatedAnnealing[B, I]) Name() string {
	return s.name
}

func (s *SimulatedAnnealing[B, I]) Pack(ctx context.Context, items []Item[I], bins []*Bin[B, I]) (*PackingResult[B, I], error) {
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

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	currentSolution := make([]int, len(items))
	for i := range currentSolution {
		currentSolution[i] = i
	}
	for i := len(currentSolution) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		currentSolution[i], currentSolution[j] = currentSolution[j], currentSolution[i]
	}

	bestSolution := make([]int, len(currentSolution))
	copy(bestSolution, currentSolution)

	currentFitness := s.evaluateSolution(currentSolution, items, bins)
	bestFitness := currentFitness

	temp := s.config.InitialTemp

	for iter := 0; iter < s.config.Iterations && temp > s.config.MinTemp; iter++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		newSolution := s.getNeighbor(currentSolution, rng)
		newFitness := s.evaluateSolution(newSolution, items, bins)

		delta := newFitness - currentFitness

		if delta > 0 || rng.Float64() < math.Exp(delta/temp) {
			currentSolution = newSolution
			currentFitness = newFitness

			if currentFitness > bestFitness {
				bestFitness = currentFitness
				copy(bestSolution, currentSolution)
			}
		}

		temp *= s.config.CoolingRate
	}

	bestItems := make([]Item[I], len(bestSolution))
	for i, idx := range bestSolution {
		bestItems[i] = items[idx]
	}

	result, err := s.base.Pack(ctx, bestItems, bins)
	if err != nil {
		return nil, err
	}

	stats := result.Statistics
	stats.ExecutionTime = time.Since(start).String()
	stats.Algorithm = s.name

	return &PackingResult[B, I]{
		PackedBins:    result.PackedBins,
		UnpackedItems: result.UnpackedItems,
		Statistics:    stats,
	}, nil
}

func (s *SimulatedAnnealing[B, I]) evaluateSolution(solution []int, items []Item[I], bins []*Bin[B, I]) float64 {
	orderedItems := make([]Item[I], len(solution))
	for i, idx := range solution {
		orderedItems[i] = items[idx]
	}

	binsCopy := make([]*Bin[B, I], len(bins))
	for i, bin := range bins {
		binsCopy[i] = bin.Clone()
	}

	return s.fitness(orderedItems, binsCopy)
}

func (s *SimulatedAnnealing[B, I]) getNeighbor(solution []int, rng *rand.Rand) []int {
	n := len(solution)
	newSolution := make([]int, n)
	copy(newSolution, solution)

	i := rng.Intn(n)
	j := rng.Intn(n)
	for i == j {
		j = rng.Intn(n)
	}

	newSolution[i], newSolution[j] = newSolution[j], newSolution[i]
	return newSolution
}
