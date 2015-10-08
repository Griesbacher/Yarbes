package main

import (
	"flag"
	"fmt"
	"github.com/griesbacher/SystemX/Config"
	"github.com/griesbacher/SystemX/RuleSystem"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/griesbacher/SystemX/NetworkInterfaces/Incoming"
)

type Stoppable interface {
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
	var rpcI *Incoming.RuleSystemRpcInterface



	if Config.GetServerConfig().LogServer.Enabled {
		if Config.GetServerConfig().LogServer.RpcInterface != "" {

		}
	}

	if Config.GetServerConfig().RuleSystem.Enabled {
		ruleSystem = RuleSystem.NewRuleSystem()
		ruleSystem.Start()
		if Config.GetServerConfig().RuleSystem.RpcInterface != "" {
			fmt.Println("Starting: RuleSystem RPC")
			rpcI = Incoming.NewRuleSystemRpcInterface(ruleSystem.EventQueue)
			if Config.GetServerConfig().LogServer.RpcInterface != Config.GetServerConfig().RuleSystem.RpcInterface {
				rpcI.Start()
			}
		}
	}

	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, syscall.SIGINT)
	signal.Notify(interruptChannel, syscall.SIGTERM)
	go func() {
		<-interruptChannel
		cleanUp([]Stoppable{ruleSystem, rpcI})
		os.Exit(1)
	}()

	//wait for the end to come
	for {
		time.Sleep(time.Duration(5) * time.Minute)
	}
}

func cleanUp(itemsToStop []Stoppable) {
	for _, item := range itemsToStop {
		item.Stop()
	}
}
