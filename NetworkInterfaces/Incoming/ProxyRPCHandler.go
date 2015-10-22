package Incoming

import (
	"github.com/griesbacher/SystemX/Module"
	"github.com/griesbacher/SystemX/NetworkInterfaces"
)

//ProxyRPCHandler is a RPC handler which accepts LogMessages
type ProxyRPCHandler struct {
	external *Module.ExternalModule
}

//Call executes the given script and returns the result
func (handler *ProxyRPCHandler) Call(call *NetworkInterfaces.RPCCall, result *Module.Result) error {
	if call == nil {
		return ErrorInputWasNil
	} else if result == nil {
		return ErrorResultWasNil
	}
	callResult, err := handler.external.Call(call.Module, call.EventAsString)
	result.Event = callResult.Event

	result.ReturnCode = callResult.ReturnCode
	result.LogMessages = callResult.LogMessages
	return err
}
