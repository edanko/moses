package main

import (
	"fmt"
	"log"
	"os"

	repo "github.com/edanko/moses/repository/mongo"
	"github.com/edanko/moses/service/nest"
	"github.com/edanko/moses/service/profile"
	"github.com/edanko/moses/service/remnant"
	"github.com/edanko/moses/service/report"
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

	profileCollection := db.Collection("profiles")
	profileRepo := repo.NewProfileRepo(profileCollection)
	profileService := profile.NewService(profileRepo)

	remnantCollection := db.Collection("remnants")
	remnantRepo := repo.NewRemnantRepo(remnantCollection)
	remnantService := remnant.NewService(remnantRepo)

	/* 	profileCollection := db.Collection("profiles")
	   	profileRepo := repo.NewProfileRepo(profileCollection)
	   	profileService := profile.NewService(profileRepo) */

	spacingCollection := db.Collection("spacing")
	spacingRepo := repo.NewSpacingRepo(spacingCollection)
	spacingService := spacing.NewService(spacingRepo)

	nestCollection := db.Collection("nests")
	nestRepo := repo.NewNestRepo(nestCollection)
	nestService := nest.NewService(nestRepo, remnantService, profileService, spacingService)

	stockCollection := db.Collection("stock")
	stockRepo := repo.NewStockRepo(stockCollection)
	stockService := stock.NewService(stockRepo)

	rep := report.NewService(nestService, remnantService, spacingService, stockService)

	nests, err := nestService.GetAll()
	if err != nil {
		log.Fatalln("no nests found")
	}

	str, err := rep.Bars(nests)
	if err != nil {
		panic(err)
	}

	//fmt.Println("bar-list")
	fmt.Println(str)

	/* 	str, err = rep.Nesting(nests)
	   	if err != nil {
	   		panic(err)
	   	}
	   	fmt.Println("nesting-list")
	   	fmt.Println(str) */

}
