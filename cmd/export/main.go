package main

import (
	"fmt"
	"log"
	"os"

	repo "github.com/edanko/moses/repository/mongo"
	"github.com/edanko/moses/service/nest"
	"github.com/edanko/moses/service/profile"
	"github.com/edanko/moses/service/remnant"
	"github.com/edanko/moses/service/spacing"
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

	spacingCollection := db.Collection("spacing")
	spacingRepo := repo.NewSpacingRepo(spacingCollection)
	spacingService := spacing.NewService(spacingRepo)

	nestCollection := db.Collection("nests")
	nestRepo := repo.NewNestRepo(nestCollection)
	nestService := nest.NewService(nestRepo, remnantService, profileService, spacingService)

	nests, err := nestService.GetAll()
	if err != nil {
		panic(err)
	}

	for _, nest := range nests {
		fmt.Println(nest.ID)
	}
}
