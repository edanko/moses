package nesting

import (
	"context"
	"math"
	"math/rand"

	"github.com/edanko/moses/internal/domain"
	"github.com/edanko/moses/pkg/binpack"
)

// AntColonyOptimizationAlgorithm implements ant colony optimization for nesting
type AntColonyOptimizationAlgorithm struct {
	numAnts       int
	maxIterations int
	alpha         float64 // pheromone importance
	beta          float64 // heuristic importance
	evaporation   float64 // pheromone evaporation rate
	q0            float64 // exploitation vs exploration parameter
}

// NewAntColonyOptimizationAlgorithm creates a new ant colony optimization algorithm
func NewAntColonyOptimizationAlgorithm() *AntColonyOptimizationAlgorithm {
	return &AntColonyOptimizationAlgorithm{
		numAnts:       10,
		maxIterations: 50,
		alpha:         1.0,
		beta:          2.0,
		evaporation:   0.5,
		q0:            0.9,
	}
}

// Name returns the algorithm name
func (a *AntColonyOptimizationAlgorithm) Name() string {
	return "ant-colony-optimization"
}

// SetParameters sets the algorithm parameters
func (a *AntColonyOptimizationAlgorithm) SetParameters(numAnts, maxIterations int, alpha, beta, evaporation, q0 float64) {
	a.numAnts = numAnts
	a.maxIterations = maxIterations
	a.alpha = alpha
	a.beta = beta
	a.evaporation = evaporation
	a.q0 = q0
}

// Nest performs the nesting operation using ant colony optimization
func (a *AntColonyOptimizationAlgorithm) Nest(ctx context.Context, req *Request) (*Result, error) {
	items := toBinpackItems(req.Profiles)
	rawBins := toBinpackBins(req.StockBar)

	if len(req.Remnants) > 0 {
		remnantBins := toBinpackBins(req.Remnants)
		rawBins = append(remnantBins, rawBins...)
	}

	if len(items) == 0 || len(rawBins) == 0 {
		return &Result{
			Nests:      []domain.Nest{},
			Unnested:   req.Profiles,
			Statistics: Statistics{TotalProfiles: len(req.Profiles), Algorithm: a.Name()},
		}, nil
	}

	numItems := len(items)

	// Initialize pheromone matrix
	pheromones := make([][]float64, numItems)
	for i := range pheromones {
		pheromones[i] = make([]float64, numItems)
		for j := range pheromones[i] {
			pheromones[i][j] = 1.0
		}
	}

	rand := rand.New(rand.NewSource(1))

	bestSolution := make([]int, numItems)
	bestFitness := math.Inf(-1)

	// Main ACO loop
	for iteration := 0; iteration < a.maxIterations; iteration++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// Each ant constructs a solution
		for ant := 0; ant < a.numAnts; ant++ {
			solution := a.constructSolution(pheromones, items, rand)

			// Evaluate solution
			orderedItems := make([]binpack.Item[domain.Part], numItems)
			for i, idx := range solution {
				orderedItems[i] = items[idx]
			}

			binsCopy := make([]*binpack.Bin[domain.StockBar, domain.Part], len(rawBins))
			for i, bin := range rawBins {
				binsCopy[i] = bin.Clone()
			}

			fitness := EvaluateProfileOrder(orderedItems, binsCopy)

			if fitness > bestFitness {
				bestFitness = fitness
				copy(bestSolution, solution)
			}
		}

		a.updatePheromones(pheromones, bestSolution, bestFitness)
	}

	bestItems := make([]binpack.Item[domain.Part], len(bestSolution))
	for i, idx := range bestSolution {
		bestItems[i] = items[idx]
	}

	firstFit := binpack.NewFirstFitAlgorithm[domain.StockBar, domain.Part]()
	packingResult, err := firstFit.Pack(ctx, bestItems, rawBins)
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

func (a *AntColonyOptimizationAlgorithm) constructSolution(pheromones [][]float64, items []binpack.Item[domain.Part], rand *rand.Rand) []int {
	numItems := len(items)
	solution := make([]int, 0, numItems)
	available := make([]bool, numItems)
	for i := range available {
		available[i] = true
	}

	// Start with a random item
	current := rand.Intn(numItems)
	solution = append(solution, current)
	available[current] = false

	for len(solution) < numItems {
		// Calculate probabilities
		probabilities := make([]float64, 0)
		indices := make([]int, 0)
		total := 0.0

		for i, avail := range available {
			if avail {
				tau := math.Pow(pheromones[current][i], a.alpha)
				eta := math.Pow(1.0/items[i].Length, a.beta)
				prob := tau * eta
				probabilities = append(probabilities, prob)
				indices = append(indices, i)
				total += prob
			}
		}

		if len(indices) == 0 {
			break
		}

		// Select next item
		var next int
		if rand.Float64() < a.q0 {
			// Exploitation: select best
			bestIdx := 0
			bestProb := probabilities[0]
			for i, prob := range probabilities {
				if prob > bestProb {
					bestProb = prob
					bestIdx = i
				}
			}
			next = indices[bestIdx]
		} else {
			// Exploration: select based on probability
			r := rand.Float64() * total
			cumsum := 0.0
			for i, prob := range probabilities {
				cumsum += prob
				if cumsum >= r {
					next = indices[i]
					break
				}
			}
		}

		solution = append(solution, next)
		available[next] = false
		current = next
	}

	return solution
}

func (a *AntColonyOptimizationAlgorithm) updatePheromones(
	pheromones [][]float64,
	bestSolution []int,
	bestFitness float64,
) {
	// Evaporate pheromones
	for i := range pheromones {
		for j := range pheromones[i] {
			pheromones[i][j] *= 1 - a.evaporation
		}
	}

	// Add pheromone based on best solution
	delta := 1.0 / (1.0 + bestFitness)
	for i := 0; i < len(bestSolution)-1; i++ {
		from := bestSolution[i]
		to := bestSolution[i+1]
		pheromones[from][to] += delta
		pheromones[to][from] += delta
	}
}
