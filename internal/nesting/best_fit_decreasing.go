package nesting

import (
	"context"

	"github.com/edanko/moses/internal/domain"
	"github.com/edanko/moses/pkg/binpack"
)

// BestFitDecreasingAlgorithm implements the best-fit decreasing bin packing algorithm
type BestFitDecreasingAlgorithm struct{}

// NewBestFitDecreasingAlgorithm creates a new best-fit decreasing algorithm
func NewBestFitDecreasingAlgorithm() *BestFitDecreasingAlgorithm {
	return &BestFitDecreasingAlgorithm{}
}

// Name returns the algorithm name
func (b *BestFitDecreasingAlgorithm) Name() string {
	return "best-fit-decreasing"
}

// Nest performs the nesting operation using best-fit decreasing algorithm
func (b *BestFitDecreasingAlgorithm) Nest(ctx context.Context, req *Request) (*Result, error) {
	items := toBinpackItems(req.Profiles)
	bins := toBinpackBins(req.StockBar)

	if len(req.Remnants) > 0 {
		remnantBins := toBinpackBins(req.Remnants)
		bins = append(remnantBins, bins...)
	}

	algorithm := binpack.NewBestFitDecreasingAlgorithm[domain.StockBar, domain.Part]()
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
