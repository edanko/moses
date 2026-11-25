package nesting

import (
	"context"
	"slices"
	"testing"
	"time"
	"uuid"

	"github.com/edanko/moses/internal/domain"
)

// TestIntegrationSystem tests complete system integration
func TestIntegrationSystem(t *testing.T) {
	t.Run("algorithm-registry", func(t *testing.T) {
		// Test that all algorithms are registered
		algorithms := ListAlgorithms()

		expectedAlgorithms := []string{
			"first-fit",
			"knapsack",
			"best-fit-decreasing",
			"first-fit-decreasing",
			"genetic-algorithm",
			"simulated-annealing",
			"ant-colony-optimization",
			"branch-and-bound",
			"multi-objective-first-fit",
			"multi-objective-knapsack",
		}

		for _, expected := range expectedAlgorithms {
			found := slices.Contains(algorithms, expected)
			if !found {
				t.Errorf("Expected algorithm %s not found in registry", expected)
			}
		}

		// Check that we have all expected algorithms
		if len(algorithms) < len(expectedAlgorithms) {
			t.Errorf("Expected at least %d algorithms, got %d", len(expectedAlgorithms), len(algorithms))
		}
	})

	t.Run("algorithm-retrieval", func(t *testing.T) {
		// Test retrieving algorithms from registry
		algorithms := []string{
			"first-fit",
			"knapsack",
			"best-fit-decreasing",
			"first-fit-decreasing",
			"genetic-algorithm",
			"simulated-annealing",
			"ant-colony-optimization",
			"branch-and-bound",
			"multi-objective-first-fit",
			"multi-objective-knapsack",
		}

		for _, name := range algorithms {
			algorithm, err := GetAlgorithm(name)
			if err != nil {
				t.Errorf("Failed to get algorithm %s: %v", name, err)
				continue
			}

			if algorithm.Name() != name {
				t.Errorf("Algorithm name mismatch: expected %s, got %s", name, algorithm.Name())
			}
		}
	})

	t.Run("algorithm-consistency", func(t *testing.T) {
		// Test that algorithms behave consistently
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
		}{
			{
				name:      "FirstFit",
				algorithm: NewFirstFitAlgorithm(),
			},
			{
				name:      "BestFitDecreasing",
				algorithm: NewBestFitDecreasingAlgorithm(),
			},
			{
				name:      "FirstFitDecreasing",
				algorithm: NewFirstFitDecreasingAlgorithm(),
			},
		}

		results := make(map[string]*Result)

		for _, algo := range algorithms {
			result, err := algo.algorithm.Nest(ctx, req)
			if err != nil {
				t.Errorf("%s algorithm error: %v", algo.name, err)
				continue
			}
			results[algo.name] = result
		}

		// Check consistency
		for name, result := range results {
			if result == nil {
				t.Errorf("%s algorithm returned nil result", name)
				continue
			}

			// Basic consistency checks
			if result.Statistics.TotalProfiles != len(testData.Profiles) {
				t.Errorf("%s algorithm: total profiles mismatch", name)
			}

			if result.Statistics.NestedProfiles > len(testData.Profiles) {
				t.Errorf("%s algorithm: nested profiles exceed total", name)
			}

			if result.Statistics.NestedProfiles+result.Statistics.UnnestedProfiles != len(testData.Profiles) {
				t.Errorf("%s algorithm: nested + unnested != total", name)
			}

			if result.Statistics.MaterialUtilization < 0 || result.Statistics.MaterialUtilization > 100 {
				t.Errorf("%s algorithm: invalid utilization %f", name, result.Statistics.MaterialUtilization)
			}
		}
	})
}

// TestIntegrationRealWorldScenario tests algorithms with realistic data
func TestIntegrationRealWorldScenario(t *testing.T) {
	t.Run("steel-cutting-scenario", func(t *testing.T) {
		// Create realistic steel cutting scenario
		profiles := []domain.Part{
			createTestProfile("00000000-0000-0000-0000-000000000001", 1200.0), // 12m profiles
			createTestProfile("00000000-0000-0000-0000-000000000002", 1500.0),
			createTestProfile("00000000-0000-0000-0000-000000000003", 1800.0),
			createTestProfile("00000000-0000-0000-0000-000000000004", 1000.0),
			createTestProfile("00000000-0000-0000-0000-000000000005", 900.0),
			createTestProfile("00000000-0000-0000-0000-000000000006", 1100.0),
			createTestProfile("00000000-0000-0000-0000-000000000007", 1300.0),
			createTestProfile("00000000-0000-0000-0000-000000000008", 800.0),
		}

		stock := []domain.StockBar{
			createTestStock("00000000-0000-0000-0000-000000000101", 6000.0), // 6m stock
			createTestStock("00000000-0000-0000-0000-000000000102", 6000.0),
			createTestStock("00000000-0000-0000-0000-000000000103", 4000.0), // 4m stock
		}

		remnants := []domain.StockBar{
			createTestStock("00000000-0000-0000-0000-000000000201", 2500.0), // 2.5m remnant
		}

		ctx := context.Background()
		req := &Request{
			Profiles:   profiles,
			StockBar:   stock,
			Remnants:   remnants,
			Spacing:    5.0, // 5cm spacing between cuts
			Parameters: make(map[string]any),
		}

		// Test with a subset of algorithms for comprehensive integration
		algorithms := []struct {
			name      string
			algorithm Algorithm
		}{
			{
				name:      "FirstFit",
				algorithm: NewFirstFitAlgorithm(),
			},
			{
				name:      "BestFitDecreasing",
				algorithm: NewBestFitDecreasingAlgorithm(),
			},
			{
				name:      "Knapsack",
				algorithm: NewKnapsackAlgorithm(),
			},
			{
				name:      "Genetic",
				algorithm: NewGeneticAlgorithm(),
			},
		}

		results := make(map[string]*Result)

		for _, algo := range algorithms {
			// Prepare algorithm if needed
			if ga, ok := algo.algorithm.(*GeneticAlgorithm); ok {
				ga.SetParameters(10, 20, 0.1, 0.2) // Small parameters for testing
			}

			result, err := algo.algorithm.Nest(ctx, req)
			if err != nil {
				t.Errorf("%s algorithm error: %v", algo.name, err)
				continue
			}
			results[algo.name] = result

			t.Logf("%s algorithm: %d nests, %.2f%% utilization, %d nested/%d total",
				algo.name,
				len(result.Nests),
				result.Statistics.MaterialUtilization,
				result.Statistics.NestedProfiles,
				len(profiles))
		}

		// Check that all algorithms produce reasonable results
		for name, result := range results {
			if result == nil {
				t.Errorf("%s algorithm returned nil result", name)
				continue
			}

			// Realistic scenario should have decent utilization
			if result.Statistics.MaterialUtilization < 40.0 {
				t.Errorf("%s algorithm has poor utilization (%.2f%%) for realistic scenario", name, result.Statistics.MaterialUtilization)
			}

			// Should nest most profiles
			completionRate := float64(result.Statistics.NestedProfiles) / float64(len(profiles)) * 100
			if completionRate < 60.0 {
				t.Errorf("%s algorithm has low completion rate (%.2f%%) for realistic scenario", name, completionRate)
			}
		}
	})
}

// TestIntegrationPerformance tests performance characteristics
func TestIntegrationPerformance(t *testing.T) {
	t.Run("performance-comparison", func(t *testing.T) {
		testData := getSimpleTestData()
		ctx := context.Background()
		req := &Request{
			Profiles:   testData.Profiles,
			StockBar:   testData.StockBar,
			Remnants:   testData.Remnants,
			Spacing:    10.0,
			Parameters: make(map[string]any),
		}

		// Test performance with time limits
		algorithms := []struct {
			name          string
			algorithm     Algorithm
			expectedMaxMs int64 // Maximum expected time in milliseconds
		}{
			{
				name:          "FirstFit",
				algorithm:     NewFirstFitAlgorithm(),
				expectedMaxMs: 100, // Should be very fast
			},
			{
				name:          "BestFitDecreasing",
				algorithm:     NewBestFitDecreasingAlgorithm(),
				expectedMaxMs: 100, // Should be very fast
			},
			{
				name:          "FirstFitDecreasing",
				algorithm:     NewFirstFitDecreasingAlgorithm(),
				expectedMaxMs: 100, // Should be very fast
			},
			{
				name:          "Knapsack",
				algorithm:     NewKnapsackAlgorithm(),
				expectedMaxMs: 500, // Should be fast to moderate
			},
		}

		for _, algo := range algorithms {
			start := time.Now()

			result, err := algo.algorithm.Nest(ctx, req)

			duration := time.Since(start)

			if err != nil {
				t.Errorf("%s algorithm error: %v", algo.name, err)
				continue
			}

			if result == nil {
				t.Errorf("%s algorithm returned nil result", algo.name)
				continue
			}

			// Check performance expectations
			if duration.Milliseconds() > algo.expectedMaxMs {
				t.Logf("%s algorithm exceeded expected time: %dms vs %dms (may be acceptable on this system)",
					algo.name, duration.Milliseconds(), algo.expectedMaxMs)
			}

			t.Logf("%s algorithm performance: %dms, %d nests, %.2f%% utilization",
				algo.name, duration.Milliseconds(),
				len(result.Nests),
				result.Statistics.MaterialUtilization)
		}
	})
}

// TestIntegrationStress tests algorithms with stress scenarios
func TestIntegrationStress(t *testing.T) {
	t.Run("large-dataset", func(t *testing.T) {
		// Create larger stress test dataset
		profiles := make([]domain.Part, 0)
		for i := range 50 {
			length := float64(500 + i*20) // Profiles from 500 to 1480
			profileID := uuid.NewV4()
			lEnd := domain.NewEndCut("E11", map[string]float64{"X": 0.0, "Y": 0.0})
			rEnd := domain.NewEndCut("E11", map[string]float64{"X": 0.0, "Y": 0.0})

			profiles = append(profiles, domain.NewPart(domain.PartParams{
				ID:          profileID,
				Project:     "stress-test",
				Section:     "stress-assembly-block",
				PositionNo:  "stress-position",
				Designation: testProfileDesignation(),
				Grade:       "S235",
				Length:      length,
				Quantity:    1,
				LeftEndCut:  lEnd,
				RightEndCut: rEnd,
			}))
		}

		stock := []domain.StockBar{
			createTestStock("stress-stock-1", 8000.0),
			createTestStock("stress-stock-2", 8000.0),
			createTestStock("stress-stock-3", 6000.0),
			createTestStock("stress-stock-4", 6000.0),
		}

		ctx := context.Background()
		req := &Request{
			Profiles:   profiles,
			StockBar:   stock,
			Remnants:   []domain.StockBar{},
			Spacing:    15.0,
			Parameters: make(map[string]any),
		}

		// Test with fast algorithms only for stress test
		algorithms := []struct {
			name      string
			algorithm Algorithm
		}{
			{
				name:      "FirstFitStress",
				algorithm: NewFirstFitAlgorithm(),
			},
			{
				name:      "BestFitDecreasingStress",
				algorithm: NewBestFitDecreasingAlgorithm(),
			},
		}

		for _, algo := range algorithms {
			start := time.Now()

			result, err := algo.algorithm.Nest(ctx, req)

			duration := time.Since(start)

			if err != nil {
				t.Errorf("%s algorithm error on stress test: %v", algo.name, err)
				continue
			}

			if result == nil {
				t.Errorf("%s algorithm returned nil result on stress test", algo.name)
				continue
			}

			// Stress test should complete in reasonable time
			if duration.Seconds() > 10 {
				t.Errorf("%s algorithm took too long on stress test: %.2fs", algo.name, duration.Seconds())
			}

			// Should handle large dataset
			if result.Statistics.NestedProfiles == 0 {
				t.Errorf("%s algorithm failed to nest any profiles in stress test", algo.name)
			}

			t.Logf("%s stress test: %d profiles -> %d nests in %.2fs, %.2f%% utilization",
				algo.name, len(profiles), len(result.Nests),
				duration.Seconds(), result.Statistics.MaterialUtilization)
		}
	})
}
