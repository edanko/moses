package nesting

import (
	"context"
	"math"
	"sort"
	"time"
)

// MultiObjectiveAlgorithm implements multi-objective optimization for nesting
type MultiObjectiveAlgorithm struct {
	name              string
	algorithm         Algorithm // Base algorithm to use
	objectives        []Objective
	weightUtilization float64
	weightCuts        float64
	weightStockUsage  float64
	weightCompletion  float64
	weightWaste       float64
}

// Objective represents an optimization objective
type Objective struct {
	Name     string
	Weight   float64
	Minimize bool
	Evaluate func(*Result) float64
}

// ParetoSolution represents a solution in Pareto front
type ParetoSolution struct {
	Result       *Result
	Objectives   []float64
	DominatedBy  int
	Dominates    []int
	Rank         int
	CrowdingDist float64
}

// NewMultiObjectiveAlgorithm creates a new multi-objective nesting algorithm
func NewMultiObjectiveAlgorithm(baseAlgorithm Algorithm) *MultiObjectiveAlgorithm {
	return &MultiObjectiveAlgorithm{
		name:              "multi-objective",
		algorithm:         baseAlgorithm,
		objectives:        make([]Objective, 0),
		weightUtilization: 0.3,
		weightCuts:        0.2,
		weightStockUsage:  0.2,
		weightCompletion:  0.2,
		weightWaste:       0.1,
	}
}

// Name returns the algorithm name
func (m *MultiObjectiveAlgorithm) Name() string {
	return m.name + "-" + m.algorithm.Name()
}

// SetObjectives sets the optimization objectives and their weights
func (m *MultiObjectiveAlgorithm) SetObjectives(objectives []Objective) {
	m.objectives = objectives
}

// SetWeights sets weights for default objectives
func (m *MultiObjectiveAlgorithm) SetWeights(utilization, cuts, stock, completion, waste float64) {
	m.weightUtilization = utilization
	m.weightCuts = cuts
	m.weightStockUsage = stock
	m.weightCompletion = completion
	m.weightWaste = waste
}

// AddObjective adds a custom objective
func (m *MultiObjectiveAlgorithm) AddObjective(name string, weight float64, minimize bool, evaluate func(*Result) float64) {
	m.objectives = append(m.objectives, Objective{
		Name:     name,
		Weight:   weight,
		Minimize: minimize,
		Evaluate: evaluate,
	})
}

// Nest performs multi-objective nesting optimization
func (m *MultiObjectiveAlgorithm) Nest(ctx context.Context, req *Request) (*Result, error) {
	// Check context cancellation at start
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	start := time.Now()

	// Generate multiple solutions using different strategies
	solutions := make([]*ParetoSolution, 0)

	// Strategy 1: Use base algorithm normally
	baseResult, err := m.algorithm.Nest(ctx, req)
	if err != nil {
		return nil, err
	}
	solution := &ParetoSolution{
		Result:      baseResult,
		Objectives:  m.evaluateObjectives(baseResult),
		DominatedBy: 0,
		Dominates:   make([]int, 0),
	}
	solutions = append(solutions, solution)

	// Strategy 2: Optimize for material utilization
	utilReq := m.optimizeForUtilization(req)
	utilResult, err := m.algorithm.Nest(ctx, utilReq)
	if err != nil {
		return nil, err
	}
	utilSolution := &ParetoSolution{
		Result:      utilResult,
		Objectives:  m.evaluateObjectives(utilResult),
		DominatedBy: 0,
		Dominates:   make([]int, 0),
	}
	solutions = append(solutions, utilSolution)

	// Strategy 3: Optimize for minimum cuts
	cutsReq := m.optimizeForCuts(req)
	cutsResult, err := m.algorithm.Nest(ctx, cutsReq)
	if err != nil {
		return nil, err
	}
	cutsSolution := &ParetoSolution{
		Result:      cutsResult,
		Objectives:  m.evaluateObjectives(cutsResult),
		DominatedBy: 0,
		Dominates:   make([]int, 0),
	}
	solutions = append(solutions, cutsSolution)

	// Strategy 4: Optimize for minimum stock usage
	stockReq := m.optimizeForStockUsage(req)
	stockResult, err := m.algorithm.Nest(ctx, stockReq)
	if err != nil {
		return nil, err
	}
	stockSolution := &ParetoSolution{
		Result:      stockResult,
		Objectives:  m.evaluateObjectives(stockResult),
		DominatedBy: 0,
		Dominates:   make([]int, 0),
	}
	solutions = append(solutions, stockSolution)

	// Find Pareto front
	paretoFront := m.findParetoFront(solutions)

	if len(paretoFront) == 0 {
		// Fallback to base algorithm result
		return baseResult, nil
	}

	// Select best solution using weighted sum or crowding distance
	bestSolution := m.selectBestSolution(paretoFront)

	// Update statistics to include multi-objective information
	bestSolution.Result.Statistics.ExecutionTime = time.Since(start).String()
	bestSolution.Result.Statistics.Algorithm = m.Name()

	return bestSolution.Result, nil
}

// evaluateObjectives evaluates all objectives for a solution
func (m *MultiObjectiveAlgorithm) evaluateObjectives(result *Result) []float64 {
	objectives := make([]float64, len(m.objectives))

	for i, obj := range m.objectives {
		objectives[i] = obj.Evaluate(result)
	}

	// If no custom objectives, use default ones
	if len(m.objectives) == 0 {
		objectives = []float64{
			result.Statistics.MaterialUtilization, // Maximize
			float64(len(result.Nests)),            // Minimize (number of nests = cuts)
			float64(result.Statistics.TotalProfiles - result.Statistics.UnnestedProfiles), // Maximize
			result.Statistics.TotalWaste, // Minimize
		}
	}

	return objectives
}

// optimizeForUtilization modifies request to optimize for material utilization
func (m *MultiObjectiveAlgorithm) optimizeForUtilization(req *Request) *Request {
	optimizedReq := *req
	// Reduce spacing to fit more profiles
	optimizedReq.Spacing = req.Spacing * 0.5
	return &optimizedReq
}

// optimizeForCuts modifies request to minimize number of cuts
func (m *MultiObjectiveAlgorithm) optimizeForCuts(req *Request) *Request {
	optimizedReq := *req
	// Increase spacing to encourage larger pieces per nest
	optimizedReq.Spacing = req.Spacing * 1.5
	return &optimizedReq
}

// optimizeForStockUsage modifies request to minimize stock usage
func (m *MultiObjectiveAlgorithm) optimizeForStockUsage(req *Request) *Request {
	optimizedReq := *req
	// Prioritize longer stock pieces to reduce count
	if len(req.StockBar) > 0 {
		sort.Slice(optimizedReq.StockBar, func(i, j int) bool {
			return optimizedReq.StockBar[i].Length() > optimizedReq.StockBar[j].Length()
		})
	}
	return &optimizedReq
}

// findParetoFront finds non-dominated solutions
func (m *MultiObjectiveAlgorithm) findParetoFront(solutions []*ParetoSolution) []*ParetoSolution {
	if len(solutions) == 0 {
		return solutions
	}

	// Initialize domination counts
	for i := range solutions {
		solutions[i].DominatedBy = 0
		solutions[i].Dominates = make([]int, 0)
	}

	// Check domination between all pairs
	for i := range solutions {
		for j := range solutions {
			if i != j {
				if m.dominates(solutions[i], solutions[j]) {
					solutions[i].Dominates = append(solutions[i].Dominates, j)
					solutions[j].DominatedBy++
				}
			}
		}
	}

	// Find non-dominated solutions (Pareto front)
	paretoFront := make([]*ParetoSolution, 0)
	for _, solution := range solutions {
		if solution.DominatedBy == 0 {
			paretoFront = append(paretoFront, solution)
		}
	}

	// Calculate crowding distance for diversity
	m.calculateCrowdingDistance(paretoFront)

	return paretoFront
}

// dominates checks if solution1 dominates solution2
func (m *MultiObjectiveAlgorithm) dominates(sol1, sol2 *ParetoSolution) bool {
	betterInAny := false

	objectives := m.objectives
	if len(objectives) == 0 {
		objectives = []Objective{
			{Name: "utilization", Minimize: false},
			{Name: "cuts", Minimize: true},
			{Name: "completion", Minimize: false},
			{Name: "waste", Minimize: true},
		}
	}

	for i := 0; i < len(sol1.Objectives) && i < len(sol2.Objectives); i++ {
		obj := objectives[i]

		if obj.Minimize {
			if sol1.Objectives[i] > sol2.Objectives[i] {
				return false // sol1 is worse in this objective
			}
			if sol1.Objectives[i] < sol2.Objectives[i] {
				betterInAny = true
			}
		} else {
			if sol1.Objectives[i] < sol2.Objectives[i] {
				return false // sol1 is worse in this objective
			}
			if sol1.Objectives[i] > sol2.Objectives[i] {
				betterInAny = true
			}
		}
	}

	return betterInAny
}

// calculateCrowdingDistance calculates crowding distance for diversity
func (m *MultiObjectiveAlgorithm) calculateCrowdingDistance(solutions []*ParetoSolution) {
	if len(solutions) == 0 {
		return
	}

	// Initialize distances
	for _, sol := range solutions {
		sol.CrowdingDist = 0
	}

	// Calculate for each objective
	for objIdx := 0; objIdx < len(m.objectives); objIdx++ {
		// Sort by this objective
		sort.Slice(solutions, func(i, j int) bool {
			return solutions[i].Objectives[objIdx] < solutions[j].Objectives[objIdx]
		})

		// Set infinite distance for boundary solutions
		if len(solutions) > 0 {
			solutions[0].CrowdingDist = math.Inf(1)
			if len(solutions) > 1 {
				solutions[len(solutions)-1].CrowdingDist = math.Inf(1)
			}
		}

		// Calculate crowding distance for intermediate solutions
		if len(solutions) > 2 {
			objRange := solutions[len(solutions)-1].Objectives[objIdx] - solutions[0].Objectives[objIdx]
			if objRange > 0 {
				for i := 1; i < len(solutions)-1; i++ {
					if solutions[i].CrowdingDist != math.Inf(1) {
						distance := solutions[i+1].Objectives[objIdx] - solutions[i-1].Objectives[objIdx]
						solutions[i].CrowdingDist += distance / objRange
					}
				}
			}
		}
	}
}

// selectBestSolution selects the best solution from Pareto front
func (m *MultiObjectiveAlgorithm) selectBestSolution(paretoFront []*ParetoSolution) *ParetoSolution {
	if len(paretoFront) == 0 {
		return nil
	}
	if len(paretoFront) == 1 {
		return paretoFront[0]
	}

	// Use weighted sum approach
	bestSolution := paretoFront[0]
	bestScore := 0.0

	for _, solution := range paretoFront {
		score := m.calculateWeightedScore(solution)
		if score > bestScore {
			bestScore = score
			bestSolution = solution
		}
	}

	return bestSolution
}

// calculateWeightedScore calculates weighted score for a solution
func (m *MultiObjectiveAlgorithm) calculateWeightedScore(solution *ParetoSolution) float64 {
	score := 0.0

	// Use default objectives if none specified
	if len(m.objectives) == 0 {
		utilization := solution.Result.Statistics.MaterialUtilization / 100.0
		cuts := 1.0 / (1.0 + float64(len(solution.Result.Nests)))
		completion := float64(solution.Result.Statistics.NestedProfiles) / float64(solution.Result.Statistics.TotalProfiles)
		waste := 1.0 / (1.0 + solution.Result.Statistics.TotalWaste)

		score = utilization*m.weightUtilization + cuts*m.weightCuts + completion*m.weightCompletion + waste*m.weightWaste
	} else {
		// Use custom objectives
		for i := range m.objectives {
			if i < len(solution.Objectives) {
				obj := m.objectives[i]
				if obj.Minimize {
					// Normalize minimization (lower is better)
					score += (1.0 / (1.0 + solution.Objectives[i])) * obj.Weight
				} else {
					// Normalize maximization (higher is better)
					score += solution.Objectives[i] * obj.Weight
				}
			}
		}
	}

	return score
}
