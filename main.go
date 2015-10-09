package main

import (
	"flag"
	"fmt"
	"github.com/griesbacher/SystemX/Config"
	"github.com/griesbacher/SystemX/LogServer"
	"github.com/griesbacher/SystemX/NetworkInterfaces/Incoming"
	"github.com/griesbacher/SystemX/RuleSystem"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"
)

type stoppable interface {
	Stop()
}

func main() {
	var configPath string
	flag.Usage = func() {
		fmt.Println(`SystemX by Philip Griesbacher @ 2015
Commandline Parameter:
-configPath Path to the config file. If no file path is given the default is ./serverConfig.gcfg.
		`)
	}
	flag.StringVar(&configPath, "configPath", "serverConfig.gcfg", "path to the config file")
	flag.Parse()
	Config.InitServerConfigProvider(configPath)

	var ruleSystem *RuleSystem.RuleSystem
	var ruleSystemRPCI *Incoming.RuleSystemRPCInterface
	var logServer *LogServer.Server
	var logServerRPCI *Incoming.LogServerRPCInterface

	if Config.GetServerConfig().LogServer.Enabled {
		logServer = LogServer.NewLogServer()
		logServer.Start()
		if Config.GetServerConfig().LogServer.RPCInterface != "" {
			fmt.Println("Starting: LogServer RPC")
			logServerRPCI = Incoming.NewLogServerRPCInterface(logServer.LogQueue)
			logServerRPCI.Start()
		}
	}

	if Config.GetServerConfig().RuleSystem.Enabled {
		ruleSystem = RuleSystem.NewRuleSystem()
		ruleSystem.Start()
		if Config.GetServerConfig().RuleSystem.RPCInterface != "" {
			fmt.Println("Starting: RuleSystem RPC")
			ruleSystemRPCI = Incoming.NewRuleSystemRPCInterface(ruleSystem.EventQueue)
			if Config.GetServerConfig().LogServer.RPCInterface != Config.GetServerConfig().RuleSystem.RPCInterface {
				ruleSystemRPCI.Start()
			}
		}
	}

	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, syscall.SIGINT)
	signal.Notify(interruptChannel, syscall.SIGTERM)
	go func() {
		<-interruptChannel
		cleanUp([]stoppable{logServerRPCI, logServer, ruleSystemRPCI, ruleSystem})
		os.Exit(1)
	}()
	fmt.Println("Everything's ready!")
	//wait for the end to come
	for {
		time.Sleep(time.Duration(5) * time.Second)
	}
}

func cleanUp(itemsToStop []stoppable) {
	for _, item := range itemsToStop {
		fmt.Println("Stopping: ", reflect.TypeOf(item))
		item.Stop()
	}
}
