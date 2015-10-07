package main

import (
	"flag"
	"fmt"
	"github.com/griesbacher/SystemX/Client/Interface"
	"github.com/griesbacher/SystemX/Config"
	"github.com/griesbacher/SystemX/RuleSystem/Incoming"
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

	rpcI := Interface.NewRpcInterface(Config.GetClientConfig().Client.RpcInterface)
	rpcI.Connect()

	if rpcClient := rpcI.GenRpcClient(); rpcClient != nil {
		result := new(Incoming.Result)
		if err := rpcClient.Call("RpcHandler.CreateEvent", string(b), &result); err != nil {
			panic(err)
		}
		if result.Err != nil {
			panic(result.Err)
		}
	}
	rpcI.Disconnect()
}
