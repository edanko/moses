package stock

import (
	"github.com/edanko/moses/entities"
)

type Reader interface {
	GetOne(dim, quality string) (*entities.Stock, error)
	GetAll() ([]*entities.Stock, error)
}

type Writer interface {
	Create(e *entities.Stock) (*entities.Stock, error)
	Update(e *entities.Stock) (*entities.Stock, error)
	Delete(id string) error
	DeleteAll() error
}

type StockRepository interface {
	Reader
	Writer
}

type UseCase interface {
	GetOne(dim, quality string) (*entities.Stock, error)
	GetAll() ([]*entities.Stock, error)
	Create(e *entities.Stock) (*entities.Stock, error)
	Update(e *entities.Stock) (*entities.Stock, error)
	Delete(id string) error
	DeleteAll() error
}
