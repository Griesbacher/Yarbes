package Incoming

import (
	"github.com/griesbacher/Yarbes/Event"
	"github.com/griesbacher/Yarbes/NetworkInterfaces/RPC"
)

//RuleSystemRPCHandler RPC handler to create Events
type RuleSystemRPCHandler struct {
	inter RuleSystemRPCInterface
}

//CreateEvent creates a event from the given string and sends it to the RuleSystem
func (handler *RuleSystemRPCHandler) CreateEvent(rpcEvent *RPC.Event, result *RPC.Result) error {
	if rpcEvent == nil {
		return ErrorInputWasNil
	} else if result == nil {
		return ErrorResultWasNil
	}

	event, err := Event.NewEventFromBytes([]byte(rpcEvent.EventAsString))
	if rpcEvent.Delay == nil {
		if err == nil {
			handler.inter.ruleSystem.EventQueue <- *event
		}
	} else {
		handler.inter.ruleSystem.AddDelayedEvent(event, *rpcEvent.Delay)
	}
	result.Err = err
	return err
}
