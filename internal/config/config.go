package config

import (
	"archive/zip"
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Init loads all yaml configs in current directory
func Init() error {
	if fileExists("config.zip") {
		r, err := zip.OpenReader("config.zip")
		if err != nil {
			log.Fatal(err)
		}
		defer r.Close()

		viper.SetConfigType("yaml")

		for _, f := range r.File {
			rc, err := f.Open()
			if err != nil {
				log.Fatal(err)
			}
			defer rc.Close()

			rd := bufio.NewReader(rc)

			err = viper.MergeConfig(rd)
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Println("[*] config read from config.zip")
	} else {

		var configFiles []string
		err := filepath.Walk(".", func(path string, info os.FileInfo, e error) error {
			if e != nil {
				return e
			}
			if info.Mode().IsRegular() {
				if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml" {
					configFiles = append(configFiles, path)
				}
			}
			return nil
		})
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if len(configFiles) == 0 {
			fmt.Println("No config files found!")
			os.Exit(0)
		}
		// TODO: do something with this
		viper.SetConfigFile("../../../config/config.yaml")
		if err = viper.MergeInConfig(); err != nil {

			viper.SetConfigFile("../../config/config.yaml")
			if err = viper.MergeInConfig(); err != nil {
				return err
			} else {
				goto cont
			}
		}
	cont:
		for _, f := range configFiles {
			viper.SetConfigFile(f)
			err = viper.MergeInConfig()
		}
		fmt.Println("[*] config read from local dir")
		return err
	}
	return nil
}

func Spacing(endcutName, profType string) float64 {

	profType = strings.ToLower(profType)

	onlyType := profType[:2]

	switch onlyType {
	case "rp":
		onlyType = "hp"
		//case "fb":
		//onlyType = "fb"
		//case "ae":
		//	onlyType = "ae"
		//default:
		//	onlyType = "hp"
	}

	key := "spacing." + onlyType + "." + endcutName + "." + profType

	res := viper.GetFloat64(key)
	if res == 0 {
		res = viper.GetFloat64("spacing." + onlyType + ".default")
		fmt.Println(" --> add", profType, "profile params to", endcutName, "endcut params in spacing config")
	}
	return res
}

func Norm(profType string) string {
	switch strings.ToLower(profType[:2]) {
	case "fb":
		return flatBarNorm(profType)
	case "rp":
		return rpNorm(profType)
	case "hp":
		return rpNorm("rp" + profType[2:])
	case "ae":
		return aeNorm(profType)
	default:
		return viper.GetString("profiles." + profType)
	}
}

func aeNorm(profType string) string {
	var a, b, r1, r2, s, t float64

	switch profType {
	case "ae250x90x10x15":
		a = 250
		b = 90
		r1 = 19
		r2 = 9.5
		s = 10
		t = 15
	default:
		fmt.Println("[x] add case for", profType, "to aeNorm")
		os.Exit(1)
	}
	return fmt.Sprintf("SHAPE=L\r\nSTART_OF_PARAMS\r\nNO_OF_PARAMS=6\r\nNORM=%s\r\nA=%0.1f\r\nB=%0.1f\r\nR1=%0.1f\r\nR2=%0.1f\r\nS=%0.1f\r\nT=%0.1f\r\nEND_OF_PARAMS", strings.ToUpper("L"+profType[2:]), a, b, r1, r2, s, t)
}

func rpNorm(profType string) string {
	var b, c, s, r float64

	switch profType {
	case "rp100x6":
		b = 100
		c = 20
		s = 6
		r = 5
	case "rp120x6.5":
		b = 120
		c = 23.5
		s = 6.5
		r = 5
	case "rp140x7":
		b = 140
		c = 26
		s = 7
		r = 6
	case "rp140x9":
		b = 140
		c = 26
		s = 9
		r = 6
	case "rp160x8":
		b = 160
		c = 28
		s = 8
		r = 7
	case "rp160x10":
		b = 160
		c = 28
		s = 10
		r = 7
	case "rp180x9":
		b = 180
		c = 31
		s = 9
		r = 7
	case "rp180x11":
		b = 180
		c = 31
		s = 11
		r = 7
	case "rp200x10":
		b = 200
		c = 34
		s = 10
		r = 8
	case "rp200x11":
		b = 200
		c = 34
		s = 11
		r = 8
	case "rp200x12":
		b = 200
		c = 34
		s = 12
		r = 8
	case "rp220x11":
		b = 220
		c = 37
		s = 11
		r = 8.5
	case "rp220x13":
		b = 220
		c = 37
		s = 13
		r = 8.5
	case "rp240x12":
		b = 240
		c = 40
		s = 12
		r = 9
	case "rp240x14":
		b = 240
		c = 40
		s = 14
		r = 9
	default:
		fmt.Println("[x] add case for", profType, "to rpNorm")
		os.Exit(1)
	}
	return fmt.Sprintf("SHAPE=HP\r\nSTART_OF_PARAMS\r\nNO_OF_PARAMS=5\r\nNORM=%s\r\nB=%0.1f\r\nC=%0.1f\r\nS=%0.1f\r\nR=%0.1f\r\nEND_OF_PARAMS", strings.ToUpper(profType), b, c, s, r)
}

func flatBarNorm(profType string) string {
	t := strings.Split(profType[2:], "x")
	return fmt.Sprintf("SHAPE=FB\r\nSTART_OF_PARAMS\r\nNO_OF_PARAMS=3\r\nNORM=%s\r\nB=%s\r\nS=%s\r\nEND_OF_PARAMS", strings.ToUpper(profType), t[0], t[1])
}

func BarSize(profType string) float64 {
	key := "barsize." + strings.ToLower(profType)
	res := viper.GetFloat64(key)

	if res == 0 {
		switch profType[:2] {
		case "ae":
			key = "barsize.ae"
		default:
			key = "barsize.default"
		}
		res = viper.GetFloat64(key)
	}

	return res
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
