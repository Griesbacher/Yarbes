package Incoming

import (
	"github.com/griesbacher/SystemX/Logging/LogServer"
	"github.com/griesbacher/SystemX/NetworkInterfaces"
)

//LogServerRPCHandler is a RPC handler which accepts LogMessages
type LogServerRPCHandler struct {
	inter LogServerRPCInterface
}

//SendMessages takes a list of LogMessages
func (handler *LogServerRPCHandler) SendMessages(messages *[]*LogServer.LogMessage, result *NetworkInterfaces.RPCResult) error {
	for _, message := range *messages {
		handler.SendMessage(message, result)
	}
	return nil
}

//SendMessage takes a single LogMessage
func (handler *LogServerRPCHandler) SendMessage(message *LogServer.LogMessage, result *NetworkInterfaces.RPCResult) error {
	handler.inter.logQueue <- *message
	return nil
}
