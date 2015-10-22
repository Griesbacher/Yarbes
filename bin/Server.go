package bin

import (
	"fmt"
	"github.com/griesbacher/SystemX/Config"
	"github.com/griesbacher/SystemX/Logging/LogServer"
	"github.com/griesbacher/SystemX/NetworkInterfaces/Incoming"
	"github.com/griesbacher/SystemX/RuleSystem"
	"github.com/griesbacher/SystemX/Tools/Strings"
	"log"
	"os"
	"os/signal"
	"reflect"
	"runtime/pprof"
	"syscall"
	"time"
)

//Server start a server config depending on the config file
func Server(serverConfigPath, clientConfigPath, cpuProfile string) {

	if cpuProfile != "" {
		f, err := os.Create(cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
	}
	Config.InitServerConfig(serverConfigPath)
	Config.InitClientConfig(clientConfigPath)

	stoppables := []Stoppable{}

	rpcInterfaces := []string{}

	if Config.GetServerConfig().LogServer.Enabled {
		logServer := LogServer.NewLogServer()
		logServer.Start()
		stoppables = append(stoppables, logServer)
		fmt.Println("Starting: LogServer")
		if Config.GetServerConfig().LogServer.RPCInterface != "" {
			fmt.Println("Starting: LogServer - RPC Interface")
			logServerRPCI := Incoming.NewLogServerRPCInterface(logServer.LogQueue)
			logServerRPCI.Start()
			stoppables = append(stoppables, logServerRPCI)
			rpcInterfaces = append(rpcInterfaces, Config.GetServerConfig().LogServer.RPCInterface)
		}
		time.Sleep(time.Duration(100) * time.Millisecond)
	}

	if Config.GetServerConfig().RuleSystem.Enabled {
		ruleSystem := RuleSystem.NewRuleSystem()
		ruleSystem.Start()
		stoppables = append(stoppables, ruleSystem)
		fmt.Println("Starting: RuleSystem")
		if Config.GetServerConfig().RuleSystem.RPCInterface != "" {
			fmt.Println("Starting: RuleSystem - RPC Interface")
			ruleSystemRPCI := Incoming.NewRuleSystemRPCInterface(ruleSystem)
			if !Strings.Contains(rpcInterfaces, Config.GetServerConfig().RuleSystem.RPCInterface) {
				fmt.Println("Starting: RPC")
				ruleSystemRPCI.Start()
				stoppables = append(stoppables, ruleSystemRPCI)
			}
		}
	}

	if Config.GetServerConfig().Proxy.Enabled {
		fmt.Println("Starting: Proxy - RPC Interface")
		proxyRPCI := Incoming.NewProxyRPCInterface()
		if !Strings.Contains(rpcInterfaces, Config.GetServerConfig().RuleSystem.RPCInterface) {
			fmt.Println("Starting: RPC")
			proxyRPCI.Start()
			stoppables = append(stoppables, proxyRPCI)
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
