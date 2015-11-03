package Incoming

import (
	"github.com/griesbacher/Yarbes/Config"
	"github.com/griesbacher/Yarbes/Module"
)

//ProxyRPCInterface is RPC interface which offers module execution
type ProxyRPCInterface struct {
	*RPCInterface
}

//NewProxyRPCInterface creates a new ProxyRPCInterface
func NewProxyRPCInterface() *ProxyRPCInterface {
	rpc := NewRPCInterface(Config.GetServerConfig().Proxy.RPCInterface)
	ruleRPC := &ProxyRPCInterface{RPCInterface: rpc}
	rpc.publishHandler(&ProxyRPCHandler{Module.NewExternalModule()})
	return ruleRPC
}
