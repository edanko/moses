package nesting

import (
	"context"

	"github.com/edanko/moses/internal/domain"
	"github.com/edanko/moses/pkg/binpack"
)

// FirstFitDecreasingAlgorithm implements the first-fit decreasing bin packing algorithm
type FirstFitDecreasingAlgorithm struct{}

// NewFirstFitDecreasingAlgorithm creates a new first-fit decreasing algorithm
func NewFirstFitDecreasingAlgorithm() *FirstFitDecreasingAlgorithm {
	return &FirstFitDecreasingAlgorithm{}
}

// Name returns the algorithm name
func (f *FirstFitDecreasingAlgorithm) Name() string {
	return "first-fit-decreasing"
}

// Nest performs the nesting operation using first-fit decreasing algorithm
func (f *FirstFitDecreasingAlgorithm) Nest(ctx context.Context, req *Request) (*Result, error) {
	items := toBinpackItems(req.Profiles)
	bins := toBinpackBins(req.StockBar)

	if len(req.Remnants) > 0 {
		remnantBins := toBinpackBins(req.Remnants)
		bins = append(remnantBins, bins...)
	}

	// Sort items by length in descending order before packing
	binpack.SortItemsByLength(items)

	algorithm := binpack.NewFirstFitAlgorithm[domain.StockBar, domain.Part]()
	packingResult, err := algorithm.Pack(ctx, items, bins)
	if err != nil {
		return nil, err
	}

	nests, unnested := fromBinpackResult(packingResult)

	return &Result{
		Nests:      nests,
		Unnested:   unnested,
		Statistics: convertStatistics(packingResult.Statistics, len(req.Profiles)),
	}, nil
}
