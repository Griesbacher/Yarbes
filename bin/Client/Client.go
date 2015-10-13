package main

import (
	"flag"
	"fmt"
	"github.com/griesbacher/SystemX/Config"
	"github.com/griesbacher/SystemX/NetworkInterfaces"
	"github.com/griesbacher/SystemX/NetworkInterfaces/Outgoing"
	"github.com/griesbacher/SystemX/Logging"
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
	Config.InitClientConfigProvider(configPath)

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

	client, err := Logging.NewClient(Config.GetClientConfig().LogServer.RPCInterface)
	if err != nil {
		panic(err)
	}
	client.Debug("Hallo Server")
	client.Disconnect()

}
