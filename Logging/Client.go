package Logging

import (
	"errors"
	"fmt"
	"github.com/griesbacher/SystemX/Logging/Local"
	"github.com/griesbacher/SystemX/Logging/LogServer"
	"github.com/griesbacher/SystemX/NetworkInterfaces"
	"github.com/griesbacher/SystemX/NetworkInterfaces/Outgoing"
	"github.com/kdar/factorlog"
	"net/rpc"
	"os"
)

//Client combines locallogging with factorlog and remote logging via RPC
type Client struct {
	logRPC      *Outgoing.RPCInterface
	rpcClient   *rpc.Client
	name        string
	localLogger *factorlog.FactorLog
}

//NewLocalClient constructs a new client, which logs to stdout
func NewLocalClient() *Client {
	return &Client{localLogger: Local.GetLogger()}
}

//NewClient creates a localClient or a RPC if a address is given
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
		return &Client{logRPC: &logRPC, rpcClient: rpcClient, name: clientName, localLogger: Local.GetLogger()}, nil

	}
	return nil, errors.New("Could not create a RPC client")
}

//LogMultiple sends the logMessages to the remote logServer, log an error to stdout
func (client Client) LogMultiple(messages *[]*LogServer.LogMessage) {
	result := new(NetworkInterfaces.RPCResult)
	if err := client.rpcClient.Call("LogServerRPCHandler.SendMessages", messages, &result); err != nil {
		client.localLogger.Error(err)
		for message := range *messages{
			client.localLogger.Debug("Message", message)
		}
	}

	if result.Err != nil {
		client.localLogger.Error(result.Err)
	}
}

//Log sends the logMessage to the remote logServer, log an error to stdout
func (client Client) Log(message *LogServer.LogMessage) {
	result := new(NetworkInterfaces.RPCResult)
	if err := client.rpcClient.Call("LogServerRPCHandler.SendMessage", message, &result); err != nil {
		client.localLogger.Error(err)
		client.localLogger.Debug("Message", message)
	}

	if result.Err != nil {
		client.localLogger.Error(result.Err)
	}
}

//Disconnect closes the connection to the remote logServer
func (client Client) Disconnect() {
	if client.logRPC != nil{
		client.logRPC.Disconnect()
	}
}

//Debug logs the message local/remote to on debug level
func (client Client) Debug(v ...interface{}) {
	if client.rpcClient == nil {
		client.localLogger.Debug(v)
	} else {
		client.Log(LogServer.NewDebugLogMessage(client.name, fmt.Sprint(v)))
	}
}

//Info logs the message local/remote to on info level
func (client Client) Info(v ...interface{}) {
	if client.rpcClient == nil {
		client.localLogger.Info(v)
	} else {
		client.Log(LogServer.NewInfoLogMessage(client.name, fmt.Sprint(v)))
	}
}

//Warn logs the message local/remote to on warn level
func (client Client) Warn(v ...interface{}) {
	if client.rpcClient == nil {
		client.localLogger.Warn(v)
	} else {
		client.Log(LogServer.NewWarnLogMessage(client.name, fmt.Sprint(v)))
	}
}

//Error logs the message local/remote to on error level
func (client Client) Error(v ...interface{}) {
	if client.rpcClient == nil {
		client.localLogger.Error(v)
	} else {
		client.Log(LogServer.NewErrorLogMessage(client.name, fmt.Sprint(v)))
	}
}
