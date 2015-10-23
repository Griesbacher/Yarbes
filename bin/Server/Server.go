package main

import (
	"flag"
	"fmt"
	"github.com/griesbacher/SystemX/bin"
)

func main() {
	var serverConfigPath string
	var clientConfigPath string
	var cpuProfile string
	flag.Usage = func() {
		fmt.Println(`SystemX by Philip Griesbacher @ 2015
Commandline Parameter:
-serverConfigPath Path to the server config file. If no file path is given the default is ./serverConfig.gcfg.
-clientConfigPath Path to the client config file. If no file path is given the default is ./clientConfig.gcfg.
		`)
	}
	flag.StringVar(&serverConfigPath, "serverConfigPath", "serverConfig.gcfg", "path to the server config file")
	flag.StringVar(&clientConfigPath, "clientConfigPath", "clientConfig.gcfg", "path to the client config file")
	flag.StringVar(&cpuProfile, "pprof", "", "write cpu profile to given file")
	flag.Parse()

	bin.Server(serverConfigPath, clientConfigPath, cpuProfile)
}
