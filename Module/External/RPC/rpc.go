package main

import (
	"os"
	"fmt"
	"github.com/griesbacher/SystemX/NetworkInterfaces/Outgoing"
	"strings"
	"encoding/json"
	"errors"
	"github.com/griesbacher/SystemX/Config"
	"encoding/gob"
)

func main() {
	gob.Register(map[string]interface{}{})
	if len(os.Args) < 2 {
		os.Exit(1)
	}
	jsonString := os.Args[1]
	args := strings.Split(os.Args[2], ",")

	serverAddress := fmt.Sprintf("%s:%s", args[0], args[1])
	Config.InitClientConfig("clientConfig.gcfg")
	rpcClient := Outgoing.NewRPCInterface(serverAddress)
	if rpcClient == nil {
		panic(errors.New("Can not create RPC Client"))
	}
	rpcClient.Connect()

	err, result := rpcClient.MakeCall(args[2], []byte(jsonString))
	if err != nil {
		panic(err)
	}
	rpcClient.Disconnect()
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", jsonBytes)
}