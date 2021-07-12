package profile

import (
	"github.com/edanko/moses/entities"
)

type Reader interface {
	GetOne(id string) (*entities.Profile, error)
	GetAll() ([]*entities.Profile, error)
	Get(project, dimension, quality string) ([]*entities.Profile, error)
}

type Writer interface {
	Create(e *entities.Profile) (*entities.Profile, error)
	Update(e *entities.Profile) (*entities.Profile, error)
	Delete(id string) error
	DeleteAll() error
}

type ProfileRepository interface {
	Reader
	Writer
}

type UseCase interface {
	GetOne(id string) (*entities.Profile, error)
	GetAll() ([]*entities.Profile, error)
	Get(project, dimension, quality string) ([]*entities.Profile, error) // get free?
	Create(e *entities.Profile) (*entities.Profile, error)
	Update(e *entities.Profile) (*entities.Profile, error)
	Delete(id string) error
	DeleteAll() error
}
