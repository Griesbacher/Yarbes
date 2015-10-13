package Logging


import (
	"errors"
	"github.com/griesbacher/SystemX/NetworkInterfaces"
	"github.com/griesbacher/SystemX/NetworkInterfaces/Outgoing"
	"net/rpc"
	"os"
	"fmt"
	"github.com/griesbacher/SystemX/Logging/Local"
	"github.com/griesbacher/SystemX/Logging/LogServer"
	"github.com/kdar/factorlog"
)

type Client struct {
	logRPC      Outgoing.RPCInterface
	rpcClient   *rpc.Client
	name        string
	localLogger *factorlog.FactorLog
}

func NewLocalClient() *Client {
	return &Client{localLogger:Local.GetLogger()}
}

func NewClient(target string) (*Client, error) {
	if target == "" {
		//use local logger
		return NewLocalClient(), nil
	}
	logRPC := Outgoing.NewRPCInterface(target)
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
		return &Client{logRPC: logRPC, rpcClient: rpcClient, name: clientName, localLogger:Local.GetLogger()}, nil

	}
	return nil, errors.New("Could not create a RPC client")
}

func (client Client) Log(message *LogServer.LogMessage) {
	result := new(NetworkInterfaces.RPCResult)
	if err := client.rpcClient.Call("LogServerRPCHandler.SendMessage", message, &result); err != nil {
		client.localLogger.Error(err)
	}

	if result.Err != nil{
		client.localLogger.Error(result.Err)
	}
}

func (client Client) Disconnect() {
	client.logRPC.Disconnect()
}

func (client Client) Debug(v... interface{}) {
	if client.rpcClient == nil {
		client.localLogger.Debug(v)
	}else {
		client.Log(LogServer.NewDebugLogMessage(client.name, fmt.Sprint(v)))
	}
}

func (client Client) Warn(v... interface{}) {
	if client.rpcClient == nil {
		client.localLogger.Warn(v)
	}else {
		client.Log(LogServer.NewWarnLogMessage(client.name, fmt.Sprint(v)))
	}
}

func (client Client) Error(v... interface{}) {
	if client.rpcClient == nil {
		client.localLogger.Error(v)
	}else {
		client.Log(LogServer.NewErrorLogMessage(client.name, fmt.Sprint(v)))
	}
}
