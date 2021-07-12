package report

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/edanko/moses/entities"
	"github.com/edanko/moses/service/nest"
	"github.com/edanko/moses/service/remnant"
	"github.com/edanko/moses/service/stock"

	"github.com/edanko/moses/service/spacing"
)

type Service struct {
	nestService    nest.UseCase
	remnantService remnant.UseCase
	spacingService spacing.UseCase
	stockService   stock.UseCase
}

func NewService(n nest.UseCase, r remnant.UseCase, s spacing.UseCase, b stock.UseCase) *Service {
	return &Service{
		nestService:    n,
		remnantService: r,
		spacingService: s,
		stockService:   b,
	}
}

func (s *Service) Bars(nests []*entities.Nest) (string, error) {
	o := strings.Builder{}
	o.Grow(2048)

	byDim := make(map[string]map[string]map[string][]*entities.Nest)

	for _, n := range nests {
		if _, ok := byDim[n.Machine]; !ok {
			byDim[n.Machine] = make(map[string]map[string][]*entities.Nest)
		}

		if _, ok := byDim[n.Machine][n.Bar.Dim]; !ok {
			byDim[n.Machine][n.Bar.Dim] = make(map[string][]*entities.Nest)
		}

		byDim[n.Machine][n.Bar.Dim][n.Bar.Quality] = append(byDim[n.Machine][n.Bar.Dim][n.Bar.Quality], n)
	}

	for m, dd := range byDim {
		fmt.Printf("nests for %s machine\n", m)

		for d, qq := range dd {
			fmt.Printf("dim: %s\n", d)

			for q, nn := range qq {
				fmt.Printf("quality: %s\n", q)

				var stockBars, remnants int
				for _, n := range nn {
					if n.Bar.IsRemnant {
						remnants++
					} else {
						stockBars++
					}

					fmt.Println("Имя программы:", n.Name)
					if n.Bar.IsRemnant {
						fmt.Println("Использовать отход:", n.Bar.Remnant.Marking)
					}
					fmt.Println("Длина:", n.Bar.Length)
					fmt.Println("Исп.:", n.PartsLen())
					if r := n.GetRemnant(); r != nil {
						fmt.Println("Отход:", r.Length)
						fmt.Println("Маркировка отхода:", r.Marking)
					}

					fmt.Println("Лом:", n.Scrap())

					fmt.Printf(" (%.2f %%)\n", n.PartsLen()/(n.Bar.Length-n.Bar.UsedLength)*100)

					for _, p := range n.Profiles {
						fmt.Println("Секция:", p.Section)
						fmt.Println("Позиция:", p.PosNo)
						fmt.Println("Длина:", p.Length)
						fmt.Println("Полная длина:", p.FullLength)
						fmt.Println()
					}

					fmt.Println("---")
				}
				fmt.Println("total nests:", len(nn))
				if stockBars > 0 {
					fmt.Println("stock bars:", stockBars)
				}
				if remnants > 0 {
					fmt.Println("remnants:", remnants)
				}
			}

			fmt.Println()
			fmt.Println()
			fmt.Println()
		}
	}

	for _, n := range nests {
		if n.Bar.IsRemnant {
			o.WriteString(n.Name + ": ДМО " + strconv.FormatFloat(n.Bar.Length, 'f', -1, 64) + "\n")
		} else {
			o.WriteString(n.Name + ": Целый хлыст " + strconv.FormatFloat(n.Bar.Length, 'f', -1, 64) + "\n")
		}
	}

	return o.String(), nil
}

func (s *Service) Nesting(nests []*entities.Nest) (string, error) {
	o := strings.Builder{}
	o.Grow(2048)

	return o.String(), nil
}
