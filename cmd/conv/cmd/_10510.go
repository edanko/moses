package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/edanko/moses/internal/config"
	"github.com/edanko/moses/internal/in"
	"github.com/edanko/moses/internal/nest"
	"github.com/edanko/moses/internal/remnant"
	"github.com/edanko/moses/internal/utils"

	"github.com/spf13/cobra"
)

var Cmd10510 = &cobra.Command{
	Use:   "10510",
	Short: "process csv&lst from aveva profile sketch&list",
	Run: func(cmd *cobra.Command, args []string) {
		main10510()
	},
}

func init() {
	rootCmd.AddCommand(Cmd10510)
}

func main10510() {
	if err := config.Init(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	files, err := filepath.Glob("in" + string(os.PathSeparator) + "*.csv")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if len(files) == 0 {
		fmt.Println("No input files!")
		os.Exit(0)
	}

	var nm string
	resultBarList := strings.Builder{}
	resultBarList.Grow(10000)

	resultNestingList := strings.Builder{}
	resultNestingList.Grow(10000)

	for i, f := range files {

		fmt.Println("[i] reading", f)
		n := nest.NewNester()
		parts := in.ProcessAvevaCsv(f)

		if i == 0 {
			nm = parts[0].Section
		}
		n.Parts = parts
		n.DynamicNest()
		resultBarList.WriteString(n.BarListString())
		resultNestingList.WriteString(n.NestingListString())

		//filename := path.Join("out", n.Bars[0].Parts[0].Section, n.TxtFileNameString()+".txt")
		filename := path.Join("out", time.Now().Format("06.01.02")+" "+n.Bars[0].Section(), n.TxtFileNameString()+".txt")

		err := utils.WriteStringToFile(filename, n.TxtOutputString())
		if err != nil {
			panic(err)
		}
		fmt.Println("[+]", n.TxtFileNameString(), "txt successfully created")

		/*for _, b := range in.ProcessTxt(filename) {
			m := mgf.NewMGF(b)

			//fname := path.Join("out", "mgf", b.Project, b.Section, m.NestName)
			fname := path.Join("out", "mgf", m.NestName)

			if err = utils.WriteStringToFile(fname+".mgf", m.GetMGFOutput()); err != nil {
				log.Fatalln(err)
			}
			if m.StrightCut != "" {
				if err = utils.WriteStringToFile(fname+"_strightcut.txt", m.StrightCut); err != nil {
					log.Fatalln(err)
				}
			}
		}
		fmt.Println("[+] mgf's successfully created")*/
	}

	err = utils.WriteStringToFile(path.Join("out", time.Now().Format("06.01.02")+" "+nm+"_nst.txt"), resultNestingList.String())
	if err != nil {
		panic(err)
	}
	err = utils.WriteStringToFile(path.Join("out", time.Now().Format("06.01.02")+" "+nm+"_bar.txt"), resultBarList.String())
	if err != nil {
		panic(err)
	}
	if err := remnant.WriteRemnantsToFile(); err != nil {
		log.Fatalln(err)
	}
}
