package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/edanko/moses/internal/config"
	"github.com/edanko/moses/internal/in"
	"github.com/edanko/moses/internal/nest"
	"github.com/edanko/moses/internal/remnant"
	"github.com/edanko/moses/internal/utils"

	"github.com/spf13/cobra"
)

var Cmd10510Custom = &cobra.Command{
	Use:   "10510_c",
	Short: "process csv with profile parts from my custom report",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("no input!")
			os.Exit(0)
		}
		main10510Custom(args[0])
	},
}

func init() {
	rootCmd.AddCommand(Cmd10510Custom)
}

func main10510Custom(file string) {
	if err := config.Init(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	var nm string
	resultBarList := strings.Builder{}
	resultBarList.Grow(10000)

	resultNestingList := strings.Builder{}
	resultNestingList.Grow(10000)

	profs := in.ProcessAvevaCustomReportCsv(file)

	for _, parts := range profs {
		n := nest.NewNester()

		nm = parts[0].Section
		n.Parts = parts

		n.DynamicNest()
		resultBarList.WriteString(n.BarListString())
		resultNestingList.WriteString(n.NestingListString())

		//filename := path.Join("out", n.Bars[0].Parts[0].Section, n.TxtFileNameString()+".txt")
		filename := path.Join("out", "txt", n.TxtFileNameString()+".txt")

		err := utils.WriteStringToFile(filename, n.TxtOutputString())
		if err != nil {
			panic(err)
		}
		fmt.Println("[+]", n.TxtFileNameString(), "txt successfully created")

		/*for _, b := range in.ProcessTxt(filename) {
			m := mgf.NewMGF(b)

			//fname := path.Join("out", "mgf", b.Project, b.Section, m.NestName)
			fname := path.Join("out", "mgf", m.NestName)

			if err = utils.WriteStringToFile(fname+".mgf", m.MGFOutput()); err != nil {
				log.Fatalln(err)
			}
			if m.StrightCut != "" {
				if err = utils.WriteStringToFile(fname+"_strightcut.txt", m.StrightCut); err != nil {
					log.Fatalln(err)
				}
			}
		}*/
		fmt.Println("[+] mgf's successfully created")

	}

	if err := utils.WriteStringToFile(path.Join("out", nm+"_nst.txt"), resultNestingList.String()); err != nil {
		log.Fatalln(err)
	}
	if err := utils.WriteStringToFile(path.Join("out", nm+"_bar.txt"), resultBarList.String()); err != nil {
		log.Fatalln(err)
	}
	if err := remnant.WriteRemnantsToFile(); err != nil {
		log.Fatalln(err)
	}
}
