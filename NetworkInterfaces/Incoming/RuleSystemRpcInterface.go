package Incoming

import (
	"github.com/griesbacher/SystemX/Config"
	"github.com/griesbacher/SystemX/Event"
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
