package report

import (
	"github.com/edanko/moses/entities"
)

type UseCase interface {
	Bars(nests []*entities.Profile) (string, error)
	Nesting(nests []*entities.Profile) (string, error)
}
