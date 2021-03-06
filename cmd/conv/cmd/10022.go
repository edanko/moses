package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/edanko/moses/internal/config"
	"github.com/edanko/moses/internal/in"
	"github.com/edanko/moses/internal/models"
	"github.com/edanko/moses/internal/nest"
	"github.com/edanko/moses/internal/remnant"
	"github.com/edanko/moses/internal/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Cmd10022 = &cobra.Command{
	Use:   "10022",
	Short: "process dxf profile report files from some cadmatic version",
	Run: func(cmd *cobra.Command, args []string) {
		main10022()
	},
}

func init() {
	rootCmd.AddCommand(Cmd10022)
}

func main10022() {
	// load config
	if err := config.Init(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	// load input file list
	files, err := filepath.Glob("in" + string(os.PathSeparator) + "*" + string(os.PathSeparator) + "*.dxf")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if len(files) == 0 {
		fmt.Println("No input files!")
		os.Exit(0)
	}

	parts := models.NewPartsStorage()

	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)

		go func() {
			defer wg.Done()

			f, err := os.Open(file)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			fmt.Println(file)

			part := in.ProcessMssDxf(f)

			parts.Add(part)
		}()
		wg.Wait()
	}

	nm := parts.Project() + "-"
	sc := parts.Sections()
	if len(sc) == 1 {
		nm += sc[0]
	} else {
		viper.Set("withsection", true)
		nm += "xxx"
	}

	resultBarList := strings.Builder{}
	resultBarList.Grow(10000)

	resultNestingList := strings.Builder{}
	resultNestingList.Grow(10000)

	for _, t := range parts.Dims() {
		for _, m := range parts[t].Qualitys() {

			n := nest.New()

			partsIterator := parts[t][m].Iterator()
			for partsIterator.Next() {

				p := partsIterator.Value().(*models.Part)

				for i := 0; i < p.Quantity; i++ {
					n.AddPart(p)
				}
			}
			n.Nest()
			resultBarList.WriteString(n.BarListString())
			resultNestingList.WriteString(n.NestingListString())

			filename := path.Join("out", "txt", n.TxtFileNameString()+".txt")
			err = utils.WriteStringToFile(filename, n.TxtOutputString())
			if err != nil {
				panic(err)
			}
			fmt.Println("[+]", n.TxtFileNameString(), "txt successfully created")
		}
	}

	err = utils.WriteStringToFile(path.Join("out", nm+"_nst.txt"), resultNestingList.String())
	if err != nil {
		panic(err)
	}
	err = utils.WriteStringToFile(path.Join("out", nm+"_bar.txt"), resultBarList.String())
	if err != nil {
		panic(err)
	}
	if err := remnant.WriteRemnantsToFile(); err != nil {
		log.Fatalln(err)
	}
}
