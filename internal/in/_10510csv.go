package in

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/edanko/moses/internal/models"
)

const (
	prj = "056001"
)

var (
	reProfType = regexp.MustCompile(`(?m)PP\/(\d+)x(\d+\.?\d?)`)
	reQuality  = regexp.MustCompile(`(?m)Qual :\s+(\w+)`)
)

func ProcessAvevaCsv(csvFile string) []*models.Part {
	var allProfiles []*models.Part

	dim, qual := dimAndQuality(csvFile)

	f, err := os.Open(csvFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := csv.NewReader(f)

	// skip header
	_, err = r.Read()
	if err != nil {
		panic(err)
	}

	bev := func(f string) (bev string) {
		if f == "" || f == "0" {
			return
		}

		switch f {
		case "177", "154", "180", "227", "143", "-245":
			bev = ""

		case "-131":
			bev = "-131"

		default:
			bev = "-" + f
		}

		return
	}

	for {
		l, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		p := &models.Part{}
		name := strings.Split(l[0], "-") // sec, sec2?, assembly, partid with side

		p.Section = name[0]
		p.PosNo = name[len(name)-1][:len(name[len(name)-1])-1]

		p.Length, err = strconv.ParseFloat(l[3], 64)
		if err != nil {
			panic(err)
		}

		p.Dim = dim
		p.Quality = qual
		p.Project = prj

		lEnd := strings.Builder{}
		lEnd.Grow(100)
		lEnd.WriteString(l[5])

		lb := bev(l[16])
		if lb != "" {
			lEnd.WriteString(lb)
		}

		// bevel web
		/*if rec[16] != "" && rec[16] != "0" {
			lEnd.WriteString("-")
			lEnd.WriteString(rec[16])
		}*/

		// A
		if l[6] != "" {
			lEnd.WriteString(" a=")
			lEnd.WriteString(l[6])
		}
		// B
		if l[7] != "" {
			lEnd.WriteString(" b=")
			lEnd.WriteString(l[7])
		}
		// C
		if l[8] != "" {
			lEnd.WriteString(" c=")
			lEnd.WriteString(l[8])
		}
		// Snipe
		if l[9] != "" {
			lEnd.WriteString(" snipe=")
			lEnd.WriteString(l[9])
		}
		// R1
		if l[10] != "" {
			lEnd.WriteString(" r1=")
			lEnd.WriteString(l[10])
		}
		// R2
		if l[11] != "" {
			lEnd.WriteString(" r2=")
			lEnd.WriteString(l[11])
		}
		// V1
		if l[12] != "" {
			lEnd.WriteString(" v1=")
			lEnd.WriteString(l[12])
		}
		// V2
		if l[13] != "" {
			lEnd.WriteString(" v2=")
			lEnd.WriteString(l[13])
		}
		// V3
		if l[14] != "" {
			lEnd.WriteString(" v3=")
			lEnd.WriteString(l[14])
		}

		rEnd := strings.Builder{}
		rEnd.Grow(100)
		rEnd.WriteString(l[19])

		rb := bev(l[30])
		if rb != "" {
			rEnd.WriteString(rb)
		}

		// bevel web
		/*if rec[30] != "" && rec[30] != "0" {
			rEnd.WriteString("-")
			rEnd.WriteString(rec[30])
		}*/

		// A
		if l[20] != "" {
			rEnd.WriteString(" a=")
			rEnd.WriteString(l[20])
		}
		// B
		if l[21] != "" {
			rEnd.WriteString(" b=")
			rEnd.WriteString(l[21])
		}
		// C
		if l[22] != "" {
			rEnd.WriteString(" c=")
			rEnd.WriteString(l[22])
		}
		// Snipe
		if l[23] != "" {
			rEnd.WriteString(" snipe=")
			rEnd.WriteString(l[23])
		}
		// R1
		if l[24] != "" {
			rEnd.WriteString(" r1=")
			rEnd.WriteString(l[24])
		}
		// R2
		if l[25] != "" {
			rEnd.WriteString(" r2=")
			rEnd.WriteString(l[25])
		}
		// V1
		if l[26] != "" {
			rEnd.WriteString(" v1=")
			rEnd.WriteString(l[26])
		}
		// V2
		if l[27] != "" {
			rEnd.WriteString(" v2=")
			rEnd.WriteString(l[27])
		}
		// V3
		if l[28] != "" {
			rEnd.WriteString(" v3=")
			rEnd.WriteString(l[28])
		}

		p.LEnd = lEnd.String()
		p.REnd = rEnd.String()

		p.SetFullLength()

		allProfiles = append(allProfiles, p)
	}
	return allProfiles
}

func dimAndQuality(fname string) (string, string) {
	var dim, quality string

	fname = strings.Replace(fname, ".csv", ".lst", 1)

	f, err := os.Open(fname)
	if err != nil {
		panic(err)
	}

	r := bufio.NewReader(f)

	for {
		str, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		if strings.Contains(str, "Type/Dim.") {
			tmp := reProfType.FindAllStringSubmatch(str, -1)
			dim = "RP" + tmp[0][1] + "X" + tmp[0][2]
			continue
		}

		if strings.Contains(str, "Qual") {
			tmp := reQuality.FindAllStringSubmatch(str, -1)
			quality = tmp[0][1]
			continue
		}

		if dim != "" && quality != "" {
			break
		}
	}

	if dim[len(dim)-2:] == ".0" {
		dim = dim[:len(dim)-2]
	}

	switch quality {
	case "D40_B":
		quality = "D40-B"
	}

	return dim, quality
}
