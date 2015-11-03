package main

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/griesbacher/Yarbes/Config"
	"github.com/griesbacher/Yarbes/Module"
	"github.com/griesbacher/Yarbes/NetworkInterfaces/RPC/Outgoing"
	"os"
)

func main() {
	gob.Register(map[string]interface{}{})
	if len(os.Args) < 2 {
		os.Exit(1)
	}
	jsonString := os.Args[4]

	serverAddress := fmt.Sprintf("%s:%s", os.Args[1], os.Args[2])
	Config.InitClientConfig("clientConfig.gcfg")
	rpcClient := Outgoing.NewRPCInterface(serverAddress)
	if rpcClient == nil {
		panic(errors.New("Can not create RPC Client"))
	}
	rpcClient.Connect()

	result, err := rpcClient.MakeCall(os.Args[3], []byte(jsonString))
	if err != nil {
		panic(err)
	}
	rpcClient.Disconnect()
	result.Messages = append(result.Messages, Module.Message{Severity: "debug", Message: fmt.Sprintf("This event was sent over rpc: %s", result.Event), Source: "RPC Module"})
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", jsonBytes)
}
