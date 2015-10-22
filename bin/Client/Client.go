package main

import (
	"flag"
	"fmt"
	"github.com/griesbacher/SystemX/bin"
)

func main() {
	var configPath string
	var cpuProfile string
	flag.Usage = func() {
		fmt.Println(`SystemX by Philip Griesbacher @ 2015
Commandline Parameter:
-configPath Path to the config file. If no file path is given the default is ./serverConfig.gcfg.
		`)
	}
	flag.StringVar(&configPath, "configPath", "clientConfig.gcfg", "path to the config file")
	flag.StringVar(&cpuProfile, "pprof", "", "write cpu profile to given file")
	flag.Parse()

	bin.Client(configPath, cpuProfile)
}
