package nesting

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/edanko/moses/internal/domain"
)

// TestBranchAndBoundAlgorithm tests branch and bound nesting algorithm
func TestBranchAndBoundAlgorithm(t *testing.T) {
	algorithm := NewBranchAndBoundAlgorithm()

	// Use small parameters for testing to avoid timeouts
	algorithm.SetParameters(1*time.Second, 100)

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
				t.Errorf("BranchAndBoundAlgorithm.Nest() error = %v, wantErr nil", err)
				return
			}

			if !validateResult(result, tt.scenario) {
				t.Errorf("BranchAndBoundAlgorithm.Nest() result validation failed for scenario %s", tt.name)
			}

			// Test branch and bound specific behavior - should have high utilization
			if len(result.Nests) > 0 {
				if result.Statistics.MaterialUtilization < 70.0 {
					t.Errorf("BranchAndBoundAlgorithm.Nest() expected high utilization (>70%%), got %.2f%%", result.Statistics.MaterialUtilization)
				}
			}
		})
	}
}

// TestBranchAndBoundAlgorithmEdgeCases tests edge cases for branch and bound algorithm
func TestBranchAndBoundAlgorithmEdgeCases(t *testing.T) {
	algorithm := NewBranchAndBoundAlgorithm()
	algorithm.SetParameters(500*time.Millisecond, 50) // Very short timeout for testing
	ctx := context.Background()

	t.Run("time-limit", func(t *testing.T) {
		testData := getComplexTestData()
		req := &Request{
			Profiles:   testData.Profiles,
			StockBar:   testData.StockBar,
			Remnants:   testData.Remnants,
			Spacing:    10.0,
			Parameters: make(map[string]any),
		}

		_, err := algorithm.Nest(ctx, req)

		// Should complete quickly with simple scenario
		if err != nil {
			t.Errorf("BranchAndBoundAlgorithm.Nest() error = %v, wantErr nil", err)
		}
	})

	t.Run("solution-limit", func(t *testing.T) {
		testData := getComplexTestData()
		req := &Request{
			Profiles:   testData.Profiles,
			StockBar:   testData.StockBar,
			Remnants:   testData.Remnants,
			Spacing:    10.0,
			Parameters: make(map[string]any),
		}

		// Reset with very small solution limit
		algorithm.SetParameters(100*time.Millisecond, 1)

		result, err := algorithm.Nest(ctx, req)

		if err != nil {
			t.Errorf("BranchAndBoundAlgorithm.Nest() error = %v, wantErr nil", err)
		}

		// Should complete quickly with solution limit
		if result == nil {
			t.Error("BranchAndBoundAlgorithm.Nest() expected result even with solution limit")
		}
	})
}

// TestBranchAndBoundAlgorithmPruning tests pruning functionality
func TestBranchAndBoundAlgorithmPruning(t *testing.T) {
	algorithm := NewBranchAndBoundAlgorithm()
	algorithm.SetParameters(2*time.Second, 1000)
	ctx := context.Background()

	t.Run("pruning-efficiency", func(t *testing.T) {
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
			t.Errorf("BranchAndBoundAlgorithm.Nest() error = %v, wantErr nil", err)
			return
		}

		// Should demonstrate pruning efficiency
		if result.Statistics.ExecutionTime == "" {
			t.Error("BranchAndBoundAlgorithm.Nest() expected execution time to be recorded")
		}

		// Should produce optimal or near-optimal solution
		if result.Statistics.MaterialUtilization < 60.0 {
			t.Errorf("BranchAndBoundAlgorithm.Nest() expected good utilization (>60%%), got %.2f%%", result.Statistics.MaterialUtilization)
		}
	})
}

// TestBranchAndBoundAlgorithmCancellation tests context cancellation
func TestBranchAndBoundAlgorithmCancellation(t *testing.T) {
	algorithm := NewBranchAndBoundAlgorithm()
	algorithm.SetParameters(1*time.Second, 100)

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
			t.Error("BranchAndBoundAlgorithm.Nest() expected error with cancelled context, got nil")
		}

		if !errors.Is(err, context.Canceled) {
			t.Errorf("BranchAndBoundAlgorithm.Nest() expected context.Canceled error, got %v", err)
		}
	})
}

// TestBranchAndBoundAlgorithmOptimality tests if algorithm produces optimal solutions
func TestBranchAndBoundAlgorithmOptimality(t *testing.T) {
	algorithm := NewBranchAndBoundAlgorithm()
	algorithm.SetParameters(2*time.Second, 1000) // Generous limits for testing
	ctx := context.Background()

	t.Run("optimal-solution", func(t *testing.T) {
		// Simple case where optimal solution is known
		profiles := []domain.Part{
			createTestProfile("00000000-0000-0000-0000-000000000001", 1000.0),
			createTestProfile("00000000-0000-0000-0000-000000000002", 1000.0),
		}
		stock := []domain.StockBar{
			createTestStock("00000000-0000-0000-0000-000000000101", 2100.0),
		}

		req := &Request{
			Profiles:   profiles,
			StockBar:   stock,
			Remnants:   []domain.StockBar{},
			Spacing:    0.0, // No spacing for simple case
			Parameters: make(map[string]any),
		}

		result, err := algorithm.Nest(ctx, req)

		if err != nil {
			t.Errorf("BranchAndBoundAlgorithm.Nest() error = %v, wantErr nil", err)
			return
		}

		// In this simple case, both profiles should fit in one stock
		if len(result.Nests) != 1 {
			t.Errorf("BranchAndBoundAlgorithm.Nest() expected 1 nest, got %d", len(result.Nests))
		}

		if len(result.Nests) > 0 && len(result.Nests[0].Parts()) != 2 {
			t.Errorf("BranchAndBoundAlgorithm.Nest() expected 2 profiles in nest, got %d", len(result.Nests[0].Parts()))
		}

		// Should have very high utilization
		if result.Statistics.MaterialUtilization < 90.0 {
			t.Errorf("BranchAndBoundAlgorithm.Nest() expected optimal utilization (>90%%), got %.2f%%", result.Statistics.MaterialUtilization)
		}
	})
}

// BenchmarkBranchAndBoundAlgorithm benchmarks branch and bound algorithm
func BenchmarkBranchAndBoundAlgorithm(b *testing.B) {
	algorithm := NewBranchAndBoundAlgorithm()
	algorithm.SetParameters(5*time.Second, 1000) // Reasonable limits for benchmarking
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
			b.Fatalf("BranchAndBoundAlgorithm.Nest() error = %v", err)
		}
	}
}
