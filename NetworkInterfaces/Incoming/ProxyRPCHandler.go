package Incoming

import (
	"github.com/griesbacher/SystemX/Event"
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
	newEvent, err := Event.NewEventFromBytes([]byte(call.EventAsString))
	if err != nil {
		return err
	}

	result, err = handler.external.Call(call.Module, *newEvent)
	return err
}
