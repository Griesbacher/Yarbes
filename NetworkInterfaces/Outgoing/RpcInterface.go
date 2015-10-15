package Outgoing

import (
	"crypto/tls"
	"errors"
	"github.com/griesbacher/SystemX/Config"
	"github.com/griesbacher/SystemX/NetworkInterfaces"
	"github.com/griesbacher/SystemX/TLS"
	"net/rpc"
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

//GenRPCClient sends a single byte for authentication and returns a rpc.Client
//TODO: durch methoden ersetzen
func (rpcI RPCInterface) GenRPCClient() *rpc.Client {
	rpcI.conn.Write([]byte("a"))
	return rpc.NewClient(rpcI.conn)
}

//Disconnect closes the tcp connection
func (rpcI RPCInterface) Disconnect() {
	rpcI.conn.Close()
}

//CreateEvent  encapsulates the RPC call to create a event on the server
func (rpcI RPCInterface) CreateEvent(event []byte) error {
	result := new(NetworkInterfaces.RPCResult)
	if err := rpcI.client.Call("RuleSystemRPCHandler.CreateEvent", string(event), &result); err != nil {
		return err
	}
	if result.Err != nil {
		return result.Err
	}
	return nil
}
