package main

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/edanko/moses/cmd/conv/cmd"
)

func main() {

	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close() // error handling omitted for example
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	// license check
	//if ok, err := license.IsLicensed(); !ok {
	//	log.Fatalln(err)
	//}
	cmd.Execute()

	f, err = os.Create("mem.prof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f.Close() // error handling omitted for example
	runtime.GC()    // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
}
