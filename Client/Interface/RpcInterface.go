package Interface

import (
	"crypto/tls"
	"github.com/griesbacher/SystemX/Config"
	"github.com/griesbacher/SystemX/TLS"
	"net/rpc"
)

type RpcInterface struct {
	serverAddress string
	Config        tls.Config
	conn          *tls.Conn
}

func NewRpcInterface(serverAddress string) RpcInterface {
	config := TLS.GenerateClientTLSConfig(Config.GetClientConfig().Client.TLSCert, Config.GetClientConfig().Client.TLSKey, Config.GetClientConfig().Client.TLSCaCert)
	return RpcInterface{serverAddress: serverAddress, Config: config}
}

func (rpcI *RpcInterface) Connect() error {
	conn, err := tls.Dial("tcp", rpcI.serverAddress, &rpcI.Config)
	rpcI.conn = conn
	return err
}

func (rpcI RpcInterface) GenRpcClient() *rpc.Client {
	return rpc.NewClient(rpcI.conn)
}

func (rpcI RpcInterface) Disconnect() {
	rpcI.conn.Close()
}
