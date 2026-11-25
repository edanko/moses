package mgf

import (
	"strconv"
	"strings"
)

func e11(v1, e1, gamma, ho, gammau, hu float64) string {
	s := strings.Builder{}

	if v1 > 45 {
		v1 = 45
	}
	if v1 < -45 {
		v1 = -45
	}

	s.WriteString("ENDCUT_TYPE=11")
	s.WriteString(newLine)

	s.WriteString("START_OF_PARAMS")
	s.WriteString(newLine)

	s.WriteString("NO_OF_PARAMS=6")
	s.WriteString(newLine)

	s.WriteString("V1=")
	s.WriteString(strconv.FormatFloat(90+v1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("E1=")
	s.WriteString(strconv.FormatFloat(e1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMA=")
	s.WriteString(strconv.FormatFloat(gamma, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HO=")
	s.WriteString(strconv.FormatFloat(ho, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMAU=")
	s.WriteString(strconv.FormatFloat(gammau, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HU=")
	s.WriteString(strconv.FormatFloat(hu, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("END_OF_PARAMS")
	s.WriteString(newLine)

	return s.String()
}

func sc12(a, b, c, r1, v1 float64) string {
	s := strings.Builder{}

	if v1 > 45 {
		v1 = 45
	}
	if v1 < -45 {
		v1 = -45
	}

	s.WriteString("SCALLOP_TYPE=12")
	s.WriteString(newLine)

	s.WriteString("START_OF_PARAMS")
	s.WriteString(newLine)

	s.WriteString("NO_OF_PARAMS=5")
	s.WriteString(newLine)

	s.WriteString("V1=")
	s.WriteString(strconv.FormatFloat(90+v1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("A=")
	s.WriteString(strconv.FormatFloat(a, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("B=")
	s.WriteString(strconv.FormatFloat(b, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("C=")
	s.WriteString(strconv.FormatFloat(c, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("R1=")
	s.WriteString(strconv.FormatFloat(r1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("END_OF_PARAMS")
	s.WriteString(newLine)

	return s.String()
}

func e13(v1, v2, e1, b, gamma, ho, gammau, hu float64) string {
	s := strings.Builder{}

	if v1 > 45 {
		v1 = 45
	}
	if v1 < -45 {
		v1 = -45
	}

	if v2 > 45 {
		v2 = 45
	}
	if v2 < 0 {
		v2 = 0
	}

	s.WriteString("ENDCUT_TYPE=13")
	s.WriteString(newLine)

	s.WriteString("START_OF_PARAMS")
	s.WriteString(newLine)

	s.WriteString("NO_OF_PARAMS=8")
	s.WriteString(newLine)

	s.WriteString("V1=")
	s.WriteString(strconv.FormatFloat(90+v1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("V2=")
	s.WriteString(strconv.FormatFloat(90-v2, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("E1=")
	s.WriteString(strconv.FormatFloat(e1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("B=")
	s.WriteString(strconv.FormatFloat(b, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMA=")
	s.WriteString(strconv.FormatFloat(gamma, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HO=")
	s.WriteString(strconv.FormatFloat(ho, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMAU=")
	s.WriteString(strconv.FormatFloat(gammau, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HU=")
	s.WriteString(strconv.FormatFloat(hu, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("END_OF_PARAMS")
	s.WriteString(newLine)

	return s.String()
}

func e21(v1, v3, e1, gamma, ho, gammau, hu float64) string {
	s := strings.Builder{}

	if v1 > 45 {
		v1 = 45
	}
	if v1 < -45 {
		v1 = -45
	}

	if v3 > 45 {
		v3 = 45
	}
	if v3 < 0 {
		v3 = 0
	}

	s.WriteString("ENDCUT_TYPE=21")
	s.WriteString(newLine)

	s.WriteString("START_OF_PARAMS")
	s.WriteString(newLine)

	s.WriteString("NO_OF_PARAMS=7")
	s.WriteString(newLine)

	s.WriteString("V1=")
	s.WriteString(strconv.FormatFloat(90+v1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("V3=")
	s.WriteString(strconv.FormatFloat(90-v3, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("E1=")
	s.WriteString(strconv.FormatFloat(e1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMA=")
	s.WriteString(strconv.FormatFloat(gamma, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HO=")
	s.WriteString(strconv.FormatFloat(ho, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMAU=")
	s.WriteString(strconv.FormatFloat(gammau, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HU=")
	s.WriteString(strconv.FormatFloat(hu, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("END_OF_PARAMS")
	s.WriteString(newLine)

	return s.String()
}

func e21r(v1, v3, r, gamma, ho, gammau, hu float64) string {
	s := strings.Builder{}

	if v1 > 45 {
		v1 = 45
	}
	if v1 < -45 {
		v1 = -45
	}

	if v3 > 45 {
		v3 = 45
	}
	if v3 < 0 {
		v3 = 0
	}

	s.WriteString("ENDCUT_TYPE=21R")
	s.WriteString(newLine)

	s.WriteString("START_OF_PARAMS")
	s.WriteString(newLine)

	s.WriteString("NO_OF_PARAMS=7")
	s.WriteString(newLine)

	s.WriteString("V1=")
	s.WriteString(strconv.FormatFloat(90+v1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("V3=")
	s.WriteString(strconv.FormatFloat(90-v3, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("R1=")
	s.WriteString(strconv.FormatFloat(r, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMA=")
	s.WriteString(strconv.FormatFloat(gamma, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HO=")
	s.WriteString(strconv.FormatFloat(ho, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMAU=")
	s.WriteString(strconv.FormatFloat(gammau, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HU=")
	s.WriteString(strconv.FormatFloat(hu, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("END_OF_PARAMS")
	s.WriteString(newLine)

	return s.String()
}

func e22(v1, v3, a, b, c, r1, e1, gamma, ho, gammau, hu float64) string {
	s := strings.Builder{}

	if v1 > 45 {
		v1 = 45
	}
	if v1 < -45 {
		v1 = -45
	}
	if v3 > 45 {
		v3 = 45
	}
	if v3 < -45 {
		v3 = -45
	}

	s.WriteString("ENDCUT_TYPE=22")
	s.WriteString(newLine)

	s.WriteString("START_OF_PARAMS")
	s.WriteString(newLine)

	s.WriteString("NO_OF_PARAMS=11")
	s.WriteString(newLine)

	s.WriteString("V1=")
	s.WriteString(strconv.FormatFloat(90+v1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("V3=")
	s.WriteString(strconv.FormatFloat(90-v3, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("E1=")
	s.WriteString(strconv.FormatFloat(e1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("A=")
	s.WriteString(strconv.FormatFloat(a, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("B=")
	s.WriteString(strconv.FormatFloat(b, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("C=")
	s.WriteString(strconv.FormatFloat(c, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("R1=")
	s.WriteString(strconv.FormatFloat(r1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMA=")
	s.WriteString(strconv.FormatFloat(gamma, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HO=")
	s.WriteString(strconv.FormatFloat(ho, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMAU=")
	s.WriteString(strconv.FormatFloat(gammau, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HU=")
	s.WriteString(strconv.FormatFloat(hu, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("END_OF_PARAMS")
	s.WriteString(newLine)

	return s.String()
}

func e23(v1, v2, v3, e1, b, gamma, ho, gammau, hu float64) string {
	s := strings.Builder{}

	if v1 > 45 {
		v1 = 45
	}
	if v1 < -45 {
		v1 = -45
	}
	if v3 > 45 {
		v3 = 45
	}
	if v3 < -45 {
		v3 = -45
	}

	s.WriteString("ENDCUT_TYPE=23")
	s.WriteString(newLine)

	s.WriteString("START_OF_PARAMS")
	s.WriteString(newLine)

	s.WriteString("NO_OF_PARAMS=9")
	s.WriteString(newLine)

	s.WriteString("V1=")
	s.WriteString(strconv.FormatFloat(90+v1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("V2=")
	s.WriteString(strconv.FormatFloat(v2, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("V3=")
	s.WriteString(strconv.FormatFloat(90-v3, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("E1=")
	s.WriteString(strconv.FormatFloat(e1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("B=")
	s.WriteString(strconv.FormatFloat(b, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMA=")
	s.WriteString(strconv.FormatFloat(gamma, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HO=")
	s.WriteString(strconv.FormatFloat(ho, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMAU=")
	s.WriteString(strconv.FormatFloat(gammau, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HU=")
	s.WriteString(strconv.FormatFloat(hu, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("END_OF_PARAMS")
	s.WriteString(newLine)

	return s.String()
}

func e23r(v1, v2, v3, r1, b, gamma, ho, gammau, hu float64) string {
	s := strings.Builder{}

	if v1 > 45 {
		v1 = 45
	}
	if v1 < -45 {
		v1 = -45
	}
	if v3 > 45 {
		v3 = 45
	}
	if v3 < -45 {
		v3 = -45
	}

	s.WriteString("ENDCUT_TYPE=23")
	s.WriteString(newLine)

	s.WriteString("START_OF_PARAMS")
	s.WriteString(newLine)

	s.WriteString("NO_OF_PARAMS=10")
	s.WriteString(newLine)

	s.WriteString("V1=")
	s.WriteString(strconv.FormatFloat(90+v1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("V2=")
	s.WriteString(strconv.FormatFloat(v2, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("V3=")
	s.WriteString(strconv.FormatFloat(90-v3, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("R1=")
	s.WriteString(strconv.FormatFloat(r1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("B=")
	s.WriteString(strconv.FormatFloat(b, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMA=")
	s.WriteString(strconv.FormatFloat(gamma, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HO=")
	s.WriteString(strconv.FormatFloat(ho, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMAU=")
	s.WriteString(strconv.FormatFloat(gammau, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HU=")
	s.WriteString(strconv.FormatFloat(hu, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("END_OF_PARAMS")
	s.WriteString(newLine)

	return s.String()
}

func e27(v1, v3, c, r2, a, e1, gamma, ho, gammau, hu float64) string {
	s := strings.Builder{}

	if v1 > 45 {
		v1 = 45
	}
	if v1 < -45 {
		v1 = -45
	}

	if v3 > 45 {
		v3 = 45
	}
	if v3 < 0 {
		v3 = 0
	}

	s.WriteString("ENDCUT_TYPE=27R")
	s.WriteString(newLine)

	s.WriteString("START_OF_PARAMS")
	s.WriteString(newLine)

	s.WriteString("NO_OF_PARAMS=10")
	s.WriteString(newLine)

	s.WriteString("V1=")
	s.WriteString(strconv.FormatFloat(90+v1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("V3=")
	s.WriteString(strconv.FormatFloat(90-v3, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("A=")
	s.WriteString(strconv.FormatFloat(a, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("C=")
	s.WriteString(strconv.FormatFloat(c, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("R2=")
	s.WriteString(strconv.FormatFloat(r2, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("E1=")
	s.WriteString(strconv.FormatFloat(e1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMA=")
	s.WriteString(strconv.FormatFloat(gamma, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HO=")
	s.WriteString(strconv.FormatFloat(ho, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMAU=")
	s.WriteString(strconv.FormatFloat(gammau, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HU=")
	s.WriteString(strconv.FormatFloat(hu, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("END_OF_PARAMS")
	s.WriteString(newLine)

	return s.String()
}

func e27r(v1, v3, r1, r2, a, e1, gamma, ho, gammau, hu float64) string {
	s := strings.Builder{}

	if v1 > 45 {
		v1 = 45
	}
	if v1 < -45 {
		v1 = -45
	}

	if v3 > 45 {
		v3 = 45
	}
	if v3 < 0 {
		v3 = 0
	}

	s.WriteString("ENDCUT_TYPE=27R")
	s.WriteString(newLine)

	s.WriteString("START_OF_PARAMS")
	s.WriteString(newLine)

	s.WriteString("NO_OF_PARAMS=10")
	s.WriteString(newLine)

	s.WriteString("V1=")
	s.WriteString(strconv.FormatFloat(90+v1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("V3=")
	s.WriteString(strconv.FormatFloat(90-v3, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("A=")
	s.WriteString(strconv.FormatFloat(a, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("R1=")
	s.WriteString(strconv.FormatFloat(r1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("R2=")
	s.WriteString(strconv.FormatFloat(r2, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("E1=")
	s.WriteString(strconv.FormatFloat(e1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMA=")
	s.WriteString(strconv.FormatFloat(gamma, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HO=")
	s.WriteString(strconv.FormatFloat(ho, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMAU=")
	s.WriteString(strconv.FormatFloat(gammau, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HU=")
	s.WriteString(strconv.FormatFloat(hu, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("END_OF_PARAMS")
	s.WriteString(newLine)

	return s.String()
}

func e28r(v1, v3, r1, r2, e1, a, b, gamma, ho, gammau, hu float64) string {
	s := strings.Builder{}

	if v1 > 45 {
		v1 = 45
	}
	if v1 < -45 {
		v1 = -45
	}

	if v3 > 45 {
		v3 = 45
	}
	if v3 < 0 {
		v3 = 0
	}

	s.WriteString("ENDCUT_TYPE=28R")
	s.WriteString(newLine)

	s.WriteString("START_OF_PARAMS")
	s.WriteString(newLine)

	s.WriteString("NO_OF_PARAMS=11")
	s.WriteString(newLine)

	s.WriteString("V1=")
	s.WriteString(strconv.FormatFloat(90+v1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("V3=")
	s.WriteString(strconv.FormatFloat(90-v3, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("R1=")
	s.WriteString(strconv.FormatFloat(r1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("R2=")
	s.WriteString(strconv.FormatFloat(r2, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("E1=")
	s.WriteString(strconv.FormatFloat(e1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("A=")
	s.WriteString(strconv.FormatFloat(a, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("B=")
	s.WriteString(strconv.FormatFloat(b, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMA=")
	s.WriteString(strconv.FormatFloat(gamma, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HO=")
	s.WriteString(strconv.FormatFloat(ho, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMAU=")
	s.WriteString(strconv.FormatFloat(gammau, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HU=")
	s.WriteString(strconv.FormatFloat(hu, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("END_OF_PARAMS")
	s.WriteString(newLine)

	return s.String()
}

func e29(v1, v3, r, r1, e1, vk, d, h, gamma, ho, gammau, hu float64) string {
	s := strings.Builder{}

	if v1 > 45 {
		v1 = 45
	}
	if v1 < -45 {
		v1 = -45
	}

	if v3 > 45 {
		v3 = 45
	}
	if v3 < 0 {
		v3 = 0
	}

	s.WriteString("ENDCUT_TYPE=29")
	s.WriteString(newLine)

	s.WriteString("START_OF_PARAMS")
	s.WriteString(newLine)

	s.WriteString("NO_OF_PARAMS=12")
	s.WriteString(newLine)

	s.WriteString("V1=")
	s.WriteString(strconv.FormatFloat(90+v1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("V3=")
	s.WriteString(strconv.FormatFloat(90-v3, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("R=")
	s.WriteString(strconv.FormatFloat(r, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("R1=")
	s.WriteString(strconv.FormatFloat(r1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("VK=")
	s.WriteString(strconv.FormatFloat(vk, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("E1=")
	s.WriteString(strconv.FormatFloat(d, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("D=")
	s.WriteString(strconv.FormatFloat(h, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("H=")
	s.WriteString(strconv.FormatFloat(e1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMA=")
	s.WriteString(strconv.FormatFloat(gamma, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HO=")
	s.WriteString(strconv.FormatFloat(ho, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMAU=")
	s.WriteString(strconv.FormatFloat(gammau, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HU=")
	s.WriteString(strconv.FormatFloat(hu, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("END_OF_PARAMS")
	s.WriteString(newLine)

	return s.String()
}

func e23f(a, b float64) string {
	s := strings.Builder{}

	s.WriteString("ENDCUT_TYPE=E23F")
	s.WriteString(newLine)

	s.WriteString("START_OF_PARAMS")
	s.WriteString(newLine)

	s.WriteString("NO_OF_PARAMS=2")
	s.WriteString(newLine)

	s.WriteString("A=")
	s.WriteString(strconv.FormatFloat(a, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("B=")
	s.WriteString(strconv.FormatFloat(b, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("END_OF_PARAMS")
	s.WriteString(newLine)

	return s.String()
}

func i1(x, r float64) string {
	s := strings.Builder{}

	s.WriteString("ICUT_TYPE=1")
	s.WriteString(newLine)

	s.WriteString("START_OF_PARAMS")
	s.WriteString(newLine)

	s.WriteString("NO_OF_PARAMS=3")
	s.WriteString(newLine)

	s.WriteString("X=")
	s.WriteString(strconv.FormatFloat(x, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("R=")
	s.WriteString(strconv.FormatFloat(r, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMA=0")
	s.WriteString(newLine)

	s.WriteString("END_OF_PARAMS")
	s.WriteString(newLine)

	return s.String()
}

func i2(x, y, r float64) string {
	s := strings.Builder{}

	s.WriteString("ICUT_TYPE=2")
	s.WriteString(newLine)

	s.WriteString("START_OF_PARAMS")
	s.WriteString(newLine)

	s.WriteString("NO_OF_PARAMS=4")
	s.WriteString(newLine)

	s.WriteString("X=")
	s.WriteString(strconv.FormatFloat(x, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("Y=")
	s.WriteString(strconv.FormatFloat(y, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("R=")
	s.WriteString(strconv.FormatFloat(r, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMA=0")
	s.WriteString(newLine)

	s.WriteString("END_OF_PARAMS")
	s.WriteString(newLine)

	return s.String()
}

func i4(x, a, b, r float64) string {
	s := strings.Builder{}

	s.WriteString("ICUT_TYPE=4")
	s.WriteString(newLine)

	s.WriteString("START_OF_PARAMS")
	s.WriteString(newLine)

	s.WriteString("NO_OF_PARAMS=5")
	s.WriteString(newLine)

	s.WriteString("X=")
	s.WriteString(strconv.FormatFloat(x-b/2, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("A=")
	s.WriteString(strconv.FormatFloat(a, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("B=")
	s.WriteString(strconv.FormatFloat(b, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("R=")
	s.WriteString(strconv.FormatFloat(r, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMA=0")
	s.WriteString(newLine)

	s.WriteString("END_OF_PARAMS")
	s.WriteString(newLine)

	return s.String()
}

func i11(x, a, b float64) string {
	s := strings.Builder{}

	s.WriteString("ICUT_TYPE=11")
	s.WriteString(newLine)

	s.WriteString("START_OF_PARAMS")
	s.WriteString(newLine)

	s.WriteString("NO_OF_PARAMS=6")
	s.WriteString(newLine)

	s.WriteString("X=")
	s.WriteString(strconv.FormatFloat(x, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("H=")
	s.WriteString(strconv.FormatFloat(b/2+20, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("V=0")
	s.WriteString(newLine)

	s.WriteString("A=")
	s.WriteString(strconv.FormatFloat(a, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("B=")
	s.WriteString(strconv.FormatFloat(b, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMA=0")
	s.WriteString(newLine)

	s.WriteString("END_OF_PARAMS")
	s.WriteString(newLine)

	return s.String()
}

func e210rs(v1, v2, v4, a, b, c, r1, cd, gammao, ho, gammau, hu float64) string {
	s := strings.Builder{}

	if v1 > 45 {
		v1 = 45
	}
	if v1 < -45 {
		v1 = -45
	}

	s.WriteString("ENDCUT_TYPE=E210RS")
	s.WriteString(newLine)

	s.WriteString("START_OF_PARAMS")
	s.WriteString(newLine)

	s.WriteString("NO_OF_PARAMS=12")
	s.WriteString(newLine)

	s.WriteString("V1=")
	s.WriteString(strconv.FormatFloat(90+v1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("V2=")
	s.WriteString(strconv.FormatFloat(v2, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("V4=")
	s.WriteString(strconv.FormatFloat(v4, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("A=")
	s.WriteString(strconv.FormatFloat(a, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("B=")
	s.WriteString(strconv.FormatFloat(b, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("C=")
	s.WriteString(strconv.FormatFloat(c, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("R1=")
	s.WriteString(strconv.FormatFloat(r1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("CD=")
	s.WriteString(strconv.FormatFloat(cd, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMAO=")
	s.WriteString(strconv.FormatFloat(gammao, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HO=")
	s.WriteString(strconv.FormatFloat(ho, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMAU=")
	s.WriteString(strconv.FormatFloat(gammau, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HU=")
	s.WriteString(strconv.FormatFloat(hu, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("END_OF_PARAMS")
	s.WriteString(newLine)

	return s.String()
}

func e213rs(v1, v2, v4, a, b, c, r1, cd, gammao, ho, gammau, hu float64) string {
	s := strings.Builder{}

	if v1 > 45 {
		v1 = 45
	}
	if v1 < -45 {
		v1 = -45
	}

	s.WriteString("ENDCUT_TYPE=E213RS")
	s.WriteString(newLine)

	s.WriteString("START_OF_PARAMS")
	s.WriteString(newLine)

	s.WriteString("NO_OF_PARAMS=12")
	s.WriteString(newLine)

	s.WriteString("V1=")
	s.WriteString(strconv.FormatFloat(90+v1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("V2=")
	s.WriteString(strconv.FormatFloat(v2, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("V4=")
	s.WriteString(strconv.FormatFloat(v4, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("A=")
	s.WriteString(strconv.FormatFloat(a, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("B=")
	s.WriteString(strconv.FormatFloat(b, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("C=")
	s.WriteString(strconv.FormatFloat(c, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("R1=")
	s.WriteString(strconv.FormatFloat(r1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("CD=")
	s.WriteString(strconv.FormatFloat(cd, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMAO=")
	s.WriteString(strconv.FormatFloat(gammao, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HO=")
	s.WriteString(strconv.FormatFloat(ho, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMAU=")
	s.WriteString(strconv.FormatFloat(gammau, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HU=")
	s.WriteString(strconv.FormatFloat(hu, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("END_OF_PARAMS")
	s.WriteString(newLine)

	return s.String()
}

func e221rs(v1, v3, a, b, c, r1, h, lf, gamma, ho, gammau, hu float64) string {
	s := strings.Builder{}

	if v1 > 45 {
		v1 = 45
	}
	if v1 < -45 {
		v1 = -45
	}

	if v3 > 45 {
		v3 = 45
	}
	if v3 < 0 {
		v3 = 0
	}

	s.WriteString("ENDCUT_TYPE=E221RS")
	s.WriteString(newLine)

	s.WriteString("START_OF_PARAMS")
	s.WriteString(newLine)

	s.WriteString("NO_OF_PARAMS=12")
	s.WriteString(newLine)

	s.WriteString("V1=")
	s.WriteString(strconv.FormatFloat(90+v1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("V3=")
	s.WriteString(strconv.FormatFloat(90-v3, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("A=")
	s.WriteString(strconv.FormatFloat(a, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("B=")
	s.WriteString(strconv.FormatFloat(b, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("C=")
	s.WriteString(strconv.FormatFloat(c, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("R1=")
	s.WriteString(strconv.FormatFloat(r1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("H=")
	s.WriteString(strconv.FormatFloat(h, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("LF=")
	s.WriteString(strconv.FormatFloat(lf, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMA=")
	s.WriteString(strconv.FormatFloat(gamma, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HO=")
	s.WriteString(strconv.FormatFloat(ho, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMAU=")
	s.WriteString(strconv.FormatFloat(gammau, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HU=")
	s.WriteString(strconv.FormatFloat(hu, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("END_OF_PARAMS")
	s.WriteString(newLine)

	return s.String()
}

func e372rs(v3, a, b, c, r1, c1, l, a1, orient, r2, dropBevel, gammao, ho, gammau, hu, ho2, gamma2 float64) string {
	s := strings.Builder{}

	if v3 > 45 {
		v3 = 45
	}
	if v3 < -45 {
		v3 = -45
	}

	s.WriteString("ENDCUT_TYPE=E372RS")
	s.WriteString(newLine)

	s.WriteString("START_OF_PARAMS")
	s.WriteString(newLine)

	s.WriteString("NO_OF_PARAMS=17")
	s.WriteString(newLine)

	s.WriteString("V3=")
	s.WriteString(strconv.FormatFloat(90-v3, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("A=")
	s.WriteString(strconv.FormatFloat(a, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("B=")
	s.WriteString(strconv.FormatFloat(b, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("C=")
	s.WriteString(strconv.FormatFloat(c, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("R1=")
	s.WriteString(strconv.FormatFloat(r1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("C1=")
	s.WriteString(strconv.FormatFloat(c1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("L=")
	s.WriteString(strconv.FormatFloat(l, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("A1=")
	s.WriteString(strconv.FormatFloat(a1, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("Orient=")
	s.WriteString(strconv.FormatFloat(orient, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("R2=")
	s.WriteString(strconv.FormatFloat(r2, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("Drop_Bevel=")
	s.WriteString(strconv.FormatFloat(dropBevel, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMAO=")
	s.WriteString(strconv.FormatFloat(gammao, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HO=")
	s.WriteString(strconv.FormatFloat(ho, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMAU=")
	s.WriteString(strconv.FormatFloat(gammau, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HU=")
	s.WriteString(strconv.FormatFloat(hu, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("HO2=")
	s.WriteString(strconv.FormatFloat(ho2, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("GAMMA2=")
	s.WriteString(strconv.FormatFloat(gamma2, 'f', -1, 64))
	s.WriteString(newLine)

	s.WriteString("END_OF_PARAMS")
	s.WriteString(newLine)

	return s.String()
}
