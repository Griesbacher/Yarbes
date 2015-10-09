package LogServer

import (
	"github.com/griesbacher/SystemX/Config"
	"github.com/griesbacher/SystemX/NetworkInterfaces"
	"github.com/griesbacher/SystemX/NetworkInterfaces/Outgoing"
	"net/rpc"
	"os"
)

type Client struct {
	logRPC    Outgoing.RPCInterface
	rpcClient *rpc.Client
	name      string
}

func NewClient() *Client {
	logRPC := Outgoing.NewRPCInterface(Config.GetClientConfig().LogServer.RPCInterface)
	if err := logRPC.Connect(); err != nil {
		panic(err)
	}
	if rpcClient := logRPC.GenRPCClient(); rpcClient != nil {
		var clientName string
		for name, _ := range logRPC.Config.NameToCertificate {
			clientName = name
			break
		}
		if clientName == "" {
			var err error
			clientName, err = os.Hostname()
			if err != nil {
				panic(err)
			}
		}
		return &Client{logRPC: logRPC, rpcClient: rpcClient, name: clientName}

	}
	return &Client{}
}

func (client Client) Log(message string) {
	result := new(NetworkInterfaces.RPCResult)
	logMessage := NewLogMessage(client.name, message)
	if err := client.rpcClient.Call("LogServerRPCHandler.SendMessage", &logMessage, &result); err != nil {
		panic(err)
	}
	if result.Err != nil {
		panic(result.Err)
	}
}

func (client Client) Disconnect() {
	client.logRPC.Disconnect()
}
