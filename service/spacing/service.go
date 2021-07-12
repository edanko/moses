package spacing

import (
	"fmt"

	"github.com/edanko/moses/entities"
)

type Service struct {
	repo SpacingRepository
}

func NewService(r SpacingRepository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) Create(e *entities.Spacing) (*entities.Spacing, error) {
	err := e.Validate()
	if err != nil {
		return nil, err
	}
	n, err := s.repo.Create(e)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (s *Service) GetOne(machine, dim string, e *entities.End) (*entities.Spacing, error) {
	sp, err := s.repo.GetOne(machine, dim, e)
	if err != nil && err == entities.ErrNotFound {
		bevel := e.WebBevel != nil || e.FlangeBevel != nil
		scallop := e.Scallop != nil
		return nil, fmt.Errorf("spacing for machine \"%s\", dim \"%s\", name \"%s\", bevel \"%t\", scallop \"%t\" not found", machine, dim, e.Name, bevel, scallop)
	} else if err != nil {
		return nil, err
	}
	return sp, nil
}

func (s *Service) GetAll(machine string) ([]*entities.Spacing, error) {
	spacings, err := s.repo.GetAll(machine)
	if err != nil {
		return nil, err
	}
	if len(spacings) == 0 {
		return nil, entities.ErrNotFound
	}
	return spacings, nil
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *Service) DeleteAll() error {
	return s.repo.DeleteAll()
}

func (s *Service) Update(e *entities.Spacing) (*entities.Spacing, error) {
	err := e.Validate()
	if err != nil {
		return nil, err
	}
	return s.repo.Update(e)
}
