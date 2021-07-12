package spacing

import (
	"github.com/edanko/moses/entities"
)

type Reader interface {
	GetAll(machine string) ([]*entities.Spacing, error)
	GetOne(machine, dim string, e *entities.End) (*entities.Spacing, error)
}

type Writer interface {
	Create(e *entities.Spacing) (*entities.Spacing, error)
	Update(e *entities.Spacing) (*entities.Spacing, error)
	Delete(id string) error
	DeleteAll() error
}

type SpacingRepository interface {
	Reader
	Writer
}

type UseCase interface {
	GetAll(machine string) ([]*entities.Spacing, error)
	GetOne(machine, dim string, e *entities.End) (*entities.Spacing, error)
	Create(e *entities.Spacing) (*entities.Spacing, error)
	Update(e *entities.Spacing) (*entities.Spacing, error)
	Delete(id string) error
	DeleteAll() error
}
