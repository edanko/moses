package nesting

import (
	"context"
	"fmt"

	"github.com/edanko/moses/internal/domain"
)

// Service runs the registered nesting algorithms.
type Service struct {
	registry *Registry
}

// NewService creates a nesting service.
func NewService() *Service {
	return &Service{registry: globalRegistry}
}

// Nest runs the named algorithm. spacing is the default per-profile spacing
// used when the request has no spacing function.
func (s *Service) Nest(
	ctx context.Context,
	algorithmName string,
	profiles []domain.Part,
	stock []domain.StockBar,
	remnants []domain.StockBar,
	spacing float64,
	parameters map[string]any,
) (*Result, error) {
	algorithm, err := s.registry.Get(algorithmName)
	if err != nil {
		return nil, fmt.Errorf("get algorithm %q: %w", algorithmName, err)
	}

	request := &Request{
		Profiles:   profiles,
		StockBar:   stock,
		Remnants:   remnants,
		Spacing:    spacing,
		Parameters: parameters,
	}
	return algorithm.Nest(ctx, request)
}

// ListAlgorithms returns the registered algorithm names.
func (s *Service) ListAlgorithms() []string {
	return s.registry.List()
}

// GetAlgorithmInfo returns the algorithm name.
func (s *Service) GetAlgorithmInfo(name string) (string, error) {
	algorithm, err := s.registry.Get(name)
	if err != nil {
		return "", err
	}
	return algorithm.Name(), nil
}
