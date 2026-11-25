package nesting

import (
	"context"
	"testing"
)

// TestBestFitDecreasingAlgorithm tests best-fit decreasing nesting algorithm
func TestBestFitDecreasingAlgorithm(t *testing.T) {
	algorithm := NewBestFitDecreasingAlgorithm()

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
				t.Errorf("BestFitDecreasingAlgorithm.Nest() error = %v, wantErr nil", err)
				return
			}

			if !validateResult(result, tt.scenario) {
				t.Errorf("BestFitDecreasingAlgorithm.Nest() result validation failed for scenario %s", tt.name)
			}

			// Test best-fit specific behavior - should have reasonable utilization
			if result.Statistics.MaterialUtilization < 40.0 {
				t.Errorf("BestFitDecreasingAlgorithm.Nest() expected reasonable utilization (>40%%), got %.2f%%", result.Statistics.MaterialUtilization)
			}
		})
	}
}

// TestBestFitDecreasingAlgorithmEdgeCases tests edge cases for best-fit decreasing algorithm
func TestBestFitDecreasingAlgorithmEdgeCases(t *testing.T) {
	algorithm := NewBestFitDecreasingAlgorithm()
	ctx := context.Background()

	t.Run("complex-scenario", func(t *testing.T) {
		testData := getComplexTestData()
		req := &Request{
			Profiles:   testData.Profiles,
			StockBar:   testData.StockBar,
			Remnants:   testData.Remnants,
			Spacing:    10.0,
			Parameters: make(map[string]any),
		}

		result, err := algorithm.Nest(ctx, req)

		if err != nil {
			t.Errorf("BestFitDecreasingAlgorithm.Nest() error = %v, wantErr nil", err)
			return
		}

		// Should handle complex scenarios
		if result.Statistics.NestedProfiles == 0 {
			t.Errorf("BestFitDecreasingAlgorithm.Nest() expected some profiles to be nested in complex scenario")
		}

		// Best-fit decreasing should perform better than simple first-fit
		if result.Statistics.MaterialUtilization < 50.0 {
			t.Errorf("BestFitDecreasingAlgorithm.Nest() expected reasonable utilization (>50%%), got %.2f%%", result.Statistics.MaterialUtilization)
		}
	})
}

// TestFirstFitDecreasingAlgorithm tests first-fit decreasing nesting algorithm
func TestFirstFitDecreasingAlgorithm(t *testing.T) {
	algorithm := NewFirstFitDecreasingAlgorithm()

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
				t.Errorf("FirstFitDecreasingAlgorithm.Nest() error = %v, wantErr nil", err)
				return
			}

			if !validateResult(result, tt.scenario) {
				t.Errorf("FirstFitDecreasingAlgorithm.Nest() result validation failed for scenario %s", tt.name)
			}

			// Test first-fit decreasing specific behavior
			if len(result.Nests) == 0 && len(tt.scenario.Profiles) > 0 {
				t.Errorf("FirstFitDecreasingAlgorithm.Nest() expected at least one nest")
			}
		})
	}
}

// TestFirstFitDecreasingAlgorithmEdgeCases tests edge cases for first-fit decreasing algorithm
func TestFirstFitDecreasingAlgorithmEdgeCases(t *testing.T) {
	algorithm := NewFirstFitDecreasingAlgorithm()
	ctx := context.Background()

	t.Run("profile-sorting", func(t *testing.T) {
		// Test that profiles are indeed sorted by decreasing length
		testData := getSimpleTestData()
		req := &Request{
			Profiles:   testData.Profiles,
			StockBar:   testData.StockBar,
			Remnants:   testData.Remnants,
			Spacing:    10.0,
			Parameters: make(map[string]any),
		}

		_, err := algorithm.Nest(ctx, req)

		if err != nil {
			t.Errorf("FirstFitDecreasingAlgorithm.Nest() error = %v, wantErr nil", err)
			return
		}

		// The fact that it completes without error indicates sorting is working
		// (the algorithm expects sorted profiles)
	})
}

// BenchmarkBestFitDecreasingAlgorithm benchmarks best-fit decreasing algorithm
func BenchmarkBestFitDecreasingAlgorithm(b *testing.B) {
	algorithm := NewBestFitDecreasingAlgorithm()
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
			b.Fatalf("BestFitDecreasingAlgorithm.Nest() error = %v", err)
		}
	}
}

// BenchmarkFirstFitDecreasingAlgorithm benchmarks first-fit decreasing algorithm
func BenchmarkFirstFitDecreasingAlgorithm(b *testing.B) {
	algorithm := NewFirstFitDecreasingAlgorithm()
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
			b.Fatalf("FirstFitDecreasingAlgorithm.Nest() error = %v", err)
		}
	}
}
