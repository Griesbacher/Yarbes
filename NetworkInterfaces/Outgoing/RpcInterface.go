package Outgoing

import (
	"crypto/tls"
	"errors"
	"github.com/griesbacher/SystemX/Config"
	"github.com/griesbacher/SystemX/Logging/LogServer"
	"github.com/griesbacher/SystemX/NetworkInterfaces"
	"github.com/griesbacher/SystemX/TLS"
	"net/rpc"
	"time"
)

//RPCInterface represents a outgoing RPC connection, with which a rpc.Client can be created
type RPCInterface struct {
	serverAddress string
	Config        *tls.Config
	conn          *tls.Conn
	client        *rpc.Client
}

//NewRPCInterface constructs a new RPCInterface
func NewRPCInterface(serverAddress string) *RPCInterface {
	config := TLS.GenerateClientTLSConfig(Config.GetClientConfig().TLS.Cert, Config.GetClientConfig().TLS.Key, Config.GetClientConfig().TLS.CaCert)
	return &RPCInterface{serverAddress: serverAddress, Config: config}
}

//Connect establishes a tcp connection and single byte for authentication and creates a rpc.Client
func (rpcI *RPCInterface) Connect() error {
	conn, err := tls.Dial("tcp", rpcI.serverAddress, rpcI.Config)
	if err != nil {
		return err
	}
	rpcI.conn = conn
	rpcI.conn.Write([]byte("a"))
	rpcI.client = rpc.NewClient(rpcI.conn)
	if rpcI.client == nil {
		return errors.New("Could not create rpc.Client")
	}
	return nil
}

//Disconnect closes the tcp connection
func (rpcI RPCInterface) Disconnect() {
	if rpcI.conn != nil {
		rpcI.conn.Close()
	}
}

//CreateEvent encapsulates the RPC call to create a Event on the server
func (rpcI RPCInterface) CreateEvent(event []byte) error {
	result := new(NetworkInterfaces.RPCResult)
	rpcEvent := NetworkInterfaces.RPCEvent{string(event), nil}
	if err := rpcI.client.Call("RuleSystemRPCHandler.CreateEvent", &rpcEvent, &result); err != nil {
		return err
	}
	return result.Err
}

//CreateDelayedEvent encapsulates the RPC call to create a DelayedEvent on the server
func (rpcI RPCInterface) CreateDelayedEvent(event []byte, delay *time.Duration) error {
	result := new(NetworkInterfaces.RPCResult)
	rpcEvent := NetworkInterfaces.RPCEvent{string(event), delay}
	if err := rpcI.client.Call("RuleSystemRPCHandler.CreateEvent", &rpcEvent, &result); err != nil {
		return err
	}
	return result.Err
}

//SendMessage sends a message to the logserver
func (rpcI RPCInterface) SendMessage(message *LogServer.LogMessage) error {
	result := new(NetworkInterfaces.RPCResult)
	if err := rpcI.client.Call("LogServerRPCHandler.SendMessage", message, &result); err != nil {
		return err
	}
	return result.Err
}

//SendMessages sends a multiple messages to the logserver
func (rpcI RPCInterface) SendMessages(messages *[]*LogServer.LogMessage) error {
	result := new(NetworkInterfaces.RPCResult)
	if err := rpcI.client.Call("LogServerRPCHandler.SendMessages", messages, &result); err != nil {
		return err
	}
	return result.Err
}
