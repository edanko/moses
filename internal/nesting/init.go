package nesting

// init registers the default algorithms
func init() {
	// Register greedy algorithms
	knapsackAlg := NewKnapsackAlgorithm()
	if err := RegisterAlgorithm(knapsackAlg); err != nil {
		panic("Failed to register knapsack algorithm: " + err.Error())
	}

	firstFitAlg := NewFirstFitAlgorithm()
	if err := RegisterAlgorithm(firstFitAlg); err != nil {
		panic("Failed to register first-fit algorithm: " + err.Error())
	}

	bestFitDecreasingAlg := NewBestFitDecreasingAlgorithm()
	if err := RegisterAlgorithm(bestFitDecreasingAlg); err != nil {
		panic("Failed to register best-fit decreasing algorithm: " + err.Error())
	}

	firstFitDecreasingAlg := NewFirstFitDecreasingAlgorithm()
	if err := RegisterAlgorithm(firstFitDecreasingAlg); err != nil {
		panic("Failed to register first-fit decreasing algorithm: " + err.Error())
	}

	// Register metaheuristic algorithms
	geneticAlg := NewGeneticAlgorithm()
	if err := RegisterAlgorithm(geneticAlg); err != nil {
		panic("Failed to register genetic algorithm: " + err.Error())
	}

	simulatedAnnealingAlg := NewSimulatedAnnealingAlgorithm()
	if err := RegisterAlgorithm(simulatedAnnealingAlg); err != nil {
		panic("Failed to register simulated annealing algorithm: " + err.Error())
	}

	antColonyAlg := NewAntColonyOptimizationAlgorithm()
	if err := RegisterAlgorithm(antColonyAlg); err != nil {
		panic("Failed to register ant colony optimization algorithm: " + err.Error())
	}

	// Register exact algorithms
	branchAndBoundAlg := NewBranchAndBoundAlgorithm()
	if err := RegisterAlgorithm(branchAndBoundAlg); err != nil {
		panic("Failed to register branch and bound algorithm: " + err.Error())
	}

	// Register multi-objective wrapper algorithms
	multiObjectiveFirstFit := NewMultiObjectiveAlgorithm(firstFitAlg)
	if err := RegisterAlgorithm(multiObjectiveFirstFit); err != nil {
		panic("Failed to register multi-objective first-fit algorithm: " + err.Error())
	}

	multiObjectiveKnapsack := NewMultiObjectiveAlgorithm(knapsackAlg)
	if err := RegisterAlgorithm(multiObjectiveKnapsack); err != nil {
		panic("Failed to register multi-objective knapsack algorithm: " + err.Error())
	}
}
