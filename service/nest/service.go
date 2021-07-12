package nest

import (
	"github.com/edanko/moses/entities"
	"github.com/edanko/moses/service/profile"
	"github.com/edanko/moses/service/remnant"
	"github.com/edanko/moses/service/spacing"
)

type Service struct {
	repo           NestRepository
	remnantService remnant.UseCase
	profileService profile.UseCase
	spacingService spacing.UseCase
}

func NewService(n NestRepository, r remnant.UseCase, p profile.UseCase, s spacing.UseCase) *Service {
	return &Service{
		repo:           n,
		remnantService: r,
		profileService: p,
		spacingService: s,
	}
}

func (s *Service) Create(e *entities.Nest) (*entities.Nest, error) {
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

func (s *Service) GetOne(id string) (*entities.Nest, error) {
	n, err := s.repo.GetOne(id)
	if n == nil {
		return nil, entities.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (s *Service) Get(project, dimension, quality string) ([]*entities.Nest, error) {
	nests, err := s.repo.Get(project, dimension, quality)
	if err != nil {
		return nil, err
	}
	if len(nests) == 0 {
		return nil, entities.ErrNotFound
	}
	return nests, nil
}

func (s *Service) GetAll() ([]*entities.Nest, error) {
	nests, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	if len(nests) == 0 {
		return nil, entities.ErrNotFound
	}

	for _, n := range nests {
		for _, p := range n.ProfilesIds {
			profile, err := s.profileService.GetOne(p.Hex())
			if err != nil {
				return nil, err
			}

			l, err := s.spacingService.GetOne(n.Machine, n.Bar.Dim, profile.L)
			if err != nil {
				return nil, err
			}
			r, err := s.spacingService.GetOne(n.Machine, n.Bar.Dim, profile.R)
			if err != nil {
				return nil, err
			}

			profile.FullLength = l.Length + profile.Length + r.Length

			n.Profiles = append(n.Profiles, profile)
		}
		if n.Bar.IsRemnant {
			remnant, err := s.remnantService.GetOne(n.Bar.RemnantID.Hex())
			if err != nil {
				return nil, err
			}
			n.Bar.Remnant = remnant
		}
	}

	return nests, nil
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *Service) DeleteAll() error {
	return s.repo.DeleteAll()
}

func (s *Service) Update(e *entities.Nest) (*entities.Nest, error) {
	err := e.Validate()
	if err != nil {
		return nil, err
	}
	return s.repo.Update(e)
}
