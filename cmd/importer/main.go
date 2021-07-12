package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/edanko/gen2dxf/pkg/wcog"
	repo "github.com/edanko/moses/repository/mongo"
	"github.com/edanko/moses/service/importer"
	"github.com/edanko/moses/service/nest"
	"github.com/edanko/moses/service/profile"
	"github.com/edanko/moses/service/remnant"
	"github.com/edanko/moses/service/spacing"
	"github.com/joho/godotenv"
)

// for now just aveva gen test
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

	// txt

	/* 	var files []string
	   	err = filepath.Walk(".", func(path string, info os.FileInfo, e error) error {
	   		if e != nil {
	   			return e
	   		}

	   		if info.Mode().IsRegular() {
	   			switch filepath.Ext(path) {
	   			case ".txt":
	   				files = append(files, path)
	   			}
	   		}

	   		return nil
	   	})

	   	if err != nil {
	   		fmt.Println(err.Error())
	   		os.Exit(1)
	   	}
	   	if len(files) == 0 {
	   		fmt.Println("No input files!")
	   		os.Exit(0)
	   	}

	   	svc := importer.NewService(profileService, nestService)

	   	for _, f := range files {
	   		fmt.Println("[i] reading", filepath.Base(f))
	   		err = svc.ImportTxt(f)

	   		if err != nil {
	   			log.Fatalln(err)
	   		}
	   	} */

	// gen, wcog
	var files []string
	var wcogs []string
	err = filepath.Walk(".", func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}

		if info.Mode().IsRegular() {
			switch filepath.Ext(path) {
			case ".csv":
				wcogs = append(wcogs, path)
			case ".gen":
				files = append(files, path)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if len(files) == 0 {
		fmt.Println("No input files!")
		os.Exit(0)
	}

	wcog := wcog.ReadWCOGs(wcogs)

	svc := importer.NewService(profileService, nestService)

	for _, f := range files {
		fmt.Println("[i] reading", filepath.Base(f))
		err = svc.ImportGen(f, wcog)

		if err != nil {
			log.Fatalln(err)
		}
	}
	// aveva csv
	/* var files []string
	err = filepath.Walk(".", func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}

		if info.Mode().IsRegular() {
			switch filepath.Ext(path) {
			case ".csv":
				files = append(files, path)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if len(files) == 0 {
		fmt.Println("No input files!")
		os.Exit(0)
	}

	svc := importer.NewService(profileService, nestService)

	for _, f := range files {
		fmt.Println("[i] reading", filepath.Base(f))
		err = svc.ImportCsv(f)

		if err != nil {
			log.Fatalln(err)
		}
	} */

}
