package importer

import (
	"github.com/edanko/gen2dxf/pkg/wcog"
	"github.com/edanko/moses/internal/formats"
	"github.com/edanko/moses/service/nest"
	"github.com/edanko/moses/service/profile"
)

type Service struct {
	profileService profile.UseCase
	nestService    nest.UseCase
}

func NewService(p profile.UseCase, n nest.UseCase) *Service {
	return &Service{
		profileService: p,
		nestService:    n,
	}
}

func (s *Service) ImportGen(f string, wcog *wcog.WCOG) error {
	profiles, err := formats.ProcessGen(f, wcog)
	if err != nil {
		return err
	}

	for _, profile := range profiles {
		_, err := s.profileService.Create(profile)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) ImportCsv(f string) error {
	profiles, err := formats.ProcessCsv(f)
	if err != nil {
		return err
	}

	for _, profile := range profiles {
		_, err := s.profileService.Create(profile)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) ImportTxt(f string) error {
	nests, err := formats.ProcessTxt(f)
	if err != nil {
		return err
	}

	for _, nest := range nests {
		for _, profile := range nest.Profiles {
			created, err := s.profileService.Create(profile)
			if err != nil {
				return err
			}
			profile.ID = created.ID
		}

		_, err := s.nestService.Create(nest)
		if err != nil {
			return err
		}

	}

	return nil
}
