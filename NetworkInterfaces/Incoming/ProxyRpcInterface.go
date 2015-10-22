package Incoming

import (
	"github.com/griesbacher/SystemX/Config"
	"github.com/griesbacher/SystemX/Module"
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
