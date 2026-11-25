package nesting

import (
	"context"

	"github.com/edanko/moses/internal/domain"
	"github.com/edanko/moses/pkg/binpack"
)

// KnapsackAlgorithm implements the knapsack-based bin packing algorithm
type KnapsackAlgorithm struct{}

// NewKnapsackAlgorithm creates a new knapsack algorithm
func NewKnapsackAlgorithm() *KnapsackAlgorithm {
	return &KnapsackAlgorithm{}
}

// Name returns the algorithm name
func (k *KnapsackAlgorithm) Name() string {
	return "knapsack"
}

// Nest performs the nesting operation using knapsack algorithm
func (k *KnapsackAlgorithm) Nest(ctx context.Context, req *Request) (*Result, error) {
	items := toBinpackItems(req.Profiles)
	bins := toBinpackBins(req.StockBar)

	if len(req.Remnants) > 0 {
		remnantBins := toBinpackBins(req.Remnants)
		bins = append(remnantBins, bins...)
	}

	algorithm := binpack.NewKnapsackAlgorithm[domain.StockBar, domain.Part]()
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
