package profile

import (
	"github.com/edanko/moses/entities"
)

type Service struct {
	repo ProfileRepository
}

func NewService(r ProfileRepository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) Create(e *entities.Profile) (*entities.Profile, error) {
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

func (s *Service) GetOne(id string) (*entities.Profile, error) {
	r, err := s.repo.GetOne(id)
	if r == nil {
		return nil, entities.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *Service) Get(project, dimension, quality string) ([]*entities.Profile, error) {
	profiles, err := s.repo.Get(project, dimension, quality)
	if err != nil {
		return nil, err
	}
	if len(profiles) == 0 {
		return nil, entities.ErrNotFound
	}
	return profiles, nil
}

func (s *Service) GetAll() ([]*entities.Profile, error) {
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

func (s *Service) Update(e *entities.Profile) (*entities.Profile, error) {
	err := e.Validate()
	if err != nil {
		return nil, err
	}
	return s.repo.Update(e)
}
