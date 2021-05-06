package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/edanko/moses/internal/config"
	"github.com/edanko/moses/internal/in"
	"github.com/edanko/moses/internal/out/mgf"
	"github.com/edanko/moses/internal/utils"
	"github.com/spf13/viper"
)

func main() {

	// license check
	//if ok, err := license.IsLicensed(); !ok {
	//	log.Fatalln(err.Error())
	//}

	// load config
	if err := config.Init(); err != nil {
		log.Fatalln(err.Error())
	}

	if os.Args[1] == "--dont-rename" {
		viper.Set("dontrename", true)
	}

	// load input file list
	var files []string
	err := filepath.Walk("in", func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}

		if info.Mode().IsRegular() {
			if filepath.Ext(path) == ".txt" && !strings.HasPrefix(filepath.Base(path), "_") {
				files = append(files, path)
			}
		}

		return nil
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(files) == 0 {
		fmt.Println("No input files!")
		os.Exit(0)
	}

	// TODO: add nester

	for _, file := range files {
		fmt.Printf(" * processing %s...\n", strings.Split(file, string(os.PathSeparator))[1])

		for _, b := range in.ProcessTxt(file) {

			fname := path.Join("out", b.Project, b.Section, b.NestName)

			if err = utils.WriteStringToFile(fname+".mgf", mgf.Mgf(b)); err != nil {
				log.Fatalln(err)
			}

			st := mgf.StrightCut(b)
			if st != "" {
				if err = utils.WriteStringToFile(fname+"_strightcut.txt", st); err != nil {
					log.Fatalln(err)
				}
			}
		}

		if !viper.GetBool("dontrename") {
			newInputFileName := filepath.Dir(file) + string(os.PathSeparator) + "_" + filepath.Base(file)
			err := os.Rename(file, newInputFileName)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}
