package nesting

import (
	"context"
	"errors"
	"testing"
)

// TestFirstFitAlgorithm tests the first-fit nesting algorithm
func TestFirstFitAlgorithm(t *testing.T) {
	algorithm := NewFirstFitAlgorithm()

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
				t.Errorf("FirstFitAlgorithm.Nest() error = %v, wantErr nil", err)
				return
			}

			if !validateResult(result, tt.scenario) {
				t.Errorf("FirstFitAlgorithm.Nest() result validation failed for scenario %s", tt.name)
			}

			// Test specific first-fit behavior
			if len(result.Nests) > 0 {
				// First-fit should place profiles in first available stock
				if len(result.Nests) > len(tt.scenario.StockBar)+len(tt.scenario.Remnants) {
					t.Errorf("FirstFitAlgorithm.Nest() created more nests than available stock pieces")
				}
			}
		})
	}
}

// TestFirstFitAlgorithmEdgeCases tests edge cases for first-fit algorithm
func TestFirstFitAlgorithmEdgeCases(t *testing.T) {
	algorithm := NewFirstFitAlgorithm()
	ctx := context.Background()

	t.Run("no-profiles", func(t *testing.T) {
		testData := getZeroProfilesTestData()
		req := &Request{
			Profiles:   testData.Profiles,
			StockBar:   testData.StockBar,
			Remnants:   testData.Remnants,
			Spacing:    10.0,
			Parameters: make(map[string]any),
		}

		result, err := algorithm.Nest(ctx, req)

		if err != nil {
			t.Errorf("FirstFitAlgorithm.Nest() error = %v, wantErr nil", err)
			return
		}

		if len(result.Nests) != 0 {
			t.Errorf("FirstFitAlgorithm.Nest() expected 0 nests, got %d", len(result.Nests))
		}
	})

	t.Run("no-stock", func(t *testing.T) {
		testData := getZeroStockTestData()
		req := &Request{
			Profiles:   testData.Profiles,
			StockBar:   testData.StockBar,
			Remnants:   testData.Remnants,
			Spacing:    10.0,
			Parameters: make(map[string]any),
		}

		result, err := algorithm.Nest(ctx, req)

		if err != nil {
			t.Errorf("FirstFitAlgorithm.Nest() error = %v, wantErr nil", err)
			return
		}

		if len(result.Nests) != 0 {
			t.Errorf("FirstFitAlgorithm.Nest() expected 0 nests with no stock, got %d", len(result.Nests))
		}

		if result.Statistics.NestedProfiles != 0 {
			t.Errorf("FirstFitAlgorithm.Nest() expected 0 nested profiles with no stock, got %d", result.Statistics.NestedProfiles)
		}
	})

	t.Run("no-fit", func(t *testing.T) {
		testData := getNoFitTestData()
		req := &Request{
			Profiles:   testData.Profiles,
			StockBar:   testData.StockBar,
			Remnants:   testData.Remnants,
			Spacing:    10.0,
			Parameters: make(map[string]any),
		}

		result, err := algorithm.Nest(ctx, req)

		if err != nil {
			t.Errorf("FirstFitAlgorithm.Nest() error = %v, wantErr nil", err)
			return
		}

		if result.Statistics.NestedProfiles != 0 {
			t.Errorf("FirstFitAlgorithm.Nest() expected 0 nested profiles when none fit, got %d", result.Statistics.NestedProfiles)
		}
	})
}

// TestFirstFitAlgorithmCancellation tests context cancellation
func TestFirstFitAlgorithmCancellation(t *testing.T) {
	algorithm := NewFirstFitAlgorithm()

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())

	testData := getComplexTestData()
	req := &Request{
		Profiles:   testData.Profiles,
		StockBar:   testData.StockBar,
		Remnants:   testData.Remnants,
		Spacing:    10.0,
		Parameters: make(map[string]any),
	}

	// Cancel context immediately
	cancel()

	_, err := algorithm.Nest(ctx, req)

	if err == nil {
		t.Error("FirstFitAlgorithm.Nest() expected error with cancelled context, got nil")
	}

	if !errors.Is(err, context.Canceled) {
		t.Errorf("FirstFitAlgorithm.Nest() expected context.Canceled error, got %v", err)
	}
}

// TestKnapsackAlgorithm tests the knapsack nesting algorithm
func TestKnapsackAlgorithm(t *testing.T) {
	algorithm := NewKnapsackAlgorithm()

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
				t.Errorf("KnapsackAlgorithm.Nest() error = %v, wantErr nil", err)
				return
			}

			if !validateResult(result, tt.scenario) {
				t.Errorf("KnapsackAlgorithm.Nest() result validation failed for scenario %s", tt.name)
			}

			// Test knapsack-specific behavior - should have good utilization
			if result.Statistics.MaterialUtilization < 70.0 {
				t.Errorf("KnapsackAlgorithm.Nest() expected good utilization (>70%%), got %.2f%%", result.Statistics.MaterialUtilization)
			}
		})
	}
}

// TestKnapsackAlgorithmEdgeCases tests edge cases for knapsack algorithm
func TestKnapsackAlgorithmEdgeCases(t *testing.T) {
	algorithm := NewKnapsackAlgorithm()
	ctx := context.Background()

	t.Run("large-profiles", func(t *testing.T) {
		testData := getEdgeCaseTestData()
		req := &Request{
			Profiles:   testData.Profiles,
			StockBar:   testData.StockBar,
			Remnants:   testData.Remnants,
			Spacing:    10.0,
			Parameters: make(map[string]any),
		}

		result, err := algorithm.Nest(ctx, req)

		if err != nil {
			t.Errorf("KnapsackAlgorithm.Nest() error = %v, wantErr nil", err)
			return
		}

		// Should handle edge cases gracefully
		if result.Statistics.NestedProfiles > len(testData.Profiles) {
			t.Errorf("KnapsackAlgorithm.Nest() nested profiles exceed total profiles")
		}
	})
}

// BenchmarkFirstFitAlgorithm benchmarks the first-fit algorithm
func BenchmarkFirstFitAlgorithm(b *testing.B) {
	algorithm := NewFirstFitAlgorithm()
	testData := getComplexTestData()
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
			b.Fatalf("KnapsackAlgorithm.Nest() error = %v", err)
		}
	}
}

// BenchmarkKnapsackAlgorithm benchmarks the knapsack algorithm
func BenchmarkKnapsackAlgorithm(b *testing.B) {
	algorithm := NewKnapsackAlgorithm()
	testData := getComplexTestData()
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
			b.Fatalf("KnapsackAlgorithm.Nest() error = %v", err)
		}
	}
}
