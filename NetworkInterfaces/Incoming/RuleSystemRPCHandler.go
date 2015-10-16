package Incoming

import (
	"github.com/griesbacher/SystemX/Event"
	"github.com/griesbacher/SystemX/NetworkInterfaces"
)

//RuleSystemRPCHandler RPC handler to create Events
type RuleSystemRPCHandler struct {
	inter RuleSystemRPCInterface
}

//CreateEvent creates a event from the given string and sends it to the RuleSystem
func (handler *RuleSystemRPCHandler) CreateEvent(args *string, result *NetworkInterfaces.RPCResult) error {
	event, err := Event.NewEventFromBytes([]byte(*args))
	if err == nil {
		handler.inter.eventQueue <- *event
	}
	result.Err = err
	return err
}
