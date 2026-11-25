package nesting

import (
	"context"

	"github.com/edanko/moses/internal/domain"
	"github.com/edanko/moses/pkg/binpack"
)

// FirstFitAlgorithm implements the first-fit bin packing algorithm
type FirstFitAlgorithm struct{}

// NewFirstFitAlgorithm creates a new first-fit algorithm
func NewFirstFitAlgorithm() *FirstFitAlgorithm {
	return &FirstFitAlgorithm{}
}

// Name returns the algorithm name
func (f *FirstFitAlgorithm) Name() string {
	return "first-fit"
}

// Nest performs the nesting operation using first-fit algorithm
func (f *FirstFitAlgorithm) Nest(ctx context.Context, req *Request) (*Result, error) {
	items := toBinpackItems(req.Profiles)
	bins := toBinpackBins(req.StockBar)

	if len(req.Remnants) > 0 {
		remnantBins := toBinpackBins(req.Remnants)
		bins = append(remnantBins, bins...)
	}

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
