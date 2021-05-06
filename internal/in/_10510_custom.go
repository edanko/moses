package in

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/edanko/moses/internal/models"
)

func ProcessAvevaCustomReportCsv(file string) map[string][]*models.Part {
	parts := make(map[string][]*models.Part, 1000)

	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ';'

	// skip header
	_, err = r.Read()
	if err != nil {
		panic(err)
	}

	for {
		l, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		p := &models.Part{
			Project: "056001",
			Section: l[0],
			PosNo:   l[1],
			Dim:     l[2],
			Quality: l[3],
			LEnd:    l[5],
			REnd:    l[6],
		}

		switch p.Quality {
		case "D40_B":
			p.Quality = "D40-B"
		}

		p.Length, err = strconv.ParseFloat(l[4], 64)
		if err != nil {
			panic(err)
		}

		if p.Section == "" || p.PosNo == "" || p.Dim == "" || p.Quality == "" {
			continue
		}

		// skip os bevel
		if strings.Contains(p.LEnd, "--") {
			p.LEnd = p.LEnd[:2] + p.LEnd[7:]
		}
		if strings.Contains(p.REnd, "--") {
			p.REnd = p.REnd[:2] + p.REnd[7:]
		}

		p.SetFullLength()

		parts[p.Dim] = append(parts[p.Dim], p)
	}
	return parts
}
