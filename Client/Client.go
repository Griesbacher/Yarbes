package main

import (
	"flag"
	"fmt"
	"github.com/griesbacher/SystemX/Config"
	"github.com/griesbacher/SystemX/LogServer"
	"github.com/griesbacher/SystemX/NetworkInterfaces"
	"github.com/griesbacher/SystemX/NetworkInterfaces/Outgoing"
	"time"
)

func main() {
	var configPath string
	flag.Usage = func() {
		fmt.Println(`SystemX by Philip Griesbacher @ 2015
Commandline Parameter:
-configPath Path to the config file. If no file path is given the default is ./serverConfig.gcfg.
		`)
	}
	flag.StringVar(&configPath, "configPath", "clientConfig.gcfg", "path to the config file")
	flag.Parse()
	Config.InitClientonfigProvider(configPath)

	b := []byte(`{
		   "k1" : "v1",
		   "k2" : 10,
		   "k3" : ["v4",12.3,{"k11" : "v11", "k22" : "v22"}]
		}`)

	eventRPC := Outgoing.NewRPCInterface(Config.GetClientConfig().Backend.RPCInterface)
	err := eventRPC.Connect()
	if err != nil {
		panic(err)
	}
	if rpcClient := eventRPC.GenRPCClient(); rpcClient != nil {
		result := new(NetworkInterfaces.RPCResult)
		if err := rpcClient.Call("RuleSystemRPCHandler.CreateEvent", string(b), &result); err != nil {
			panic(err)
		}
		if result.Err != nil {
			panic(result.Err)
		}
	}
	eventRPC.Disconnect()

	logRPC := Outgoing.NewRPCInterface(Config.GetClientConfig().LogServer.RPCInterface)
	lerr := logRPC.Connect()
	if lerr != nil {
		panic(lerr)
	}
	if rpcClient := logRPC.GenRPCClient(); rpcClient != nil {
		result := new(NetworkInterfaces.RPCResult)
		for i := 0; i < 10; i++ {
			start := time.Now()
			message := LogServer.NewLogMessage("client0", "Hallo Log ")
			if err := rpcClient.Call("LogServerRPCHandler.SendMessage", &message, &result); err != nil {
				panic(err)
			}
			if result.Err != nil {
				panic(result.Err)
			}
			fmt.Println(time.Now().Sub(start))
		}
	}
	logRPC.Disconnect()
}
