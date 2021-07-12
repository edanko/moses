package nester

import (
	"github.com/edanko/moses/entities"
)

type UseCase interface {
	Nest(parts []*entities.Profile) ([]*entities.Nest, error)
	Renest(nests []*entities.Nest) ([]*entities.Nest, error)
}
