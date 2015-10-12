package LogServer

import (
	"errors"
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

func NewClient() (*Client, error) {
	logRPC := Outgoing.NewRPCInterface(Config.GetClientConfig().LogServer.RPCInterface)
	if err := logRPC.Connect(); err != nil {
		return nil, err
	}
	if rpcClient := logRPC.GenRPCClient(); rpcClient != nil {
		var clientName string
		for name := range logRPC.Config.NameToCertificate {
			clientName = name
			break
		}
		if clientName == "" {
			var err error
			clientName, err = os.Hostname()
			if err != nil {
				return nil, err
			}
		}
		return &Client{logRPC: logRPC, rpcClient: rpcClient, name: clientName}, nil

	}
	return nil, errors.New("Could not create a RPC client")
}

func (client Client) Log(message *LogMessage) error {
	result := new(NetworkInterfaces.RPCResult)
	if err := client.rpcClient.Call("LogServerRPCHandler.SendMessage", message, &result); err != nil {
		return err
	}
	return result.Err
}

func (client Client) Disconnect() {
	client.logRPC.Disconnect()
}

func (client Client) Debug(message string) error {
	logMessage := NewDebugLogMessage(client.name, message)
	return client.Log(logMessage)
}
