package nest

import (
	"github.com/edanko/moses/entities"
)

type Reader interface {
	GetOne(id string) (*entities.Nest, error)
	GetAll() ([]*entities.Nest, error)
	Get(project, dimension, quality string) ([]*entities.Nest, error)
}

type Writer interface {
	Create(e *entities.Nest) (*entities.Nest, error)
	Update(e *entities.Nest) (*entities.Nest, error)
	Delete(id string) error
	DeleteAll() error
}

type NestRepository interface {
	Reader
	Writer
}

type UseCase interface {
	GetOne(id string) (*entities.Nest, error)
	GetAll() ([]*entities.Nest, error)
	Get(project, dimension, quality string) ([]*entities.Nest, error)
	Create(e *entities.Nest) (*entities.Nest, error)
	Update(e *entities.Nest) (*entities.Nest, error)
	Delete(id string) error
	DeleteAll() error
}
