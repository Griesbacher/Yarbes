package Outgoing

import (
	"crypto/tls"
	"github.com/griesbacher/SystemX/Config"
	"github.com/griesbacher/SystemX/TLS"
	"net/rpc"
)

//RPCInterface represents a outgoing RPC connection, with which a rpc.Client can be created
type RPCInterface struct {
	serverAddress string
	Config        *tls.Config
	conn          *tls.Conn
}

//NewRPCInterface constructs a new RPCInterface
func NewRPCInterface(serverAddress string) RPCInterface {
	config := TLS.GenerateClientTLSConfig(Config.GetClientConfig().TLS.Cert, Config.GetClientConfig().TLS.Key, Config.GetClientConfig().TLS.CaCert)
	return RPCInterface{serverAddress: serverAddress, Config: config}
}

//Connect establishes a tcp connection
func (rpcI *RPCInterface) Connect() error {
	conn, err := tls.Dial("tcp", rpcI.serverAddress, rpcI.Config)
	rpcI.conn = conn
	return err
}

//GenRPCClient sends a single byte for authentication and returns a rpc.Client
func (rpcI RPCInterface) GenRPCClient() *rpc.Client {
	rpcI.conn.Write([]byte("a"))
	return rpc.NewClient(rpcI.conn)
}

//Disconnect closes the tcp connection
func (rpcI RPCInterface) Disconnect() {
	rpcI.conn.Close()
}
