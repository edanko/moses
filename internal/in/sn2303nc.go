package in

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/edanko/moses/internal/models"
)

func ProcessNc(fname string) *models.Part {
	p := models.Part{}

	fmt.Println(fname)

	f, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	s.Scan()
	if s.Text()[:3] != "PRF" {
		fmt.Println(fname, "not a profile")
		return nil
	}

	for s.Scan() {
		switch s.Text()[:3] {
		case "POS":
			pos := strings.Split(s.Text(), ",")
			p.Section = pos[2]
			p.PosNo = pos[3]

		case "BEW":
			bew := strings.Split(s.Text(), ",")
			_ = bew[1]

		case "NUM":
			n := strings.Split(s.Text(), ",")
			p.Quantity, err = strconv.Atoi(n[1])
			if err != nil {
				panic(err)
			}

		case "TYP":
			t := strings.Split(s.Text(), ",")

			switch t[1] {
			case "HP":
				p.Dim = "rp" + t[2] + "x" + t[3]
				p.Quality = t[6]
			default:
				fmt.Println("not hp here", p.PosNo)
			}

		case "LEN":
			l := strings.Split(s.Text(), ",")
			p.Length, err = strconv.ParseFloat(l[1], 64)
			if err != nil {
				panic(err)
			}

		case "MIS":
			m := strings.SplitN(s.Text(), ",", 2)
			p.LEnd = parsePlfEnd(m[1])

		case "MIE":
			m := strings.SplitN(s.Text(), ",", 2)
			p.REnd = parsePlfEnd(m[1])

		case "MIT":
			continue

		case "MII":
			l := strings.SplitN(s.Text(), ",", 2)
			res := parsePlfInner(l[1])
			p.Icuts = append(p.Icuts, res)

		case "END":
			break
		}
	}

	//p.SetFullLength()
	return &p
}
