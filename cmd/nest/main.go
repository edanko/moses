package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/edanko/moses/entities"
	repo "github.com/edanko/moses/repository/mongo"
	"github.com/edanko/moses/service/nest"
	"github.com/edanko/moses/service/nester"
	"github.com/edanko/moses/service/profile"
	"github.com/edanko/moses/service/remnant"
	"github.com/edanko/moses/service/spacing"
	"github.com/edanko/moses/service/stock"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	db, err := repo.NewMongoDB(os.Getenv("MONGODB_URI"), os.Getenv("MONGODB_DATABASE"))
	if err != nil {
		log.Fatalln(err)
	}

	remnantCollection := db.Collection("remnants")
	remnantRepo := repo.NewRemnantRepo(remnantCollection)
	remnantService := remnant.NewService(remnantRepo)

	profileCollection := db.Collection("profiles")
	profileRepo := repo.NewProfileRepo(profileCollection)
	profileService := profile.NewService(profileRepo)

	spacingCollection := db.Collection("spacing")
	spacingRepo := repo.NewSpacingRepo(spacingCollection)
	spacingService := spacing.NewService(spacingRepo)

	nestCollection := db.Collection("nests")
	nestRepo := repo.NewNestRepo(nestCollection)
	nestService := nest.NewService(nestRepo, remnantService, profileService, spacingService)

	stockCollection := db.Collection("stock")
	stockRepo := repo.NewStockRepo(stockCollection)
	stockService := stock.NewService(stockRepo)

	nester := nester.NewService(nestService, remnantService, spacingService, stockService)

	profiles, _ := profileService.GetAll()

	/* 	sp := &entities.Spacing{}
	   	sp.Machine = "img"
	   	sp.Dim = "RP180*11.0"
	   	sp.Name = "23"
	   	sp.HasBevel = false
	   	sp.HasScallop = false
	   	sp.Length = 10

	   	_, err = spacingService.Create(sp)
	   	if err != nil {
	   		log.Fatalln(err)
	   	} */

	/* 	s := &entities.Stock{}
	   	s.Dim = "RP200*12"
	   	s.Quality = "A40"
	   	s.Length = 12000

	   	_, err = stockService.Create(s)
	   	if err != nil {
	   		log.Fatalln(err)
	   	} */

	//byQuality := make(map[string]*entities.Profile)
	byDim := make(map[string]map[string][]*entities.Profile)

	for _, p := range profiles {

		/* 		if _, exist := byDim[p.Dim]; !exist {
			byDim = make(map[string]map[string][]*entities.Profile)
		} */
		if _, exist := byDim[p.Dim][p.Quality]; !exist {
			byDim[p.Dim] = make(map[string][]*entities.Profile)
		}

		byDim[p.Dim][p.Quality] = append(byDim[p.Dim][p.Quality], p)
	}

	machine := "img"

	wg := &sync.WaitGroup{}

	for _, dims := range byDim {
		for _, profiles := range dims {
			wg.Add(1)
			go func(ps []*entities.Profile) {
				defer wg.Done()

				nests, err := nester.Nest(machine, ps)
				if err != nil {
					log.Fatalln(err)
				}

				for i, nest := range nests {
					for _, p := range nest.Profiles {
						nest.ProfilesIds = append(nest.ProfilesIds, p.ID)
					}

					nest.Machine = machine
					nest.Project = nest.Profiles[0].Project
					nest.Launch = "launch placeholder"
					nest.Name = "nest_" + strconv.Itoa(i)

					// save nest
					_, err = nestService.Create(nest)
					if err != nil {
						panic(err)
					}

					// use remnant
					if nest.Bar.IsRemnant {
						rem, err := remnantService.GetOne(nest.Bar.RemnantID.Hex())
						if err != nil {
							panic(err)
						}

						rem.Used = true
						rem.UsedIn = nest.Name

						_, err = remnantService.Update(rem)
						if err != nil {
							panic(err)
						}
					}

					// save remnant
					if rem := nest.GetRemnant(); rem != nil {
						_, err = remnantService.Create(rem)
						if err != nil {
							panic(err)
						}
					}
				}

			}(profiles)
		}
	}

	wg.Wait()
}
