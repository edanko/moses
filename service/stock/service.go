package stock

import (
	"fmt"

	"github.com/edanko/moses/entities"
)

type Service struct {
	repo StockRepository
}

func NewService(r StockRepository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) Create(e *entities.Stock) (*entities.Stock, error) {
	err := e.Validate()
	if err != nil {
		return nil, err
	}
	r, err := s.repo.Create(e)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *Service) GetOne(dim, quality string) (*entities.Stock, error) {
	r, err := s.repo.GetOne(dim, quality)
	if r == nil {
		return nil, fmt.Errorf("stock bar length for dim \"%s\" and quality \"%s\"", dim, quality)
	}
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *Service) GetAll() ([]*entities.Stock, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, entities.ErrNotFound
	}
	return users, nil
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *Service) DeleteAll() error {
	return s.repo.DeleteAll()
}

func (s *Service) Update(e *entities.Stock) (*entities.Stock, error) {
	err := e.Validate()
	if err != nil {
		return nil, err
	}
	return s.repo.Update(e)
}
