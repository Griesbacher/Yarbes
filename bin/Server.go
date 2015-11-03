package bin

import (
	"encoding/gob"
	"fmt"
	"github.com/griesbacher/Yarbes/Config"
	"github.com/griesbacher/Yarbes/Logging/LogServer"
	httpIncoming "github.com/griesbacher/Yarbes/NetworkInterfaces/HTTP/Incoming"
	rpcIncoming "github.com/griesbacher/Yarbes/NetworkInterfaces/RPC/Incoming"
	"github.com/griesbacher/Yarbes/RuleSystem"
	"github.com/griesbacher/Yarbes/Tools/Strings"
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
	gob.Register(map[string]interface{}{})
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
	httpInterfaces := []string{}

	if Config.GetServerConfig().LogServer.Enabled {
		logServer := LogServer.NewLogServer()
		logServer.Start()
		stoppables = append(stoppables, logServer)
		fmt.Println("Starting: LogServer")
		if Config.GetServerConfig().LogServer.RPCInterface != "" {
			fmt.Println("Starting: LogServer - RPC Interface")
			logServerRPCI := rpcIncoming.NewLogServerRPCInterface(logServer.LogQueue)
			logServerRPCI.Start()
			stoppables = append(stoppables, logServerRPCI)
			rpcInterfaces = append(rpcInterfaces, Config.GetServerConfig().LogServer.RPCInterface)
		}
		if Config.GetServerConfig().LogServer.HTTPInterface != "" {
			fmt.Println("Starting: LogServer - HTTP Interface")
			logServerRPCI := httpIncoming.NewLogServerHTTPInterface(logServer.InfluxClient)
			logServerRPCI.Start()
			stoppables = append(stoppables, logServerRPCI)
			httpInterfaces = append(httpInterfaces, Config.GetServerConfig().LogServer.HTTPInterface)
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
			ruleSystemRPCI := rpcIncoming.NewRuleSystemRPCInterface(ruleSystem)
			if !Strings.Contains(rpcInterfaces, Config.GetServerConfig().RuleSystem.RPCInterface) {
				fmt.Println("Starting: RPC")
				ruleSystemRPCI.Start()
				stoppables = append(stoppables, ruleSystemRPCI)
			}
		}
	}

	if Config.GetServerConfig().Proxy.Enabled {
		fmt.Println("Starting: Proxy - RPC Interface")
		proxyRPCI := rpcIncoming.NewProxyRPCInterface()
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
	if cpuProfile != "" {
		pprof.StopCPUProfile()
	}
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
