package Incoming

import (
	"github.com/griesbacher/SystemX/Config"
	"github.com/griesbacher/SystemX/Event"
	"github.com/griesbacher/SystemX/NetworkInterfaces"
)

type RuleSystemRPCInterface struct {
	rpcInterface *RPCInterface
	eventQueue   chan Event.Event
}

func NewRuleSystemRPCInterface(eventQueue chan Event.Event) *RuleSystemRPCInterface {
	rpc := NewRPCInterface(Config.GetServerConfig().RuleSystem.RPCInterface)
	ruleRPC := &RuleSystemRPCInterface{rpcInterface: rpc, eventQueue: eventQueue}
	rpc.publishHandler(&RuleSystemRPCHandler{*ruleRPC})
	return ruleRPC
}

func (rpcI RuleSystemRPCInterface) Start() {
	rpcI.rpcInterface.Start()
}

func (rpcI RuleSystemRPCInterface) Stop() {
	rpcI.rpcInterface.Stop()
}

type RuleSystemRPCHandler struct {
	inter RuleSystemRPCInterface
}

func (handler *RuleSystemRPCHandler) CreateEvent(args *string, result *NetworkInterfaces.RPCResult) error {
	event, err := Event.NewEventFromBytes([]byte(*args))
	if err == nil {
		handler.inter.eventQueue <- *event
	}
	result.Err = err
	return err
}
