package main

import (
	"flag"
	"fmt"
	"github.com/griesbacher/SystemX/Config"
	"github.com/griesbacher/SystemX/NetworkInterfaces"
	"github.com/griesbacher/SystemX/NetworkInterfaces/Outgoing"
	"github.com/griesbacher/SystemX/LogServer"
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

		eventRpc := Outgoing.NewRpcInterface(Config.GetClientConfig().Backend.RpcInterface)
		err := eventRpc.Connect()
		if err != nil {
			panic(err)
		}
		if rpcClient := eventRpc.GenRpcClient(); rpcClient != nil {
			result := new(NetworkInterfaces.RpcResult)
			if err := rpcClient.Call("RuleSystemRpcHandler.CreateEvent", string(b), &result); err != nil {
				panic(err)
			}
			if result.Err != nil {
				panic(result.Err)
			}
		}
		eventRpc.Disconnect()

	logRpc := Outgoing.NewRpcInterface(Config.GetClientConfig().LogServer.RpcInterface)
	lerr := logRpc.Connect()
	if lerr != nil {
		panic(lerr)
	}
	if rpcClient := logRpc.GenRpcClient(); rpcClient != nil {
		result := new(NetworkInterfaces.RpcResult)
		for i := 0; i<10; i++ {
			start := time.Now()
			message := LogServer.NewLogMessage("client0", "Hallo Log ")
			if err := rpcClient.Call("LogServerRpcHandler.SendMessage", &message, &result); err != nil {
				panic(err)
			}
			if result.Err != nil {
				panic(result.Err)
			}
			fmt.Println(time.Now().Sub(start))
		}
	}
	logRpc.Disconnect()
}
