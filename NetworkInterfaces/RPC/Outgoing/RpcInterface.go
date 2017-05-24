package Outgoing

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/griesbacher/Yarbes/Config"
	"github.com/griesbacher/Yarbes/Logging/LogServer"
	"github.com/griesbacher/Yarbes/Module"
	"github.com/griesbacher/Yarbes/NetworkInterfaces/RPC"
	"github.com/griesbacher/Yarbes/TLS"
	"net"
	"net/rpc"
	"reflect"
	"time"
)

//RPCInterface represents a outgoing RPC connection, with which a rpc.Client can be created
type RPCInterface struct {
	serverAddress string
	Config        *tls.Config
	conn          interface{}
	client        *rpc.Client
}

//NewRPCInterface constructs a new RPCInterface
func NewRPCInterface(serverAddress string) *RPCInterface {
	var config *tls.Config
	if Config.GetClientConfig().TLS.Enable {
		config = TLS.GenerateClientTLSConfig(Config.GetClientConfig().TLS.Cert, Config.GetClientConfig().TLS.Key, Config.GetClientConfig().TLS.CaCert)
	}
	return &RPCInterface{serverAddress: serverAddress, Config: config}
}

//Connect establishes a tcp connection and single byte for authentication and creates a rpc.Client
func (rpcI *RPCInterface) Connect() error {
	if Config.GetClientConfig().TLS.Enable {
		conn, err := tls.Dial("tcp", rpcI.serverAddress, rpcI.Config)
		if err != nil {
			return err
		}
		conn.Write([]byte("a"))
		rpcI.client = rpc.NewClient(conn)
		rpcI.conn = conn
	} else {
		conn, err := net.Dial("tcp", rpcI.serverAddress)
		if err != nil {
			return err
		}
		conn.Write([]byte("a"))
		rpcI.client = rpc.NewClient(conn)
		rpcI.conn = conn
	}
	if rpcI.client == nil {
		return errors.New("Could not create rpc.Client")
	}
	return nil
}

//Disconnect closes the tcp connection
func (rpcI RPCInterface) Disconnect() {
	if rpcI.conn != nil {
		switch conn := rpcI.conn.(type) {
		case *tls.Conn:
			(*conn).Close()
		case *net.Conn:
			(*conn).Close()
		case *net.TCPConn:
			(*conn).Close()
		default:
			fmt.Println("Diconnect error")
			fmt.Println(reflect.TypeOf(conn))
		}
	}
}

//CreateEvent encapsulates the RPC call to create a Event on the server
func (rpcI RPCInterface) CreateEvent(event []byte) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = rec.(error)
		}
	}()
	result := new(RPC.Result)
	Event := RPC.Event{EventAsString: string(event)}
	if err := rpcI.client.Call("RuleSystemRPCHandler.CreateEvent", &Event, &result); err != nil {
		return err
	}
	return result.Err
}

//CreateDelayedEvent encapsulates the RPC call to create a DelayedEvent on the server
func (rpcI RPCInterface) CreateDelayedEvent(event []byte, delay *time.Duration) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = rec.(error)
		}
	}()
	result := new(RPC.Result)
	Event := RPC.Event{EventAsString: string(event), Delay: delay}
	if err := rpcI.client.Call("RuleSystemRPCHandler.CreateEvent", &Event, &result); err != nil {
		return err
	}
	return result.Err
}

//SendMessage sends a message to the logserver
func (rpcI RPCInterface) SendMessage(message *LogServer.LogMessage) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = rec.(error)
		}
	}()
	result := new(RPC.Result)
	if err := rpcI.client.Call("LogServerRPCHandler.SendMessage", message, &result); err != nil {
		return err
	}
	return result.Err
}

//SendMessages sends a multiple messages to the logserver
func (rpcI RPCInterface) SendMessages(messages *[]*LogServer.LogMessage) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = rec.(error)
		}
	}()
	result := new(RPC.Result)
	if err := rpcI.client.Call("LogServerRPCHandler.SendMessages", messages, &result); err != nil {
		return err
	}
	return result.Err
}

//MakeCall executes a given Module on the other side
func (rpcI RPCInterface) MakeCall(command string, event []byte) (result *Module.Result, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = rec.(error)
		}
	}()
	result = new(Module.Result)
	call := &RPC.Call{Event: &RPC.Event{EventAsString: string(event)}, Module: command}
	err = rpcI.client.Call("ProxyRPCHandler.Call", call, &result)
	return result, err

}
