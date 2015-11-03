package Incoming

import (
	"github.com/griesbacher/Yarbes/Config"
	"github.com/griesbacher/Yarbes/Logging/LogServer"
)

//LogServerRPCInterface is RPC interface which offers logging
type LogServerRPCInterface struct {
	*RPCInterface
	logQueue chan LogServer.LogMessage
}

//NewLogServerRPCInterface creates a new LogServerRPCInterface
func NewLogServerRPCInterface(logQueue chan LogServer.LogMessage) *LogServerRPCInterface {
	rpc := NewRPCInterface(Config.GetServerConfig().LogServer.RPCInterface)
	ruleRPC := &LogServerRPCInterface{RPCInterface: rpc, logQueue: logQueue}
	rpc.publishHandler(&LogServerRPCHandler{*ruleRPC})
	return ruleRPC
}
