package nesting

import (
	"context"
	"math"
	"math/rand"

	"github.com/edanko/moses/internal/domain"
	"github.com/edanko/moses/pkg/binpack"
)

// GeneticAlgorithm implements a genetic algorithm for nesting
type GeneticAlgorithm struct {
	populationSize int
	generations    int
	mutationRate   float64
	crossoverRate  float64
}

// NewGeneticAlgorithm creates a new genetic algorithm
func NewGeneticAlgorithm() *GeneticAlgorithm {
	return &GeneticAlgorithm{
		populationSize: 50,
		generations:    100,
		mutationRate:   0.1,
		crossoverRate:  0.8,
	}
}

// Name returns the algorithm name
func (g *GeneticAlgorithm) Name() string {
	return "genetic-algorithm"
}

// SetParameters sets the algorithm parameters
func (g *GeneticAlgorithm) SetParameters(populationSize, generations int, mutationRate, crossoverRate float64) {
	g.populationSize = populationSize
	g.generations = generations
	g.mutationRate = mutationRate
	g.crossoverRate = crossoverRate
}

// Nest performs the nesting operation using genetic algorithm
func (g *GeneticAlgorithm) Nest(ctx context.Context, req *Request) (*Result, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

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
			Statistics: Statistics{TotalProfiles: len(req.Profiles), Algorithm: g.Name()},
		}, nil
	}

	rand := rand.New(rand.NewSource(1))

	// Initialize population with random permutations
	population := make([][]int, g.populationSize)
	for i := range population {
		population[i] = make([]int, len(items))
		for j := range items {
			population[i][j] = j
		}
		// Shuffle
		for j := len(items) - 1; j > 0; j-- {
			k := rand.Intn(j + 1)
			population[i][j], population[i][k] = population[i][k], population[i][j]
		}
	}

	bestFitness := math.Inf(-1)
	var bestSolution []int

	// Evolution loop
	for gen := 0; gen < g.generations; gen++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// Evaluate fitness for each individual
		fitnesses := make([]float64, g.populationSize)
		for i, individual := range population {
			fitness := g.evaluateIndividual(individual, items, bins)
			fitnesses[i] = fitness
			if fitness > bestFitness {
				bestFitness = fitness
				bestSolution = make([]int, len(individual))
				copy(bestSolution, individual)
			}
		}

		// Create new population
		newPopulation := make([][]int, 0, g.populationSize)
		newPopulation = append(newPopulation, bestSolution) // Elitism

		for len(newPopulation) < g.populationSize {
			parent1 := g.tournamentSelect(population, fitnesses, rand)
			parent2 := g.tournamentSelect(population, fitnesses, rand)

			if rand.Float64() < g.crossoverRate {
				child1, child2 := g.crossover(parent1, parent2, rand)
				newPopulation = append(newPopulation, child1, child2)
			} else {
				newPopulation = append(newPopulation, parent1, parent2)
			}
		}

		population = newPopulation[:g.populationSize]

		// Mutation
		for i := 1; i < g.populationSize; i++ {
			if rand.Float64() < g.mutationRate {
				g.mutate(population[i], rand)
			}
		}
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

func (g *GeneticAlgorithm) evaluateIndividual(order []int, items []binpack.Item[domain.Part], bins []*binpack.Bin[domain.StockBar, domain.Part]) float64 {
	orderedItems := make([]binpack.Item[domain.Part], len(order))
	for i, idx := range order {
		orderedItems[i] = items[idx]
	}

	binsCopy := make([]*binpack.Bin[domain.StockBar, domain.Part], len(bins))
	for i, bin := range bins {
		binsCopy[i] = bin.Clone()
	}

	return EvaluateProfileOrder(orderedItems, binsCopy)
}

func (g *GeneticAlgorithm) tournamentSelect(population [][]int, fitnesses []float64, rand *rand.Rand) []int {
	tournamentSize := 3
	bestIdx := rand.Intn(len(population))
	bestFitness := fitnesses[bestIdx]

	for i := 1; i < tournamentSize; i++ {
		idx := rand.Intn(len(population))
		if fitnesses[idx] > bestFitness {
			bestIdx = idx
			bestFitness = fitnesses[idx]
		}
	}

	solution := make([]int, len(population[bestIdx]))
	copy(solution, population[bestIdx])
	return solution
}

func (g *GeneticAlgorithm) crossover(parent1, parent2 []int, rand *rand.Rand) ([]int, []int) {
	n := len(parent1)
	point := rand.Intn(n-1) + 1

	child1 := make([]int, n)
	child2 := make([]int, n)

	used1 := make(map[int]bool)
	used2 := make(map[int]bool)

	// Copy first part
	for i := range point {
		child1[i] = parent1[i]
		child2[i] = parent2[i]
		used1[parent1[i]] = true
		used2[parent2[i]] = true
	}

	// Fill remaining from other parent
	idx1, idx2 := point, point
	for i := range n {
		if !used1[parent2[i]] {
			child1[idx1] = parent2[i]
			idx1++
		}
		if !used2[parent1[i]] {
			child2[idx2] = parent1[i]
			idx2++
		}
	}

	return child1, child2
}

func (g *GeneticAlgorithm) mutate(individual []int, rand *rand.Rand) {
	n := len(individual)
	i := rand.Intn(n)
	j := rand.Intn(n)
	individual[i], individual[j] = individual[j], individual[i]
}
