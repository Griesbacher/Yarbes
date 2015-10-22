package bin

import (
	"flag"
	"fmt"
	"github.com/griesbacher/SystemX/Config"
	"github.com/griesbacher/SystemX/Logging/LogServer"
	"github.com/griesbacher/SystemX/NetworkInterfaces/Incoming"
	"github.com/griesbacher/SystemX/RuleSystem"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"
	"log"
	"runtime/pprof"
)

//Server start a server config depending on the config file
func Server() {
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
	flag.StringVar(&cpuProfile, "cpuprofile", "", "write cpu profile to file")
	flag.Parse()
	if cpuProfile != "" {
		f, err := os.Create(cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	Config.InitServerConfig(serverConfigPath)
	Config.InitClientConfig(clientConfigPath)

	var ruleSystem *RuleSystem.RuleSystem
	var ruleSystemRPCI *Incoming.RuleSystemRPCInterface
	var logServer *LogServer.Server
	var logServerRPCI *Incoming.LogServerRPCInterface

	stoppables := []Stoppable{}

	if Config.GetServerConfig().LogServer.Enabled {
		logServer = LogServer.NewLogServer()
		logServer.Start()
		stoppables = append(stoppables, logServer)
		fmt.Println("Starting: LogServer")
		if Config.GetServerConfig().LogServer.RPCInterface != "" {
			fmt.Println("Starting: LogServer - RPC Interface")
			logServerRPCI = Incoming.NewLogServerRPCInterface(logServer.LogQueue)
			logServerRPCI.Start()
			stoppables = append(stoppables, logServerRPCI)
		}
		time.Sleep(time.Duration(100) * time.Millisecond)
	}

	if Config.GetServerConfig().RuleSystem.Enabled {
		ruleSystem = RuleSystem.NewRuleSystem()
		ruleSystem.Start()
		stoppables = append(stoppables, ruleSystem)
		fmt.Println("Starting: RuleSystem")
		if Config.GetServerConfig().RuleSystem.RPCInterface != "" {
			fmt.Println("Starting: RuleSystem - RPC Interface")
			ruleSystemRPCI = Incoming.NewRuleSystemRPCInterface(ruleSystem)
			if Config.GetServerConfig().LogServer.RPCInterface != Config.GetServerConfig().RuleSystem.RPCInterface || (logServerRPCI == nil || !logServerRPCI.IsRunning()) {
				fmt.Println("Starting: RPC")
				ruleSystemRPCI.Start()
				stoppables = append(stoppables, ruleSystemRPCI)
			}
		}
	}

	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, syscall.SIGINT)
	signal.Notify(interruptChannel, syscall.SIGTERM)
	quit := make(chan bool)
	go func() {
		<-interruptChannel
		cleanUp(stoppables)
		quit <- true
	}()
	fmt.Println("Everything's ready!")
	//wait for the end to come
	<-quit
	pprof.StopCPUProfile()
	fmt.Println("Bye")
}

func cleanUp(itemsToStop []Stoppable) {
	for _, item := range itemsToStop {
		if item != nil && item.IsRunning() {
			fmt.Println("Stopping: ", reflect.TypeOf(item))
			item.Stop()
		}
	}
}
