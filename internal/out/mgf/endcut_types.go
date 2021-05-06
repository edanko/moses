package mgf

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

func Params(t, end string) (string, string) {

	e := strings.Split(end, " ")

	p := make(map[string]float64)

	for i := 1; i < len(e); i++ {
		s := strings.Split(e[i], "=")

		val, err := strconv.ParseFloat(s[1], 64)
		if err != nil {
			panic(err)
		}

		p[s[0]] = val
	}

	if t[:2] == "rp" {

		spl := strings.Split(t[2:], "x")
		t, _ := strconv.ParseFloat(spl[1], 64)
		h, _ := strconv.ParseFloat(spl[0], 64)

		switch e[0] {
		case "0":
			return e21(p["angle"], 0, 0, 0, 0, 0, 0), ""

		case "bk":
			if p["angle"] > 0 {
				return e21(p["angle"], 0, 0, 0, 0, 0, 0), sc12(p["h"], p["h"], 0, 0, p["angle"])
			} else {
				return e21(p["angle"], 0, p["h"], 0, 0, 0, 0), ""
			}

		case "0sc":
			return e21r(0, 0, p["r"], 0, 0, 0, 0), ""

		case "e23f":
			return e23f(p["a"], p["b"]), ""

		case "e23f_sc":
			return e23f(p["a"], p["b"]), sc12(p["r"], p["r"], 0, p["r"], p["angle"])

		case "fl":
			return e21r(p["angle"], 45, p["r"], 0, 0, 0, 0), ""

		case "fv":
			return e21r(p["angle"], 45, p["r"], p["a"], t-2, 0, 0), ""

		case "fv.1":
			return e21r(p["angle"], 45, p["r"], 0, 0, 0, 0), ""

		case "sb1":
			return e213rs(0, 30, 30, p["r"], p["r"], 0, p["r"], -1, 30, t-2, 0, 0), ""

		case "sb1.1":
			return e21r(p["angle"], 0, p["r"], 0, 0, 0, 0), ""

		case "sb2":
			return e213rs(0, 45, 45, p["r"], p["r"], 0, p["r"], -1, 45, t-2, 0, 0), ""

		case "sb3", "sb20":
			return e21r(p["angle"], 45, p["r"], 30, t-2, 0, 0), ""

		case "sb3.1", "sb20.1":
			return e21r(p["angle"], 45, p["r"], 0, 0, 0, 0), ""

		case "sb3v":
			return e22(p["angle"], 45, p["r"], p["r"]+p["u"], p["u"], p["r"], 0, 30, t-2, 0, 0), ""

		case "sb4", "sb4-01":
			return e27r(0, 0, p["down"], p["up"], p["up"], 0, 45, t-2, 0, 0), ""

		case "sb4_2":
			return e27(0, 0, p["e"], p["up"], p["up"], 0, 45, t-2, 0, 0), ""

		case "sb5", "sb5-01":
			return e21r(p["angle"], 45, p["r"], 45, t-2, 0, 0), ""

		case "sb5_2":
			return e21(p["angle"], 45, p["e"], 45, t-2, 0, 0), ""

		case "sb5v":
			return e22(p["angle"], 45, p["r"], p["r"]+p["u"], p["u"], p["r"], 0, 45, t-2, 0, 0), ""

		case "sb6":
			return e213rs(0, 45, 45, p["r"], p["r"], 0, p["r"], -1, 45, t-2, 0, 0), ""

		case "sb6.1":
			return e21r(p["angle"], 0, p["r"], 0, 0, 0, 0), ""

		case "sb8":
			return e210rs(p["angle"], 30, 30, 0, 0, 0, 0, 13, 30, t-2, 0, 0), ""

		case "sb8.1":
			return e210rs(p["angle"], 30, 30, 0, 0, 0, 0, 13, 0, 0, 0, 0), ""

		case "sb9":
			return e213rs(p["angle"], 45, 45, p["r"], p["r"], 0, p["r"], 15, 0, 0, 0, 0), ""

		case "sb9-01":
			return e213rs(p["angle"], 45, 45, p["r"], p["r"], 0, p["r"], 15, 45, t-2, 0, 0), ""

		case "sb10":
			return e28r(0, 0, p["r"], 1, 0, p["a"], h-20, 0, 0, 0, 0), ""

		case "sb11":
			return e21r(0, 45, p["r"], 0, 0, 0, 0), ""

		case "sb11_2":
			return e21(0, 45, p["e"], 0, 0, 0, 0), ""

		case "sc50":
			return e22(0, 0, p["r"]+50, p["r"], 0, p["r"], 0, 0, 0, 0, 0), ""

		case "bl":
			return e21r(p["angle"], 0, p["r"], 0, 0, 0, 0), ""

		case "blg":
			return e21r(p["angle"], 0, p["r"], 45, t-2, 0, 0), ""

		case "br":
			if p["h"] != 0 {
				return e29(p["angle"], 0, 5, 5, 0, 30, d(h, p["bulb"]), p["bulb"], 0, 0, 0, 0), sc12(p["h"], p["h"], 0, 0, p["angle"])
			} else {
				return e29(p["angle"], 0, 5, 5, 0, 30, d(h, p["bulb"]), p["bulb"], 0, 0, 0, 0), ""
			}

		case "brbz":
			if p["h"] != 0 {
				return e29(p["angle"], p["angle2"], 5, 5, 0, 30, d(h, p["bulb"]), p["bulb"], 0, 0, 0, 0), sc12(p["h"], p["h"], 0, 0, p["angle"])
			} else {
				return e29(p["angle"], p["angle2"], 5, 5, 0, 30, d(h, p["bulb"]), p["bulb"], 0, 0, 0, 0), ""
			}

		case "bvbz":
			return e221rs(p["angle"], p["angle2"], p["b"], p["h"], p["h"]-p["r"], p["r"], p["h2"]-t, -1, 0, 0, 0, 0), ""

		case "ba":
			return e21(p["angle"], 0, 0, 0, 0, 0, 0), ""

		case "bc":
			return e23(0, 90-p["angle"], 0, 0, 0, 0, 0, 0, 0), ""

		case "bck":
			return e23(p["angle"], 90-p["angle2"], 0, p["b"], p["h"], 0, 0, 0, 0), ""

		case "bcr":
			return e23r(p["angle"], 90-p["angle2"], 0, p["r"], p["h"], 0, 0, 0, 0), ""

		case "bac":
			return e23(0, bac(h, p["3"], p["2"]), 0, 0, p["2"], 0, 0, 0, 0), sc12(p["1"], p["2"], 0, 0, p["angle"])

		case "bv":
			return e22(p["angle"], 0, p["b"], p["h"], p["h"]-p["r"], p["r"], 0, 0, 0, 0, 0), ""

		case "bvf":
			if p["r"] == 0 {
				return e21(p["angle"], 25, 0, 0, 0, 0, 0), ""
			} else {
				return e22(p["angle"], 25, p["b"], p["h"], p["h"]-p["r"], p["r"], 0, 0, 0, 0, 0), ""
			}

		case "bvg":
			return e22(p["angle"], 0, p["b"], p["h"], p["h"]-p["r"], p["r"], 0, 0, 0, 0, 0), ""

		case "bw":
			return e29(p["angle"], 0, 5, 5, 0, 30, d(h, p["bulb"]), p["bulb"], 0, 0, 0, 0), sc12(p["r"], p["r"], 0, p["r"], p["angle"])

		case "btn":
			if p["h"] != 0 {
				return e29(p["angle"], 0, 5, 5, 0, 30, d(h, p["bulb"]), p["bulb"], 0, 0, 0, 0), sc12(p["l"], p["h"], 0, p["r"], p["angle"])
			} else {
				return e29(p["angle"], 0, 5, 5, 0, 30, d(h, p["bulb"]), p["bulb"], 0, 0, 0, 0), ""
			}

		case "bmn":
			return e22(p["angle"], 0, p["b"], p["r"], 0, p["r"], 0, 0, 0, 0, 0), ""

		case "bmnbz":
			return e221rs(p["angle"], p["angle2"], p["b"], p["r"], 0, p["r"], p["h2"]-t, -1, 0, 0, 0, 0), ""

		case "blgbz":
			return e221rs(p["angle"], p["angle2"], p["r"], p["r"], 0, p["r"], p["h2"]-t, -1, 45, t-2, 0, 0), ""

		case "blbz":
			return e221rs(p["angle"], p["angle2"], p["r"], p["r"], 0, p["r"], p["h2"]-t, -1, 0, 0, 0, 0), ""

		case "bkbz":
			return e221rs(p["angle"], p["angle2"], p["h"], p["h"], 0, 0, p["h2"]-t, -1, 0, 0, 0, 0), ""

		case "babz":
			return e221rs(p["angle"], p["angle2"], 0, 0, 0, 0, p["h2"]-t, -1, 0, 0, 0, 0), ""

		case "bxaf", "bxag":
			var c, r2 float64

			if p["r"] > 0 {
				c = p["h"] - p["r"]
			}

			if t/2 < 1 {
				r2 = 1
			} else {
				r2 = t / 2
			}

			return e372rs(25, p["l"], p["h"], c, p["r"], h/4, l(h, p["hprof"]), h-p["hprof"], 1, r2, 0, 0, 0, 0, 0, 0, 0), ""

		case "0_45":
			return e21(0, 45, 0, 0, 0, 0, 0), ""

		case "1":
			return e27r(0, 0, p["down"], p["up"], p["up"], 0, 0, 0, 0, 0), ""

		case "b-102":
			return e213rs(p["angle"], 30, 30, p["r"], p["r"], 0, p["r"], 12, 30, t-2, 0, 0), ""

		case "21":
			if p["snipe"] != 0 {
				var e1 float64
				if p["angle"] > 90 {
					e1 = 0
				} else {
					e1 = p["snipe"]
				}
				return e21(p["v1"], p["v2"], e1, 0, 0, 0, 0), sc12(p["snipe"], p["snipe"], 0, 0, p["v1"])
			} else {
				return e21r(p["v1"], p["v2"], p["r1"], 0, 0, 0, 0), ""
			}

		case "21-116", "21-216", "21-145", "21--145", "21-173", "21--173", "21-245", "21-273", "21-275", "21-346":
			if p["snipe"] != 0 {
				var e1 float64
				if p["angle"] > 90 {
					e1 = 0
				} else {
					e1 = p["snipe"]
				}
				return e21(p["v1"], p["v2"], e1, 45, t-2, 0, 0), sc12(p["snipe"], p["snipe"], 0, 0, p["v1"])
			} else {
				return e21r(p["v1"], p["v2"], p["r1"], 45, t-2, 0, 0), ""
			}

		case "21-131", "21-132":
			return e213rs(p["v1"], 25, 25, p["r1"], p["r1"], 0, p["r1"], -1, 25, t-2, 0, 0), ""

		case "21--132":
			return e213rs(p["v1"], 25, 25, p["r1"], p["r1"], 0, p["r1"], -1, 0, 0, 25, t-2), ""

		case "23":
			if p["r"] > 0 {
				return e23r(p["v1"], p["v3"], p["v2"], p["r1"], p["b"], 0, 0, 0, 0), ""
			} else {
				return e23(p["v1"], p["v3"], p["v2"], 0, p["b"], 0, 0, 0, 0), ""
			}

		case "27":
			return e27r(p["v1"], p["v2"], p["r1"], p["r2"], p["a"], 0, 0, 0, 0, 0), ""

		default:
			fmt.Println("unknown endcut:", e[0])
			//return "", ""
		}
	}
	if t[:2] == "fb" {
		switch e[0] {
		case "ba":
			return e11(p["angle"], 0, 0, 0, 0, 0), ""

		case "bk":
			return e11(p["angle"], p["h"], 0, 0, 0, 0), ""

		case "baosv":
			return e11(p["angle"], 0, 0, 0, p["deg"], 2), ""

		case "batsv":
			return e11(p["angle"], 0, p["deg"], 2, 0, 0), ""

		case "bc":
			return e13(90, p["angle"], 0, p["h"], 0, 0, 0, 0), ""

		case "bck":
			return e13(p["angle"], p["angle2"], p["b"], p["h"], 0, 0, 0, 0), ""

		default:
			return "", ""
		}
	}

	return "", ""

}

func IcutString(s string) string {

	e := strings.Split(s, " ")

	params := make(map[string]float64)

	for i := 1; i < len(e); i++ {
		s := strings.Split(e[i], "=")

		val, err := strconv.ParseFloat(s[1], 64)
		if err != nil {
			panic(err)
		}

		params[s[0]] = val
	}

	switch e[0] {
	case "ae_half_hole", "hole", "i1":
		return i1(params["x"], params["r"])

	case "notch":
		return i1(params["u"], params["r"])

	case "ae_hole", "i2":
		return i2(params["x"], params["y"], params["r"])

	case "hc":
		return i4(params["x"], params["a"], params["b"], params["r"])

	case "hin":
		return i11(params["x"], params["l"], params["b"])

	}
	return ""

}

func d(height, bulb float64) float64 {
	var sub float64
	switch bulb {
	default:
		log.Fatalln("unknown bulb", bulb, "in br or bw template")
	case 21:
		sub = 20.01
	case 24:
		sub = 21.74
	case 27:
		sub = 25.05
	case 29:
		sub = 27.78
	case 32:
		sub = 29.52
	case 35:
		sub = 32.83
	case 38:
		sub = 35.35
	}

	return height - sub
}

func bac(h, arg1, arg2 float64) float64 {
	return 90 - math.Atan(arg1/(h-arg2))*180/math.Pi
}

func l(height, hprof float64) float64 {
	switch height - hprof {
	case 20, 24:
		return 100
	case 40:
		return 188
	case 60:
		return 260
	case 80:
		return 290
	case 100:
		return 374
	case 120:
		return 449
	default:
		log.Fatalln("unknown height and hprof difference:", height-hprof)
	}
	return 0
}

// get angle for cx template
func cx_doc(h, t, legUp float64) float64 {
	if legUp == 0 {
		legUp = h
	}
	legSide := (h - (3 * t)) / 2

	return 90 - math.Atan(legUp/legSide)*180/math.Pi
}

// for fl bv template
func a(h, angle float64) float64 {
	if angle == 0 {
		return 0
	}

	var sign float64

	// just inverted direction
	if angle < 0 {
		sign = 1
	} else {
		sign = -1
	}

	return sign * h * math.Sin(math.Abs(angle)*math.Pi/180)
}
