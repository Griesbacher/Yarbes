package main

import (
	"github.com/griesbacher/SystemX/bin"
	"time"
	"flag"
	"os"
	"log"
"runtime/pprof"
)

var cpuProfile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	go bin.Server()
	time.Sleep(time.Duration(1) * time.Second)
	bin.Client()
}
