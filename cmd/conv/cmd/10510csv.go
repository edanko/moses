package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/edanko/moses/internal/config"
	"github.com/edanko/moses/internal/in"
	"github.com/edanko/moses/internal/models"
	"github.com/edanko/moses/internal/nest"
	"github.com/edanko/moses/internal/remnant"
	"github.com/edanko/moses/internal/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	parts := models.NewPartsStorage()

	for _, file := range files {
		fmt.Println(file)
		part := in.ProcessAvevaCsv(file)
		for _, p := range part {
			parts.Add(p)
		}
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

	var allBars []*models.Bar

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

			allBars = append(allBars, n.Bars...)

			filename := path.Join("out", "txt", n.TxtFileNameString()+".txt")
			err = utils.WriteStringToFile(filename, n.TxtOutputString())
			if err != nil {
				panic(err)
			}
			fmt.Println("[+]", n.TxtFileNameString()+".txt", "successfully created")
		}
	}

	barlist := "Ведомость расхода материала"
	partlist := "Партлист"

	f := excelize.NewFile()

	f.SetSheetName("Sheet1", barlist)

	_ = f.SetCellValue(barlist, "A1", "Карта раскроя")       // 1
	_ = f.SetCellValue(barlist, "B1", "Тип оборудования")    // 2
	_ = f.SetCellValue(barlist, "C1", "Исп. ДМО")            // 3
	_ = f.SetCellValue(barlist, "D1", "Типоразмер")          // 4
	_ = f.SetCellValue(barlist, "E1", "Марка материала")     // 5
	_ = f.SetCellValue(barlist, "F1", "Длина, мм")           // 6
	_ = f.SetCellValue(barlist, "G1", "Масса заготовки, кг") // 7
	_ = f.SetCellValue(barlist, "H1", "Коэф. раскроя, %")    // 8
	_ = f.SetCellValue(barlist, "I1", "Маркировка ДМО")      // 9
	_ = f.SetCellValue(barlist, "J1", "Длина ДМО, мм")       // 10
	_ = f.SetCellValue(barlist, "K1", "Масса ДМО, кг")       // 11

	row := 2

	for _, b := range allBars {
		cell, _ := excelize.CoordinatesToCellName(1, row)
		_ = f.SetCellValue(barlist, cell, b.NestName())

		cell, _ = excelize.CoordinatesToCellName(2, row)
		_ = f.SetCellValue(barlist, cell, "Камера резки профиля")

		cell, _ = excelize.CoordinatesToCellName(3, row)
		_ = f.SetCellValue(barlist, cell, b.UsedRemnant)

		cell, _ = excelize.CoordinatesToCellName(4, row)
		_ = f.SetCellValue(barlist, cell, b.Dim())

		cell, _ = excelize.CoordinatesToCellName(5, row)
		_ = f.SetCellValue(barlist, cell, b.Quality())

		cell, _ = excelize.CoordinatesToCellName(6, row)
		_ = f.SetCellValue(barlist, cell, b.Capacity)

		cell, _ = excelize.CoordinatesToCellName(7, row)
		_ = f.SetCellValue(barlist, cell, mass(b.Dim(), b.Capacity))

		cell, _ = excelize.CoordinatesToCellName(8, row)
		_ = f.SetCellValue(barlist, cell, b.Length/b.Capacity*100)

		if b.Capacity-b.Length > 1000 {
			cell, _ = excelize.CoordinatesToCellName(9, row)
			_ = f.SetCellValue(barlist, cell, b.NestName()+"R01")

			cell, _ = excelize.CoordinatesToCellName(10, row)
			_ = f.SetCellValue(barlist, cell, b.RemnantLength())

			cell, _ = excelize.CoordinatesToCellName(11, row)
			_ = f.SetCellValue(barlist, cell, mass(b.Dim(), b.RemnantLength()))
		}

		row++
	}

	f.SetActiveSheet(f.NewSheet(partlist))

	_ = f.SetCellValue(partlist, "A1", "Чертеж")             // 1
	_ = f.SetCellValue(partlist, "B1", "Заказ")              // 2
	_ = f.SetCellValue(partlist, "C1", "Секция")             // 3
	_ = f.SetCellValue(partlist, "D1", "Позиция")            // 4
	_ = f.SetCellValue(partlist, "E1", "Карта раскроя")      // 5
	_ = f.SetCellValue(partlist, "F1", "Типоразмер детали")  // 6
	_ = f.SetCellValue(partlist, "G1", "Марка материала")    // 7
	_ = f.SetCellValue(partlist, "H1", "Кол-во, шт")         // 8
	_ = f.SetCellValue(partlist, "I1", "Длина, мм")          // 9
	_ = f.SetCellValue(partlist, "J1", "Масса 1 детали, кг") // 10
	_ = f.SetCellValue(partlist, "K1", "Общая масса, кг")    // 11
	_ = f.SetCellValue(partlist, "L1", "Маршрут обработки")  // 12
	_ = f.SetCellValue(partlist, "M1", "Примечание")         // 13

	row = 2

	for _, b := range allBars {
		for _, p := range b.Parts {
			//cell, _ := excelize.CoordinatesToCellName(1, row)
			//_ = f.SetCellValue(partlist, cell, b.NestName())

			cell, _ := excelize.CoordinatesToCellName(2, row)
			_ = f.SetCellValue(partlist, cell, p.Project)

			cell, _ = excelize.CoordinatesToCellName(3, row)
			_ = f.SetCellValue(partlist, cell, p.Section)

			cell, _ = excelize.CoordinatesToCellName(4, row)

			pos, _ := strconv.Atoi(p.PosNo)
			_ = f.SetCellValue(partlist, cell, pos)

			cell, _ = excelize.CoordinatesToCellName(5, row)
			_ = f.SetCellValue(partlist, cell, b.NestName())

			cell, _ = excelize.CoordinatesToCellName(6, row)
			_ = f.SetCellValue(partlist, cell, p.Dim)

			cell, _ = excelize.CoordinatesToCellName(7, row)
			_ = f.SetCellValue(partlist, cell, p.Quality)

			cell, _ = excelize.CoordinatesToCellName(8, row)
			_ = f.SetCellValue(partlist, cell, p.Quantity)

			cell, _ = excelize.CoordinatesToCellName(9, row)
			_ = f.SetCellValue(partlist, cell, p.Length)

			cell, _ = excelize.CoordinatesToCellName(10, row)
			_ = f.SetCellValue(partlist, cell, mass(b.Dim(), p.Length))

			cell, _ = excelize.CoordinatesToCellName(11, row)
			_ = f.SetCellValue(partlist, cell, mass(b.Dim(), p.Length)*float64(p.Quantity))

			row++
		}
	}

	if err := f.SaveAs(path.Join("out", nm+".xlsx")); err != nil {
		fmt.Println(err)
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
