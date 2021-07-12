package formats

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/edanko/moses/entities"
)

var (
	reProfType = regexp.MustCompile(`(?m)PP\/(\d+)x(\d+\.?\d?)`)
	reQuality  = regexp.MustCompile(`(?m)Qual :\s+(\w+)`)
)

func ProcessCsv(fname string) (map[string]*entities.Profile, error) {
	profs := make(map[string]*entities.Profile)

	dim, qual := dimAndQuality(fname)

	f, err := os.Open(fname)
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

	/* 	bev := func(f string) (bev string) {
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
	} */

	for {
		l, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		p := &entities.Profile{}
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

		if p, exists := profs[p.Project+p.Section+p.PosNo]; exists {
			p.Quantity++
			continue
		}

		p.Quantity = 1

		p.L = &entities.End{}
		p.L.Name = l[5]

		params := make(map[string]float64, 10)

		if n := stof(l[6]); n > 0 {
			params["A"] = n
		}
		if n := stof(l[7]); n > 0 {
			params["B"] = n
		}
		if n := stof(l[8]); n > 0 {
			params["C"] = n
		}
		if n := stof(l[9]); n > 0 {
			params["Ks"] = n
		}
		if n := stof(l[10]); n > 0 {
			params["R1"] = n
		}
		if n := stof(l[11]); n > 0 {
			params["R2"] = n
		}
		if n := stof(l[12]); n > 0 {
			params["V1"] = n
		}
		if n := stof(l[13]); n > 0 {
			params["V2"] = n
		}
		if n := stof(l[14]); n > 0 {
			params["V3"] = n
		}
		if n := stof(l[15]); n > 0 {
			params["V4"] = n
		}

		p.L.Params = params
		//lb := bev(l[16])

		p.R = &entities.End{}
		p.R.Name = l[19]

		params = make(map[string]float64, 10)

		if n := stof(l[20]); n > 0 {
			params["A"] = n
		}
		if n := stof(l[21]); n > 0 {
			params["B"] = n
		}
		if n := stof(l[22]); n > 0 {
			params["C"] = n
		}
		if n := stof(l[23]); n > 0 {
			params["Ks"] = n
		}
		if n := stof(l[24]); n > 0 {
			params["R1"] = n
		}
		if n := stof(l[25]); n > 0 {
			params["R2"] = n
		}
		if n := stof(l[26]); n > 0 {
			params["V1"] = n
		}
		if n := stof(l[27]); n > 0 {
			params["V2"] = n
		}
		if n := stof(l[28]); n > 0 {
			params["V3"] = n
		}
		if n := stof(l[29]); n > 0 {
			params["V4"] = n
		}

		p.R.Params = params

		//rb := bev(l[30])

		p.Source = path.Base(fname)

		profs[p.Project+p.Section+p.PosNo] = p
	}
	return profs, err
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
			dim = "RP" + tmp[0][1] + "*" + tmp[0][2]
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

	return dim, quality
}

func stof(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}
	return f
}
