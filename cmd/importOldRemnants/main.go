package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/edanko/moses/entities"
	repo "github.com/edanko/moses/repository/mongo"
	"github.com/edanko/moses/service/remnant"
	"github.com/joho/godotenv"
	"github.com/jszwec/csvutil"
)

type csvRemnant struct {
	Id          int     `csv:"id"`
	Prj         string  `csv:"prj"`
	Sec         string  `csv:"sec"`
	Dim         string  `csv:"dim"`
	Quality     string  `csv:"quality"`
	RemLen      float64 `csv:"rem_len"`
	MarkingText string  `csv:"marking_text"`
	UsedIn      string  `csv:"used_in"`
}

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

	remnantCollection := db.Collection("remnants")
	remnantRepo := repo.NewRemnantRepo(remnantCollection)
	remnantService := remnant.NewService(remnantRepo)

	// aveva csv
	var files []string
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

	var total, notUsed int

	for _, f := range files {
		fmt.Println("[i] reading", filepath.Base(f))

		f, err := os.Open(f)
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()

		csvReader := csv.NewReader(bufio.NewReader(f))
		csvReader.Comma = ';'

		dec, err := csvutil.NewDecoder(csvReader)
		if err != nil {
			log.Fatal(err)
		}

		for {
			var r csvRemnant
			if err := dec.Decode(&r); err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			rem := &entities.Remnant{}
			rem.Project = r.Prj
			rem.From = r.Sec
			rem.Dim = strings.Replace(r.Dim, "X", "*", 1)

			if !strings.Contains(rem.Dim, ".") {
				rem.Dim += ".0"
			}

			rem.Quality = r.Quality
			rem.Length = r.RemLen
			rem.Marking = r.MarkingText
			rem.UsedIn = r.UsedIn
			rem.Used = r.UsedIn != ""

			total++
			if !rem.Used {
				notUsed++
			}

			_, err = remnantService.Create(rem)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	fmt.Println("imported:", total, ", not used:", notUsed)

}
