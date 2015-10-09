package Incoming

import (
	"github.com/griesbacher/SystemX/Config"
	"github.com/griesbacher/SystemX/LogServer"
	"github.com/griesbacher/SystemX/NetworkInterfaces"
)

type LogServerRPCInterface struct {
	rpcInterface *RPCInterface
	logQueue     chan LogServer.LogMessage
}

func NewLogServerRPCInterface(logQueue chan LogServer.LogMessage) *LogServerRPCInterface {
	rpc := NewRPCInterface(Config.GetServerConfig().LogServer.RPCInterface)
	ruleRPC := &LogServerRPCInterface{rpcInterface: rpc, logQueue: logQueue}
	rpc.publishHandler(&LogServerRPCHandler{*ruleRPC})
	return ruleRPC
}

func (rpcI LogServerRPCInterface) Start() {
	rpcI.rpcInterface.Start()
}

func (rpcI LogServerRPCInterface) Stop() {
	rpcI.rpcInterface.Stop()
}

type LogServerRPCHandler struct {
	inter LogServerRPCInterface
}

func (handler *LogServerRPCHandler) SendMessages(messages *[]LogServer.LogMessage, result *NetworkInterfaces.RPCResult) error {
	for _, message := range *messages {
		handler.SendMessage(&message, result)
	}
	return nil
}

func (handler *LogServerRPCHandler) SendMessage(message *LogServer.LogMessage, result *NetworkInterfaces.RPCResult) error {
	handler.inter.logQueue <- *message
	return nil
}
