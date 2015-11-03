package Incoming

import (
	"github.com/griesbacher/Yarbes/Module"
	"github.com/griesbacher/Yarbes/NetworkInterfaces"
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
	callResult, err := handler.external.Call(call.Module, "", call.EventAsString)
	result.Event = callResult.Event

	result.RemoteReturnCode = callResult.ReturnCode
	result.Messages = callResult.Messages
	return err
}