package binpack

import (
	"context"
	"math"
	"math/rand"
	"time"
)

type GeneticAlgorithmConfig struct {
	PopulationSize int
	Generations    int
	MutationRate   float64
	CrossoverRate  float64
}

type GeneticAlgorithm[B any, I any] struct {
	name    string
	config  GeneticAlgorithmConfig
	base    Algorithm[B, I]
	fitness FitnessFunc[B, I]
}

type FitnessFunc[B any, I any] func(items []Item[I], bins []*Bin[B, I]) float64

func NewGeneticAlgorithm[B any, I any](config GeneticAlgorithmConfig, base Algorithm[B, I], fitness FitnessFunc[B, I]) *GeneticAlgorithm[B, I] {
	if config.PopulationSize == 0 {
		config.PopulationSize = 50
	}
	if config.Generations == 0 {
		config.Generations = 100
	}
	if config.MutationRate == 0 {
		config.MutationRate = 0.1
	}
	if config.CrossoverRate == 0 {
		config.CrossoverRate = 0.8
	}
	return &GeneticAlgorithm[B, I]{
		name:    "genetic-algorithm",
		config:  config,
		base:    base,
		fitness: fitness,
	}
}

func (g *GeneticAlgorithm[B, I]) Name() string {
	return g.name
}

func (g *GeneticAlgorithm[B, I]) Pack(ctx context.Context, items []Item[I], bins []*Bin[B, I]) (*PackingResult[B, I], error) {
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

	population := g.initializePopulation(len(items), rng)
	bestSolution := make([]int, len(items))
	bestFitness := math.Inf(-1)

	for gen := 0; gen < g.config.Generations; gen++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		fitnesses := make([]float64, len(population))
		for i, individual := range population {
			fitnesses[i] = g.evaluateIndividual(individual, items, bins)
			if fitnesses[i] > bestFitness {
				bestFitness = fitnesses[i]
				copy(bestSolution, individual)
			}
		}

		population = g.evolvePopulation(population, fitnesses, rng)
	}

	bestItems := make([]Item[I], len(bestSolution))
	for i, idx := range bestSolution {
		bestItems[i] = items[idx]
	}

	result, err := g.base.Pack(ctx, bestItems, bins)
	if err != nil {
		return nil, err
	}

	stats := result.Statistics
	stats.ExecutionTime = time.Since(start).String()
	stats.Algorithm = g.name

	return &PackingResult[B, I]{
		PackedBins:    result.PackedBins,
		UnpackedItems: result.UnpackedItems,
		Statistics:    stats,
	}, nil
}

func (g *GeneticAlgorithm[B, I]) initializePopulation(n int, rng *rand.Rand) [][]int {
	population := make([][]int, g.config.PopulationSize)
	for i := range population {
		population[i] = make([]int, n)
		for j := range n {
			population[i][j] = j
		}
		for j := n - 1; j > 0; j-- {
			k := rng.Intn(j + 1)
			population[i][j], population[i][k] = population[i][k], population[i][j]
		}
	}
	return population
}

func (g *GeneticAlgorithm[B, I]) evaluateIndividual(individual []int, items []Item[I], bins []*Bin[B, I]) float64 {
	orderedItems := make([]Item[I], len(individual))
	for i, idx := range individual {
		orderedItems[i] = items[idx]
	}

	binsCopy := make([]*Bin[B, I], len(bins))
	for i, bin := range bins {
		binsCopy[i] = bin.Clone()
	}

	return g.fitness(orderedItems, binsCopy)
}

func (g *GeneticAlgorithm[B, I]) evolvePopulation(population [][]int, fitnesses []float64, rng *rand.Rand) [][]int {
	newPopulation := make([][]int, 0, len(population))

	for len(newPopulation) < g.config.PopulationSize {
		parent1 := g.selectParent(population, fitnesses, rng)
		parent2 := g.selectParent(population, fitnesses, rng)

		child1, child2 := g.crossover(parent1, parent2, rng)

		if rng.Float64() < g.config.MutationRate {
			child1 = g.mutate(child1, rng)
		}
		if rng.Float64() < g.config.MutationRate {
			child2 = g.mutate(child2, rng)
		}

		newPopulation = append(newPopulation, child1)
		if len(newPopulation) < g.config.PopulationSize {
			newPopulation = append(newPopulation, child2)
		}
	}

	return newPopulation
}

func (g *GeneticAlgorithm[B, I]) selectParent(population [][]int, fitnesses []float64, rng *rand.Rand) []int {
	tournamentSize := 3
	bestIdx := rng.Intn(len(population))
	bestFitness := fitnesses[bestIdx]

	for i := 1; i < tournamentSize; i++ {
		idx := rng.Intn(len(population))
		if fitnesses[idx] > bestFitness {
			bestFitness = fitnesses[idx]
			bestIdx = idx
		}
	}

	return population[bestIdx]
}

func (g *GeneticAlgorithm[B, I]) crossover(parent1, parent2 []int, rng *rand.Rand) (child1, child2 []int) {
	if rng.Float64() > g.config.CrossoverRate {
		return parent1, parent2
	}

	n := len(parent1)
	child1 = make([]int, n)
	child2 = make([]int, n)

	point := rng.Intn(n)
	for i := range point {
		child1[i] = parent1[i]
		child2[i] = parent2[i]
	}

	child1[point] = parent2[point]
	child2[point] = parent1[point]

	return child1, child2
}

func (g *GeneticAlgorithm[B, I]) mutate(individual []int, rng *rand.Rand) []int {
	n := len(individual)
	i := rng.Intn(n)
	j := rng.Intn(n)
	individual[i], individual[j] = individual[j], individual[i]
	return individual
}
