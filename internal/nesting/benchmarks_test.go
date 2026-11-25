package nesting

import (
	"context"
	"testing"
	"time"
)

// BenchmarkAllAlgorithms compares performance of all algorithms
func BenchmarkAllAlgorithms(b *testing.B) {
	testData := getSimpleTestData()
	ctx := context.Background()
	req := &Request{
		Profiles:   testData.Profiles,
		StockBar:   testData.StockBar,
		Remnants:   testData.Remnants,
		Spacing:    10.0,
		Parameters: make(map[string]any),
	}

	algorithms := []struct {
		name      string
		algorithm Algorithm
		prepare   func(Algorithm)
	}{
		{
			name:      "FirstFit",
			algorithm: NewFirstFitAlgorithm(),
			prepare:   func(algo Algorithm) {},
		},
		{
			name:      "Knapsack",
			algorithm: NewKnapsackAlgorithm(),
			prepare:   func(algo Algorithm) {},
		},
		{
			name:      "BestFitDecreasing",
			algorithm: NewBestFitDecreasingAlgorithm(),
			prepare:   func(algo Algorithm) {},
		},
		{
			name:      "FirstFitDecreasing",
			algorithm: NewFirstFitDecreasingAlgorithm(),
			prepare:   func(algo Algorithm) {},
		},
	}

	for _, algo := range algorithms {
		b.Run(algo.name, func(b *testing.B) {
			algo.prepare(algo.algorithm)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := algo.algorithm.Nest(ctx, req)
				if err != nil {
					b.Fatalf("%s algorithm error: %v", algo.name, err)
				}
			}
		})
	}
}

// BenchmarkMetaheuristicAlgorithms compares performance of metaheuristic algorithms
func BenchmarkMetaheuristicAlgorithms(b *testing.B) {
	testData := getSimpleTestData()
	ctx := context.Background()
	req := &Request{
		Profiles:   testData.Profiles,
		StockBar:   testData.StockBar,
		Remnants:   testData.Remnants,
		Spacing:    10.0,
		Parameters: make(map[string]any),
	}

	algorithms := []struct {
		name      string
		algorithm Algorithm
		prepare   func(Algorithm)
	}{
		{
			name:      "Genetic",
			algorithm: NewGeneticAlgorithm(),
			prepare: func(algo Algorithm) {
				if ga, ok := algo.(*GeneticAlgorithm); ok {
					ga.SetParameters(10, 20, 0.1, 0.2) // Small parameters for benchmarking
				}
			},
		},
		{
			name:      "SimulatedAnnealing",
			algorithm: NewSimulatedAnnealingAlgorithm(),
			prepare: func(algo Algorithm) {
				if sa, ok := algo.(*SimulatedAnnealingAlgorithm); ok {
					sa.SetParameters(50.0, 0.9, 1.0, 10) // Small parameters for benchmarking
				}
			},
		},
		{
			name:      "AntColonyOptimization",
			algorithm: NewAntColonyOptimizationAlgorithm(),
			prepare: func(algo Algorithm) {
				if aco, ok := algo.(*AntColonyOptimizationAlgorithm); ok {
					aco.SetParameters(5, 10, 1.0, 2.0, 0.8, 0.1) // numAnts, numIterations, alpha, beta, q0, tau0
				}
			},
		},
	}

	for _, algo := range algorithms {
		b.Run(algo.name, func(b *testing.B) {
			algo.prepare(algo.algorithm)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := algo.algorithm.Nest(ctx, req)
				if err != nil {
					b.Fatalf("%s algorithm error: %v", algo.name, err)
				}
			}
		})
	}
}

// BenchmarkExactOptimization compares performance of exact optimization
func BenchmarkExactOptimization(b *testing.B) {
	testData := getSimpleTestData()
	ctx := context.Background()
	req := &Request{
		Profiles:   testData.Profiles,
		StockBar:   testData.StockBar,
		Remnants:   testData.Remnants,
		Spacing:    10.0,
		Parameters: make(map[string]any),
	}

	algorithms := []struct {
		name      string
		algorithm Algorithm
		prepare   func(Algorithm)
	}{
		{
			name:      "BranchAndBound",
			algorithm: NewBranchAndBoundAlgorithm(),
			prepare: func(algo Algorithm) {
				if bb, ok := algo.(*BranchAndBoundAlgorithm); ok {
					bb.SetParameters(100*time.Millisecond, 50) // Very limited for benchmarking
				}
			},
		},
	}

	for _, algo := range algorithms {
		b.Run(algo.name, func(b *testing.B) {
			algo.prepare(algo.algorithm)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := algo.algorithm.Nest(ctx, req)
				if err != nil {
					b.Fatalf("%s algorithm error: %v", algo.name, err)
				}
			}
		})
	}
}

// BenchmarkMultiObjective compares performance of multi-objective algorithms
func BenchmarkMultiObjective(b *testing.B) {
	testData := getSimpleTestData()
	ctx := context.Background()
	req := &Request{
		Profiles:   testData.Profiles,
		StockBar:   testData.StockBar,
		Remnants:   testData.Remnants,
		Spacing:    10.0,
		Parameters: make(map[string]any),
	}

	algorithms := []struct {
		name      string
		algorithm Algorithm
		prepare   func(Algorithm)
	}{
		{
			name:      "MultiObjective-FirstFit",
			algorithm: NewMultiObjectiveAlgorithm(NewFirstFitAlgorithm()),
			prepare: func(algo Algorithm) {
				if mo, ok := algo.(*MultiObjectiveAlgorithm); ok {
					mo.SetWeights(0.3, 0.2, 0.2, 0.2, 0.1)
				}
			},
		},
		{
			name:      "MultiObjective-Knapsack",
			algorithm: NewMultiObjectiveAlgorithm(NewKnapsackAlgorithm()),
			prepare: func(algo Algorithm) {
				if mo, ok := algo.(*MultiObjectiveAlgorithm); ok {
					mo.SetWeights(0.4, 0.2, 0.2, 0.1, 0.1)
				}
			},
		},
	}

	for _, algo := range algorithms {
		b.Run(algo.name, func(b *testing.B) {
			algo.prepare(algo.algorithm)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := algo.algorithm.Nest(ctx, req)
				if err != nil {
					b.Fatalf("%s algorithm error: %v", algo.name, err)
				}
			}
		})
	}
}

// BenchmarkLargeScenario tests algorithms with larger data sets
func BenchmarkLargeScenario(b *testing.B) {
	testData := getComplexTestData() // Use complex data for larger scenario
	ctx := context.Background()
	req := &Request{
		Profiles:   testData.Profiles,
		StockBar:   testData.StockBar,
		Remnants:   testData.Remnants,
		Spacing:    10.0,
		Parameters: make(map[string]any),
	}

	algorithms := []struct {
		name      string
		algorithm Algorithm
		prepare   func(Algorithm)
	}{
		{
			name:      "FirstFitDecreasing-Large",
			algorithm: NewFirstFitDecreasingAlgorithm(),
			prepare:   func(algo Algorithm) {},
		},
		{
			name:      "BestFitDecreasing-Large",
			algorithm: NewBestFitDecreasingAlgorithm(),
			prepare:   func(algo Algorithm) {},
		},
		{
			name:      "Knapsack-Large",
			algorithm: NewKnapsackAlgorithm(),
			prepare:   func(algo Algorithm) {},
		},
	}

	for _, algo := range algorithms {
		b.Run(algo.name, func(b *testing.B) {
			algo.prepare(algo.algorithm)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := algo.algorithm.Nest(ctx, req)
				if err != nil {
					b.Fatalf("%s algorithm error: %v", algo.name, err)
				}
			}
		})
	}
}

// TestAlgorithmQuality compares quality metrics across algorithms
func TestAlgorithmQuality(t *testing.T) {
	testData := getSimpleTestData()
	ctx := context.Background()
	req := &Request{
		Profiles:   testData.Profiles,
		StockBar:   testData.StockBar,
		Remnants:   testData.Remnants,
		Spacing:    10.0,
		Parameters: make(map[string]any),
	}

	algorithms := []struct {
		name      string
		algorithm Algorithm
		prepare   func(Algorithm)
	}{
		{
			name:      "FirstFit",
			algorithm: NewFirstFitAlgorithm(),
			prepare:   func(algo Algorithm) {},
		},
		{
			name:      "FirstFitDecreasing",
			algorithm: NewFirstFitDecreasingAlgorithm(),
			prepare:   func(algo Algorithm) {},
		},
		{
			name:      "BestFitDecreasing",
			algorithm: NewBestFitDecreasingAlgorithm(),
			prepare:   func(algo Algorithm) {},
		},
		{
			name:      "Knapsack",
			algorithm: NewKnapsackAlgorithm(),
			prepare:   func(algo Algorithm) {},
		},
	}

	results := make(map[string]*Result)

	for _, algo := range algorithms {
		algo.prepare(algo.algorithm)

		result, err := algo.algorithm.Nest(ctx, req)
		if err != nil {
			t.Errorf("%s algorithm error: %v", algo.name, err)
			continue
		}

		results[algo.name] = result
	}

	// Compare quality metrics
	t.Run("utilization-comparison", func(t *testing.T) {
		for name, result := range results {
			if result.Statistics.MaterialUtilization < 50.0 {
				t.Logf("%s algorithm has low utilization: %.2f%%", name, result.Statistics.MaterialUtilization)
			}
		}

		// Check that improved algorithms generally have better utilization than basic first-fit
		ffUtil := results["FirstFit"].Statistics.MaterialUtilization
		for name, result := range results {
			if name != "FirstFit" && result.Statistics.MaterialUtilization < ffUtil-5.0 {
				t.Logf("%s algorithm has significantly lower utilization than FirstFit: %.2f%% vs %.2f%%",
					name, result.Statistics.MaterialUtilization, ffUtil)
			}
		}
	})

	t.Run("completion-rate", func(t *testing.T) {
		for name, result := range results {
			completionRate := float64(result.Statistics.NestedProfiles) / float64(len(testData.Profiles)) * 100
			if completionRate < 50.0 {
				t.Logf("%s algorithm has low completion rate: %.2f%%", name, completionRate)
			}
		}
	})
}
