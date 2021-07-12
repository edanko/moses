package remnant

import (
	"github.com/edanko/moses/entities"
)

type Reader interface {
	GetOne(id string) (*entities.Remnant, error)
	GetAll() ([]*entities.Remnant, error)
	GetNotUsed(project, dimension, quality string) ([]*entities.Remnant, error)
	GetAllNotUsed() ([]*entities.Remnant, error)
}

type Writer interface {
	Create(e *entities.Remnant) (*entities.Remnant, error)
	Update(e *entities.Remnant) (*entities.Remnant, error)
	Delete(id string) error
	DeleteAll() error
}

type RemnantRepository interface {
	Reader
	Writer
}

type UseCase interface {
	GetOne(id string) (*entities.Remnant, error)
	GetAll() ([]*entities.Remnant, error)
	GetNotUsed(project, dimension, quality string) ([]*entities.Remnant, error)
	GetAllNotUsed() ([]*entities.Remnant, error)
	Create(e *entities.Remnant) (*entities.Remnant, error)
	Update(e *entities.Remnant) (*entities.Remnant, error)
	Delete(id string) error
	DeleteAll() error
}
