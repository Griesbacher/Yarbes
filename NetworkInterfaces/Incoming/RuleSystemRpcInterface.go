package Incoming

import (
	"github.com/griesbacher/SystemX/Event"
	"github.com/griesbacher/SystemX/Config"
	"github.com/griesbacher/SystemX/NetworkInterfaces"
)

type RuleSystemRpcInterface struct {
	rpcInterface *RpcInterface
	eventQueue   chan Event.Event
}

func NewRuleSystemRpcInterface(eventQueue chan Event.Event) *RuleSystemRpcInterface {
	rpc := NewRpcInterface(Config.GetServerConfig().RuleSystem.RpcInterface)
	ruleRpc := &RuleSystemRpcInterface{rpcInterface:rpc, eventQueue:eventQueue}
	rpc.publishHandler(&RuleSystemRpcHandler{*ruleRpc})
	return ruleRpc
}

func (rpcI RuleSystemRpcInterface) Start() {
	rpcI.rpcInterface.Start()
}

func (rpcI RuleSystemRpcInterface) Stop() {
	rpcI.rpcInterface.Stop()
}

type RuleSystemRpcHandler struct {
	inter RuleSystemRpcInterface
}

func (handler *RuleSystemRpcHandler) CreateEvent(args *string, result *NetworkInterfaces.RpcResult) error {
	event, err := Event.NewEvent([]byte(*args))
	if err == nil {
		handler.inter.eventQueue <- *event
	}
	result.Err = err
	return err
}
