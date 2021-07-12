package remnant

import (
	"github.com/edanko/moses/entities"
)

type Service struct {
	repo RemnantRepository
}

func NewService(r RemnantRepository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) Create(e *entities.Remnant) (*entities.Remnant, error) {
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

func (s *Service) GetOne(id string) (*entities.Remnant, error) {
	r, err := s.repo.GetOne(id)
	if r == nil {
		return nil, entities.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *Service) GetNotUsed(project, dimension, quality string) ([]*entities.Remnant, error) {
	remnants, err := s.repo.GetNotUsed(project, dimension, quality)
	if err != nil {
		return nil, err
	}
	//if len(remnants) == 0 {
	//	return nil, entities.ErrNotFound
	//}
	return remnants, nil
}

func (s *Service) GetAllNotUsed() ([]*entities.Remnant, error) {
	remnants, err := s.repo.GetAllNotUsed()
	if err != nil {
		return nil, err
	}
	if len(remnants) == 0 {
		return nil, entities.ErrNotFound
	}
	return remnants, nil
}

func (s *Service) GetAll() ([]*entities.Remnant, error) {
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

func (s *Service) Update(e *entities.Remnant) (*entities.Remnant, error) {
	err := e.Validate()
	if err != nil {
		return nil, err
	}
	return s.repo.Update(e)
}
