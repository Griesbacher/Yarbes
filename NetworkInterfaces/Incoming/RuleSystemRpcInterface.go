package Incoming

import (
	"github.com/griesbacher/SystemX/Config"
	"github.com/griesbacher/SystemX/Event"
	"github.com/griesbacher/SystemX/NetworkInterfaces"
)

//RuleSystemRPCInterface offers a RPC interface to creates Events
type RuleSystemRPCInterface struct {
	*RPCInterface
	eventQueue chan Event.Event
}

//NewRuleSystemRPCInterface creates a new RuleSystemRPCInterface
func NewRuleSystemRPCInterface(eventQueue chan Event.Event) *RuleSystemRPCInterface {
	rpc := NewRPCInterface(Config.GetServerConfig().RuleSystem.RPCInterface)
	ruleRPC := &RuleSystemRPCInterface{RPCInterface: rpc, eventQueue: eventQueue}
	rpc.publishHandler(&RuleSystemRPCHandler{*ruleRPC})
	return ruleRPC
}

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
