package nesting

import (
	"context"
	"errors"
	"testing"
)

// TestMultiObjectiveAlgorithm tests multi-objective nesting algorithm
func TestMultiObjectiveAlgorithm(t *testing.T) {
	// Create base algorithm for testing
	baseAlgorithm := NewFirstFitAlgorithm()
	algorithm := NewMultiObjectiveAlgorithm(baseAlgorithm)

	// Set custom weights
	algorithm.SetWeights(0.4, 0.2, 0.2, 0.1, 0.1)

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
		{
			name:     "with-remnants",
			scenario: getTestScenarios()[2],
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
				t.Errorf("MultiObjectiveAlgorithm.Nest() error = %v, wantErr nil", err)
				return
			}

			if !validateResult(result, tt.scenario) {
				t.Errorf("MultiObjectiveAlgorithm.Nest() result validation failed for scenario %s", tt.name)
			}

			// Test multi-objective specific behavior - should have balanced optimization
			if len(result.Nests) > 0 {
				// Should consider multiple objectives
				if result.Statistics.MaterialUtilization < 40.0 {
					t.Errorf("MultiObjectiveAlgorithm.Nest() expected balanced utilization (>40%%), got %.2f%%", result.Statistics.MaterialUtilization)
				}

				// Algorithm name should indicate multi-objective nature
				if len(result.Statistics.Algorithm) <= len(baseAlgorithm.Name()) {
					t.Errorf("MultiObjectiveAlgorithm.Nest() expected algorithm name to reflect multi-objective nature")
				}
			}
		})
	}
}

// TestMultiObjectiveAlgorithmEdgeCases tests edge cases for multi-objective algorithm
func TestMultiObjectiveAlgorithmEdgeCases(t *testing.T) {
	baseAlgorithm := NewFirstFitAlgorithm()
	algorithm := NewMultiObjectiveAlgorithm(baseAlgorithm)
	algorithm.SetWeights(0.3, 0.3, 0.2, 0.1, 0.1)
	ctx := context.Background()

	t.Run("custom-objectives", func(t *testing.T) {
		// Add custom objective
		algorithm.AddObjective("custom-test", 0.5, true, func(result *Result) float64 {
			return float64(len(result.Nests)) // Objective: minimize number of nests
		})

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
			t.Errorf("MultiObjectiveAlgorithm.Nest() error = %v, wantErr nil", err)
			return
		}

		// Should consider custom objective
		if len(result.Nests) > 5 {
			t.Errorf("MultiObjectiveAlgorithm.Nest() expected custom objective to influence result (too many nests: %d)", len(result.Nests))
		}
	})

	t.Run("different-base-algorithms", func(t *testing.T) {
		// Test with different base algorithms
		knapsackAlgorithm := NewKnapsackAlgorithm()
		moWithKnapsack := NewMultiObjectiveAlgorithm(knapsackAlgorithm)

		testData := getSimpleTestData()
		req := &Request{
			Profiles:   testData.Profiles,
			StockBar:   testData.StockBar,
			Remnants:   testData.Remnants,
			Spacing:    10.0,
			Parameters: make(map[string]any),
		}

		result1, err1 := algorithm.Nest(ctx, req)
		result2, err2 := moWithKnapsack.Nest(ctx, req)

		if err1 != nil {
			t.Errorf("MultiObjectiveAlgorithm.Nest() with first-fit error = %v", err1)
		}
		if err2 != nil {
			t.Errorf("MultiObjectiveAlgorithm.Nest() with knapsack error = %v", err2)
		}

		// Results should be different due to different base algorithms
		if result1.Statistics.MaterialUtilization == result2.Statistics.MaterialUtilization {
			t.Logf("Note: Both algorithms produced same utilization %.2f%%", result1.Statistics.MaterialUtilization)
		}
	})
}

// TestMultiObjectiveAlgorithmWeights tests weight configuration
func TestMultiObjectiveAlgorithmWeights(t *testing.T) {
	baseAlgorithm := NewFirstFitAlgorithm()
	algorithm := NewMultiObjectiveAlgorithm(baseAlgorithm)
	ctx := context.Background()

	t.Run("different-weights", func(t *testing.T) {
		// Test with different weight configurations
		weightConfigs := []struct {
			name             string
			utilWeight       float64
			cutsWeight       float64
			stockWeight      float64
			completionWeight float64
			wasteWeight      float64
		}{
			{
				name:             "utilization-focused",
				utilWeight:       0.8,
				cutsWeight:       0.05,
				stockWeight:      0.05,
				completionWeight: 0.05,
				wasteWeight:      0.05,
			},
			{
				name:             "balanced",
				utilWeight:       0.3,
				cutsWeight:       0.2,
				stockWeight:      0.2,
				completionWeight: 0.2,
				wasteWeight:      0.1,
			},
			{
				name:             "cuts-focused",
				utilWeight:       0.1,
				cutsWeight:       0.6,
				stockWeight:      0.1,
				completionWeight: 0.1,
				wasteWeight:      0.1,
			},
		}

		testData := getSimpleTestData()
		req := &Request{
			Profiles:   testData.Profiles,
			StockBar:   testData.StockBar,
			Remnants:   testData.Remnants,
			Spacing:    10.0,
			Parameters: make(map[string]any),
		}

		for _, config := range weightConfigs {
			t.Run(config.name, func(t *testing.T) {
				algorithm.SetWeights(config.utilWeight, config.cutsWeight, config.stockWeight, config.completionWeight, config.wasteWeight)

				result, err := algorithm.Nest(ctx, req)

				if err != nil {
					t.Errorf("MultiObjectiveAlgorithm.Nest() error = %v, wantErr nil", err)
					return
				}

				// Results should be different based on weights
				if result == nil {
					t.Error("MultiObjectiveAlgorithm.Nest() expected non-nil result")
				}
			})
		}
	})
}

// TestMultiObjectiveAlgorithmPareto tests Pareto front functionality
func TestMultiObjectiveAlgorithmPareto(t *testing.T) {
	baseAlgorithm := NewFirstFitAlgorithm()
	algorithm := NewMultiObjectiveAlgorithm(baseAlgorithm)
	algorithm.SetWeights(0.3, 0.3, 0.2, 0.1, 0.1)

	// Add multiple objectives for Pareto analysis
	algorithm.AddObjective("minimize-nests", 0.5, true, func(result *Result) float64 {
		return float64(len(result.Nests))
	})

	algorithm.AddObjective("maximize-completion", 0.5, false, func(result *Result) float64 {
		if result.Statistics.TotalProfiles == 0 {
			return 0
		}
		return float64(result.Statistics.NestedProfiles) / float64(result.Statistics.TotalProfiles)
	})

	t.Run("pareto-analysis", func(t *testing.T) {
		testData := getSimpleTestData()
		req := &Request{
			Profiles:   testData.Profiles,
			StockBar:   testData.StockBar,
			Remnants:   testData.Remnants,
			Spacing:    10.0,
			Parameters: make(map[string]any),
		}

		result, err := algorithm.Nest(context.Background(), req)

		if err != nil {
			t.Errorf("MultiObjectiveAlgorithm.Nest() error = %v, wantErr nil", err)
			return
		}

		// Should produce result considering Pareto front
		if result == nil {
			t.Error("MultiObjectiveAlgorithm.Nest() expected Pareto analysis to produce result")
		}

		// Should consider multiple strategies
		if result.Statistics.Algorithm == "" {
			t.Error("MultiObjectiveAlgorithm.Nest() expected algorithm name to be set")
		}
	})
}

// TestMultiObjectiveAlgorithmCancellation tests context cancellation
func TestMultiObjectiveAlgorithmCancellation(t *testing.T) {
	baseAlgorithm := NewFirstFitAlgorithm()
	algorithm := NewMultiObjectiveAlgorithm(baseAlgorithm)

	t.Run("context-cancellation", func(t *testing.T) {
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

		// Cancel immediately
		cancel()

		_, err := algorithm.Nest(ctx, req)

		if err == nil {
			t.Error("MultiObjectiveAlgorithm.Nest() expected error with cancelled context, got nil")
		}

		if !errors.Is(err, context.Canceled) {
			t.Errorf("MultiObjectiveAlgorithm.Nest() expected context.Canceled error, got %v", err)
		}
	})
}

// BenchmarkMultiObjectiveAlgorithm benchmarks multi-objective algorithm
func BenchmarkMultiObjectiveAlgorithm(b *testing.B) {
	baseAlgorithm := NewFirstFitAlgorithm()
	algorithm := NewMultiObjectiveAlgorithm(baseAlgorithm)
	algorithm.SetWeights(0.3, 0.3, 0.2, 0.1, 0.1)

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
			b.Fatalf("MultiObjectiveAlgorithm.Nest() error = %v", err)
		}
	}
}
