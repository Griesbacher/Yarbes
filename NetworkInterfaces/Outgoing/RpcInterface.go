package Outgoing

import (
	"crypto/tls"
	"github.com/griesbacher/SystemX/Config"
	"github.com/griesbacher/SystemX/TLS"
	"net/rpc"
)

type RPCInterface struct {
	serverAddress string
	Config        *tls.Config
	conn          *tls.Conn
}

func NewRPCInterface(serverAddress string) RPCInterface {
	config := TLS.GenerateClientTLSConfig(Config.GetClientConfig().TLS.Cert, Config.GetClientConfig().TLS.Key, Config.GetClientConfig().TLS.CaCert)
	return RPCInterface{serverAddress: serverAddress, Config: config}
}

func (rpcI *RPCInterface) Connect() error {
	conn, err := tls.Dial("tcp", rpcI.serverAddress, rpcI.Config)
	rpcI.conn = conn
	return err
}

func (rpcI RPCInterface) GenRPCClient() *rpc.Client {
	rpcI.conn.Write([]byte("a"))
	return rpc.NewClient(rpcI.conn)
}

func (rpcI RPCInterface) Disconnect() {
	rpcI.conn.Close()
}
