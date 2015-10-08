package Incoming

import (
	"github.com/griesbacher/SystemX/Config"
	"github.com/griesbacher/SystemX/LogServer"
	"github.com/griesbacher/SystemX/NetworkInterfaces"
)

type LogServerRpcInterface struct {
	rpcInterface *RpcInterface
	logQueue     chan LogServer.LogMessage
}

func NewLogServerRpcInterface(logQueue chan LogServer.LogMessage) *LogServerRpcInterface {
	rpc := NewRpcInterface(Config.GetServerConfig().LogServer.RpcInterface)
	ruleRpc := &LogServerRpcInterface{rpcInterface: rpc, logQueue: logQueue}
	rpc.publishHandler(&LogServerRpcHandler{*ruleRpc})
	return ruleRpc
}

func (rpcI LogServerRpcInterface) Start() {
	rpcI.rpcInterface.Start()
}

func (rpcI LogServerRpcInterface) Stop() {
	rpcI.rpcInterface.Stop()
}

type LogServerRpcHandler struct {
	inter LogServerRpcInterface
}

func (handler *LogServerRpcHandler) SendMessages(messages *[]LogServer.LogMessage, result *NetworkInterfaces.RpcResult) error {
	for _, message := range *messages {
		handler.SendMessage(&message, result)
	}
	return nil
}

func (handler *LogServerRpcHandler) SendMessage(message *LogServer.LogMessage, result *NetworkInterfaces.RpcResult) error {
	handler.inter.logQueue <- *message
	return nil
}
