package nesting

import (
	"errors"
	"sync"
)

var (
	ErrAlgorithmNotFound = errors.New("nesting algorithm not found")
	ErrAlgorithmExists   = errors.New("nesting algorithm already exists")
)

// Registry manages available nesting algorithms
type Registry struct {
	algorithms map[string]Algorithm
	mu         sync.RWMutex
}

// NewRegistry creates a new algorithm registry
func NewRegistry() *Registry {
	return &Registry{
		algorithms: make(map[string]Algorithm),
	}
}

// Register adds a new algorithm to the registry
func (r *Registry) Register(algorithm Algorithm) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := algorithm.Name()
	if _, exists := r.algorithms[name]; exists {
		return ErrAlgorithmExists
	}

	r.algorithms[name] = algorithm
	return nil
}

// Get retrieves an algorithm by name
func (r *Registry) Get(name string) (Algorithm, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	algorithm, exists := r.algorithms[name]
	if !exists {
		return nil, ErrAlgorithmNotFound
	}

	return algorithm, nil
}

// List returns all registered algorithm names
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.algorithms))
	for name := range r.algorithms {
		names = append(names, name)
	}

	return names
}

// global registry instance
var globalRegistry = NewRegistry()

// RegisterAlgorithm registers an algorithm globally
func RegisterAlgorithm(algorithm Algorithm) error {
	return globalRegistry.Register(algorithm)
}

// GetAlgorithm gets an algorithm globally
func GetAlgorithm(name string) (Algorithm, error) {
	return globalRegistry.Get(name)
}

// ListAlgorithms lists all registered algorithms globally
func ListAlgorithms() []string {
	return globalRegistry.List()
}
