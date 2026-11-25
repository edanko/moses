package nesting

import (
	"uuid"

	"github.com/edanko/moses/internal/domain"
)

func testProfileDesignation() domain.ProfileDesignation {
	profileType, err := domain.ParseProfileDesignation("RP100X50")
	if err != nil {
		panic(err)
	}
	return profileType
}

// TestData provides common test data structures
type TestData struct {
	Profiles []domain.Part
	StockBar []domain.StockBar
	Remnants []domain.StockBar
}

// TestScenario represents a test scenario
type TestScenario struct {
	Name            string
	Profiles        []domain.Part
	StockBar        []domain.StockBar
	Remnants        []domain.StockBar
	Spacing         float64
	ExpectedNests   int
	ExpectedUtilMin float64
	ExpectedUtilMax float64
}

// createTestProfile creates a test profile with basic parameters
func testID(value string) uuid.UUID {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.NewV4()
	}
	return id
}

func createTestProfile(id string, length float64) domain.Part {
	profileID := testID(id)
	lEnd := domain.NewEndCut("E11", map[string]float64{"X": 0.0, "Y": 0.0})
	rEnd := domain.NewEndCut("E11", map[string]float64{"X": 0.0, "Y": 0.0})

	return domain.NewPart(domain.PartParams{
		ID:          profileID,
		Project:     "test-project",
		Section:     "test-assembly-block",
		PositionNo:  "test-position",
		Designation: testProfileDesignation(),
		Grade:       "S235",
		Length:      length,
		Quantity:    1,
		LeftEndCut:  lEnd,
		RightEndCut: rEnd,
	})
}

// createTestStock creates a test stock piece
func createTestStock(id string, length float64) domain.StockBar {
	stockID := testID(id)

	return domain.NewStockBar(domain.StockBarParams{
		ID:          stockID,
		Project:     "test-project",
		Designation: testProfileDesignation(),
		Grade:       "S235",
		Length:      length,
		Type:        domain.StockTypeBar,
	})
}

// getSimpleTestData creates simple test data for basic testing
func getSimpleTestData() TestData {
	return TestData{
		Profiles: []domain.Part{
			createTestProfile("00000000-0000-0000-0000-000000000001", 1000.0),
			createTestProfile("00000000-0000-0000-0000-000000000002", 800.0),
			createTestProfile("00000000-0000-0000-0000-000000000003", 1200.0),
			createTestProfile("00000000-0000-0000-0000-000000000004", 600.0),
		},
		StockBar: []domain.StockBar{
			createTestStock("00000000-0000-0000-0000-000000000101", 3000.0),
			createTestStock("00000000-0000-0000-0000-000000000102", 2000.0),
		},
		Remnants: []domain.StockBar{
			createTestStock("00000000-0000-0000-0000-000000000201", 1500.0),
		},
	}
}

// getComplexTestData creates complex test data for advanced testing
func getComplexTestData() TestData {
	profiles := make([]domain.Part, 0)
	for i := range 20 {
		length := float64(500 + i*50) // Variable lengths from 500 to 1450
		profileID := uuid.NewV4()
		lEnd := domain.NewEndCut("E11", map[string]float64{"X": 0.0, "Y": 0.0})
		rEnd := domain.NewEndCut("E11", map[string]float64{"X": 0.0, "Y": 0.0})

		profiles = append(profiles, domain.NewPart(domain.PartParams{
			ID:          profileID,
			Project:     "complex-test-project",
			Section:     "complex-assembly-block",
			PositionNo:  "position-" + string(rune(i+'0')),
			Designation: testProfileDesignation(),
			Grade:       "S235",
			Length:      length,
			Quantity:    1,
			LeftEndCut:  lEnd,
			RightEndCut: rEnd,
		}))
	}

	stock := []domain.StockBar{
		createTestStock("00000000-0000-0000-0000-000000000301", 6000.0),
		createTestStock("00000000-0000-0000-0000-000000000302", 6000.0),
		createTestStock("00000000-0000-0000-0000-000000000303", 4000.0),
	}

	remnants := []domain.StockBar{
		createTestStock("00000000-0000-0000-0000-000000000401", 2500.0),
		createTestStock("00000000-0000-0000-0000-000000000402", 1800.0),
	}

	return TestData{
		Profiles: profiles,
		StockBar: stock,
		Remnants: remnants,
	}
}

// getEdgeCaseTestData creates edge case test data
func getEdgeCaseTestData() TestData {
	return TestData{
		Profiles: []domain.Part{
			createTestProfile("00000000-0000-0000-0000-000000000501", 100.0),  // Very small
			createTestProfile("00000000-0000-0000-0000-000000000502", 5000.0), // Very large
			createTestProfile("00000000-0000-0000-0000-000000000503", 2999.0), // Almost fits in 3000
			createTestProfile("00000000-0000-0000-0000-000000000504", 3001.0), // Too large for 3000
		},
		StockBar: []domain.StockBar{
			createTestStock("00000000-0000-0000-0000-000000000501", 3000.0),
		},
		Remnants: []domain.StockBar{},
	}
}

// getNoFitTestData creates test data where profiles won't fit
func getNoFitTestData() TestData {
	return TestData{
		Profiles: []domain.Part{
			createTestProfile("00000000-0000-0000-0000-000000000601", 4000.0),
			createTestProfile("00000000-0000-0000-0000-000000000602", 4500.0),
		},
		StockBar: []domain.StockBar{
			createTestStock("00000000-0000-0000-0000-000000000601", 3000.0),
		},
		Remnants: []domain.StockBar{},
	}
}

// getZeroProfilesTestData creates test data with no profiles
func getZeroProfilesTestData() TestData {
	return TestData{
		Profiles: []domain.Part{},
		StockBar: []domain.StockBar{
			createTestStock("00000000-0000-0000-0000-000000000701", 3000.0),
		},
		Remnants: []domain.StockBar{},
	}
}

// getZeroStockTestData creates test data with no stock
func getZeroStockTestData() TestData {
	return TestData{
		Profiles: []domain.Part{
			createTestProfile("00000000-0000-0000-0000-000000000801", 1000.0),
		},
		StockBar: []domain.StockBar{},
		Remnants: []domain.StockBar{},
	}
}

// getTestScenarios returns a list of predefined test scenarios
func getTestScenarios() []TestScenario {
	return []TestScenario{
		{
			Name: "simple-case",
			Profiles: []domain.Part{
				createTestProfile("00000000-0000-0000-0000-000000000001", 1000.0),
				createTestProfile("00000000-0000-0000-0000-000000000002", 1500.0),
			},
			StockBar: []domain.StockBar{
				createTestStock("00000000-0000-0000-0000-000000000101", 3000.0),
			},
			Remnants:        []domain.StockBar{},
			Spacing:         10.0,
			ExpectedNests:   1,
			ExpectedUtilMin: 80.0,
			ExpectedUtilMax: 90.0,
		},
		{
			Name: "multiple-stocks",
			Profiles: []domain.Part{
				createTestProfile("00000000-0000-0000-0000-000000000002", 1200.0),
				createTestProfile("00000000-0000-0000-0000-000000000003", 1800.0),
				createTestProfile("00000000-0000-0000-0000-000000000004", 1000.0),
			},
			StockBar: []domain.StockBar{
				createTestStock("00000000-0000-0000-0000-000000000102", 2000.0),
				createTestStock("00000000-0000-0000-0000-000000000103", 2000.0),
			},
			Remnants:        []domain.StockBar{},
			Spacing:         5.0,
			ExpectedNests:   2,
			ExpectedUtilMin: 75.0,
			ExpectedUtilMax: 100.0,
		},
		{
			Name: "with-remnants",
			Profiles: []domain.Part{
				createTestProfile("00000000-0000-0000-0000-000000000005", 800.0),
				createTestProfile("00000000-0000-0000-0000-000000000006", 600.0),
			},
			StockBar: []domain.StockBar{
				createTestStock("00000000-0000-0000-0000-000000000104", 2000.0),
			},
			Remnants: []domain.StockBar{
				createTestStock("00000000-0000-0000-0000-000000000201", 1000.0),
			},
			Spacing:         0.0,
			ExpectedNests:   2,
			ExpectedUtilMin: 40.0,
			ExpectedUtilMax: 100.0,
		},
	}
}

// calculateExpectedUtilization calculates expected utilization for test validation

// validateResult validates nesting result against expected values.
func validateResult(result *Result, scenario TestScenario) bool {
	if result == nil {
		return false
	}
	if len(result.Nests) == 0 && len(scenario.Profiles) > 0 {
		return false
	}
	utilization := result.Statistics.MaterialUtilization
	if utilization < scenario.ExpectedUtilMin || utilization > scenario.ExpectedUtilMax {
		return false
	}
	if result.Statistics.NestedProfiles > len(scenario.Profiles) {
		return false
	}
	return result.Statistics.NestedProfiles+result.Statistics.UnnestedProfiles == len(scenario.Profiles)
}
