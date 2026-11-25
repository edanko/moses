package nesting

import (
	"context"

	"github.com/edanko/moses/internal/domain"
)

// Algorithm defines the interface for nesting algorithms
type Algorithm interface {
	// Name returns the algorithm name
	Name() string

	// Nest performs the nesting operation
	Nest(ctx context.Context, req *Request) (*Result, error)
}

// Request contains all data needed for nesting
type Request struct {
	Profiles   []domain.Part
	StockBar   []domain.StockBar
	Remnants   []domain.StockBar
	Spacing    float64        // default spacing when profile-specific not available
	Parameters map[string]any // algorithm-specific parameters
}

// Result contains the nesting output
type Result struct {
	Nests      []domain.Nest
	Unnested   []domain.Part // profiles that couldn't be nested
	Statistics Statistics
}

// Statistics provides nesting performance metrics
type Statistics struct {
	TotalProfiles       int
	NestedProfiles      int
	UnnestedProfiles    int
	MaterialUtilization float64 // percentage of stock material used
	TotalWaste          float64
	Algorithm           string
	ExecutionTime       string
}
