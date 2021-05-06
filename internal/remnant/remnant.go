package remnant

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/edanko/moses/internal/models"
	"github.com/edanko/moses/internal/utils"

	"github.com/jszwec/csvutil"
	"github.com/spf13/viper"
)

type Remnant struct {
	Id          int     `csv:"id"`
	Prj         string  `csv:"prj"`
	Sec         string  `csv:"sec"`
	Dim         string  `csv:"dim"`
	Quality     string  `csv:"quality"`
	RemLen      float64 `csv:"rem_len"`
	MarkingText string  `csv:"marking_text"`
	UsedIn      string  `csv:"used_in"`
}

var (
	remnants      []*Remnant
	lastRemnantId int
	d, _          = os.Getwd()
	remnantPath   = filepath.Join(d, "remnants.csv")
)

const (
	useRemnants = true
	delim       = ';'
)

func init() {
	viper.Set("useremants", useRemnants)
	//remnantPath = viper.GetString("remnantpath")

	if !useRemnants {
		fmt.Println("[x] not using remnants")
		return
	}
	remnants = loadRemanantsFromCSV()
	if len(remnants) == 0 {
		fmt.Println("[+] no remnants loaded")
		lastRemnantId = 0
	} else {
		fmt.Println("[+]", len(remnants), "remnants loaded")
		lastRemnantId = remnants[len(remnants)-1].Id
	}
}

func loadRemanantsFromCSV() []*Remnant {
	rems := make([]*Remnant, 0, 1000)

	f, err := os.Open(remnantPath)
	if err != nil {
		fmt.Println("[!] failed to open", remnantPath)
		return rems
	}
	defer f.Close()

	csvReader := csv.NewReader(bufio.NewReader(f))
	csvReader.Comma = delim

	dec, err := csvutil.NewDecoder(csvReader)
	if err != nil {
		log.Fatal(err)
	}

	for {
		var r Remnant
		if err := dec.Decode(&r); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		r.Id, err = strconv.Atoi(dec.Record()[0])
		if err != nil {
			panic(err)
		}

		rems = append(rems, &r)
	}
	return rems
}

func Remnants(prj, dim, quality string) []*models.Bar {
	rems := make([]*models.Bar, 0)

	for _, r := range remnants {
		if r.UsedIn != "" {
			continue
		}

		if r.Prj != prj || r.Dim != strings.ToUpper(dim) || r.Quality != strings.ToUpper(quality) {
			continue
		}

		rems = append(rems, &models.Bar{
			Capacity:      r.RemLen,
			Length:        0,
			MarkRemnant:   "",
			UsedRemnantID: r.Id,
			UsedRemnant:   r.MarkingText,
			Parts:         make([]*models.Part, 0),
		})
	}
	return rems
}

func WriteRemnantsToFile() error {
	b, err := csvutil.Marshal(remnants)
	if err != nil {
		return err
	}

	b = bytes.ReplaceAll(b, []byte(","), []byte(string(delim)))

	fmt.Println("[+]", len(remnants), "remnants will be stored")

	return utils.WriteStringToFile(remnantPath, string(b))
}

func Add(prj, sec, dim, qual, length, mark string) {

	lastRemnantId++

	r := &Remnant{
		Id:          lastRemnantId,
		Prj:         prj,
		Sec:         sec,
		Dim:         dim,
		Quality:     qual,
		MarkingText: mark,
	}

	var err error
	r.RemLen, err = strconv.ParseFloat(length, 64)
	if err != nil {
		panic(err)
	}

	remnants = append(remnants, r)
}

func Use(id int, mark string) {
	for _, r := range remnants {
		if r.Id == id {
			r.UsedIn = mark
		}
	}
}
