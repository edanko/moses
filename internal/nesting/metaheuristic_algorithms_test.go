package nesting

import (
	"context"
	"testing"
)

// TestGeneticAlgorithm tests genetic nesting algorithm
func TestGeneticAlgorithm(t *testing.T) {
	algorithm := NewGeneticAlgorithm()

	// Use smaller parameters for testing to reduce test time
	algorithm.SetParameters(10, 20, 0.1, 0.2)

	tests := []struct {
		name     string
		scenario TestScenario
	}{
		{
			name:     "simple-case",
			scenario: getTestScenarios()[0],
		},
		{
			name:     "multiple-stocks",
			scenario: getTestScenarios()[1],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			req := &Request{
				Profiles:   tt.scenario.Profiles,
				StockBar:   tt.scenario.StockBar,
				Remnants:   tt.scenario.Remnants,
				Spacing:    tt.scenario.Spacing,
				Parameters: make(map[string]any),
			}

			result, err := algorithm.Nest(ctx, req)

			if err != nil {
				t.Errorf("GeneticAlgorithm.Nest() error = %v, wantErr nil", err)
				return
			}

			if !validateResult(result, tt.scenario) {
				t.Errorf("GeneticAlgorithm.Nest() result validation failed for scenario %s", tt.name)
			}

			// Test genetic algorithm specific behavior
			if len(result.Nests) > 0 {
				// Should have reasonable utilization due to optimization
				if result.Statistics.MaterialUtilization < 65.0 {
					t.Errorf("GeneticAlgorithm.Nest() expected good utilization (>65%%), got %.2f%%", result.Statistics.MaterialUtilization)
				}
			}
		})
	}
}

// TestGeneticAlgorithmEdgeCases tests edge cases for genetic algorithm
func TestGeneticAlgorithmEdgeCases(t *testing.T) {
	algorithm := NewGeneticAlgorithm()
	algorithm.SetParameters(5, 10, 0.1, 0.2) // Small parameters for testing

	t.Run("cancellation", func(t *testing.T) {
		testData := getComplexTestData()
		req := &Request{
			Profiles:   testData.Profiles,
			StockBar:   testData.StockBar,
			Remnants:   testData.Remnants,
			Spacing:    10.0,
			Parameters: make(map[string]any),
		}

		// Create context with cancellation
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		_, err := algorithm.Nest(ctx, req)

		if err == nil {
			t.Error("GeneticAlgorithm.Nest() expected error with cancelled context, got nil")
		}
	})
}

// TestSimulatedAnnealingAlgorithm tests simulated annealing nesting algorithm
func TestSimulatedAnnealingAlgorithm(t *testing.T) {
	algorithm := NewSimulatedAnnealingAlgorithm()

	// Use smaller parameters for testing
	algorithm.SetParameters(100.0, 0.9, 1.0, 10)

	tests := []struct {
		name     string
		scenario TestScenario
	}{
		{
			name:     "simple-case",
			scenario: getTestScenarios()[0],
		},
		{
			name:     "multiple-stocks",
			scenario: getTestScenarios()[1],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			req := &Request{
				Profiles:   tt.scenario.Profiles,
				StockBar:   tt.scenario.StockBar,
				Remnants:   tt.scenario.Remnants,
				Spacing:    tt.scenario.Spacing,
				Parameters: make(map[string]any),
			}

			result, err := algorithm.Nest(ctx, req)

			if err != nil {
				t.Errorf("SimulatedAnnealingAlgorithm.Nest() error = %v, wantErr nil", err)
				return
			}

			if !validateResult(result, tt.scenario) {
				t.Errorf("SimulatedAnnealingAlgorithm.Nest() result validation failed for scenario %s", tt.name)
			}

			// Test simulated annealing specific behavior
			if len(result.Nests) > 0 {
				// Should have good optimization results
				if result.Statistics.MaterialUtilization < 65.0 {
					t.Errorf("SimulatedAnnealingAlgorithm.Nest() expected good utilization (>65%%), got %.2f%%", result.Statistics.MaterialUtilization)
				}
			}
		})
	}
}

// TestSimulatedAnnealingAlgorithmEdgeCases tests edge cases for simulated annealing algorithm
func TestSimulatedAnnealingAlgorithmEdgeCases(t *testing.T) {
	algorithm := NewSimulatedAnnealingAlgorithm()
	algorithm.SetParameters(50.0, 0.9, 1.0, 5) // Small parameters for testing
	ctx := context.Background()

	t.Run("temperature-schedule", func(t *testing.T) {
		testData := getSimpleTestData()
		req := &Request{
			Profiles:   testData.Profiles,
			StockBar:   testData.StockBar,
			Remnants:   testData.Remnants,
			Spacing:    10.0,
			Parameters: make(map[string]any),
		}

		result, err := algorithm.Nest(ctx, req)

		if err != nil {
			t.Errorf("SimulatedAnnealingAlgorithm.Nest() error = %v, wantErr nil", err)
			return
		}

		// Should handle temperature scheduling properly
		if result.Statistics.NestedProfiles > len(testData.Profiles) {
			t.Errorf("SimulatedAnnealingAlgorithm.Nest() nested profiles exceed total profiles")
		}
	})
}

// TestAntColonyOptimizationAlgorithm tests ant colony optimization nesting algorithm
func TestAntColonyOptimizationAlgorithm(t *testing.T) {
	algorithm := NewAntColonyOptimizationAlgorithm()

	// Use smaller parameters for testing
	algorithm.SetParameters(5, 10, 1.0, 2.0, 0.8, 0.1)

	tests := []struct {
		name     string
		scenario TestScenario
	}{
		{
			name:     "simple-case",
			scenario: getTestScenarios()[0],
		},
		{
			name:     "multiple-stocks",
			scenario: getTestScenarios()[1],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			req := &Request{
				Profiles:   tt.scenario.Profiles,
				StockBar:   tt.scenario.StockBar,
				Remnants:   tt.scenario.Remnants,
				Spacing:    tt.scenario.Spacing,
				Parameters: make(map[string]any),
			}

			result, err := algorithm.Nest(ctx, req)

			if err != nil {
				t.Errorf("AntColonyOptimizationAlgorithm.Nest() error = %v, wantErr nil", err)
				return
			}

			if !validateResult(result, tt.scenario) {
				t.Errorf("AntColonyOptimizationAlgorithm.Nest() result validation failed for scenario %s", tt.name)
			}

			// Test ant colony optimization specific behavior
			if len(result.Nests) > 0 {
				// Should achieve good optimization
				if result.Statistics.MaterialUtilization < 60.0 {
					t.Errorf("AntColonyOptimizationAlgorithm.Nest() expected reasonable utilization (>60%%), got %.2f%%", result.Statistics.MaterialUtilization)
				}
			}
		})
	}
}

// TestAntColonyOptimizationAlgorithmEdgeCases tests edge cases for ant colony optimization algorithm
func TestAntColonyOptimizationAlgorithmEdgeCases(t *testing.T) {
	algorithm := NewAntColonyOptimizationAlgorithm()
	algorithm.SetParameters(3, 5, 1.0, 2.0, 0.8, 0.1) // numAnts, numIterations, alpha, beta, q0, tau0
	ctx := context.Background()

	t.Run("pheromone-evaporation", func(t *testing.T) {
		testData := getSimpleTestData()
		req := &Request{
			Profiles:   testData.Profiles,
			StockBar:   testData.StockBar,
			Remnants:   testData.Remnants,
			Spacing:    10.0,
			Parameters: make(map[string]any),
		}

		result, err := algorithm.Nest(ctx, req)

		if err != nil {
			t.Errorf("AntColonyOptimizationAlgorithm.Nest() error = %v, wantErr nil", err)
			return
		}

		// Should handle pheromone evaporation properly
		if result.Statistics.NestedProfiles > len(testData.Profiles) {
			t.Errorf("AntColonyOptimizationAlgorithm.Nest() nested profiles exceed total profiles")
		}
	})
}

// BenchmarkGeneticAlgorithm benchmarks genetic algorithm
func BenchmarkGeneticAlgorithm(b *testing.B) {
	algorithm := NewGeneticAlgorithm()
	algorithm.SetParameters(20, 50, 0.1, 0.2) // Moderate parameters for benchmarking
	testData := getSimpleTestData()
	ctx := context.Background()

	req := &Request{
		Profiles:   testData.Profiles,
		StockBar:   testData.StockBar,
		Remnants:   testData.Remnants,
		Spacing:    10.0,
		Parameters: make(map[string]any),
	}

	for b.Loop() {
		_, err := algorithm.Nest(ctx, req)
		if err != nil {
			b.Fatalf("GeneticAlgorithm.Nest() error = %v", err)
		}
	}
}

// BenchmarkSimulatedAnnealingAlgorithm benchmarks simulated annealing algorithm
func BenchmarkSimulatedAnnealingAlgorithm(b *testing.B) {
	algorithm := NewSimulatedAnnealingAlgorithm()
	algorithm.SetParameters(200.0, 0.95, 1.0, 20) // Moderate parameters
	testData := getSimpleTestData()
	ctx := context.Background()

	req := &Request{
		Profiles:   testData.Profiles,
		StockBar:   testData.StockBar,
		Remnants:   testData.Remnants,
		Spacing:    10.0,
		Parameters: make(map[string]any),
	}

	for b.Loop() {
		_, err := algorithm.Nest(ctx, req)
		if err != nil {
			b.Fatalf("SimulatedAnnealingAlgorithm.Nest() error = %v", err)
		}
	}
}

// BenchmarkAntColonyOptimizationAlgorithm benchmarks ant colony optimization algorithm
func BenchmarkAntColonyOptimizationAlgorithm(b *testing.B) {
	algorithm := NewAntColonyOptimizationAlgorithm()
	algorithm.SetParameters(10, 20, 1.0, 2.0, 0.8, 0.1) // numAnts, numIterations, alpha, beta, q0, tau0
	testData := getSimpleTestData()
	ctx := context.Background()

	req := &Request{
		Profiles:   testData.Profiles,
		StockBar:   testData.StockBar,
		Remnants:   testData.Remnants,
		Spacing:    10.0,
		Parameters: make(map[string]any),
	}

	for b.Loop() {
		_, err := algorithm.Nest(ctx, req)
		if err != nil {
			b.Fatalf("AntColonyOptimizationAlgorithm.Nest() error = %v", err)
		}
	}
}
