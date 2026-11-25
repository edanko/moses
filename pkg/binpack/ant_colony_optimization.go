package binpack

import (
	"context"
	"math"
	"math/rand"
	"time"
)

type AntColonyConfig struct {
	Ants        int
	Iterations  int
	Alpha       float64
	Beta        float64
	Evaporation float64
	Q           float64
}

type AntColonyOptimization[B any, I any] struct {
	name    string
	config  AntColonyConfig
	base    Algorithm[B, I]
	fitness FitnessFunc[B, I]
}

func NewAntColonyOptimization[B any, I any](config AntColonyConfig, base Algorithm[B, I], fitness FitnessFunc[B, I]) *AntColonyOptimization[B, I] {
	if config.Ants == 0 {
		config.Ants = 10
	}
	if config.Iterations == 0 {
		config.Iterations = 50
	}
	if config.Alpha == 0 {
		config.Alpha = 1.0
	}
	if config.Beta == 0 {
		config.Beta = 2.0
	}
	if config.Evaporation == 0 {
		config.Evaporation = 0.1
	}
	if config.Q == 0 {
		config.Q = 1.0
	}
	return &AntColonyOptimization[B, I]{
		name:    "ant-colony-optimization",
		config:  config,
		base:    base,
		fitness: fitness,
	}
}

func (a *AntColonyOptimization[B, I]) Name() string {
	return a.name
}

func (a *AntColonyOptimization[B, I]) Pack(ctx context.Context, items []Item[I], bins []*Bin[B, I]) (*PackingResult[B, I], error) {
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

	numItems := len(items)
	pheromones := make([][]float64, numItems)
	for i := range pheromones {
		pheromones[i] = make([]float64, numItems)
		for j := range pheromones[i] {
			pheromones[i][j] = 1.0
		}
	}

	bestSolution := make([]int, numItems)
	for i := range bestSolution {
		bestSolution[i] = i
	}
	bestFitness := math.Inf(-1)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for iter := 0; iter < a.config.Iterations; iter++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		solutions := make([][]int, a.config.Ants)
		fitnesses := make([]float64, a.config.Ants)

		for ant := 0; ant < a.config.Ants; ant++ {
			solution := a.constructSolution(pheromones, items, rng)
			solutions[ant] = solution
			fitnesses[ant] = a.evaluateSolution(solution, items, bins)

			if fitnesses[ant] > bestFitness {
				bestFitness = fitnesses[ant]
				copy(bestSolution, solution)
			}
		}

		for i := range pheromones {
			for j := range pheromones[i] {
				pheromones[i][j] *= (1.0 - a.config.Evaporation)
			}
		}

		for ant := 0; ant < a.config.Ants; ant++ {
			delta := a.config.Q * fitnesses[ant]
			solution := solutions[ant]
			for k := 0; k < len(solution)-1; k++ {
				i, j := solution[k], solution[k+1]
				if i < len(pheromones) && j < len(pheromones[i]) {
					pheromones[i][j] += delta
				}
			}
		}
	}

	bestItems := make([]Item[I], len(bestSolution))
	for i, idx := range bestSolution {
		bestItems[i] = items[idx]
	}

	result, err := a.base.Pack(ctx, bestItems, bins)
	if err != nil {
		return nil, err
	}

	stats := result.Statistics
	stats.ExecutionTime = time.Since(start).String()
	stats.Algorithm = a.name

	return &PackingResult[B, I]{
		PackedBins:    result.PackedBins,
		UnpackedItems: result.UnpackedItems,
		Statistics:    stats,
	}, nil
}

func (a *AntColonyOptimization[B, I]) constructSolution(pheromones [][]float64, items []Item[I], rng *rand.Rand) []int {
	n := len(items)
	solution := make([]int, 0, n)
	used := make([]bool, n)

	current := rng.Intn(n)
	solution = append(solution, current)
	used[current] = true

	for len(solution) < n {
		probs := make([]float64, n)
		total := 0.0

		for i := range n {
			if used[i] {
				continue
			}
			pheromone := math.Pow(pheromones[current][i], a.config.Alpha)
			heuristic := math.Pow(1.0/float64(i+1), a.config.Beta)
			probs[i] = pheromone * heuristic
			total += probs[i]
		}

		if total == 0 {
			for i := range n {
				if !used[i] {
					solution = append(solution, i)
					used[i] = true
					current = i
					break
				}
			}
			continue
		}

		r := rng.Float64() * total
		sum := 0.0
		selected := -1

		for i := range n {
			if used[i] {
				continue
			}
			sum += probs[i]
			if sum >= r {
				selected = i
				break
			}
		}

		if selected == -1 {
			for i := range n {
				if !used[i] {
					selected = i
					break
				}
			}
		}

		if selected == -1 {
			break
		}

		solution = append(solution, selected)
		used[selected] = true
		current = selected
	}

	return solution
}

func (a *AntColonyOptimization[B, I]) evaluateSolution(solution []int, items []Item[I], bins []*Bin[B, I]) float64 {
	orderedItems := make([]Item[I], len(solution))
	for i, idx := range solution {
		orderedItems[i] = items[idx]
	}

	binsCopy := make([]*Bin[B, I], len(bins))
	for i, bin := range bins {
		binsCopy[i] = bin.Clone()
	}

	return a.fitness(orderedItems, binsCopy)
}
