package Incoming

import (
	"github.com/griesbacher/Yarbes/Config"
	"github.com/griesbacher/Yarbes/RuleSystem"
)

//RuleSystemRPCInterface offers a RPC interface to creates Events
type RuleSystemRPCInterface struct {
	*RPCInterface
	ruleSystem *RuleSystem.RuleSystem
}

//NewRuleSystemRPCInterface creates a new RuleSystemRPCInterface
func NewRuleSystemRPCInterface(ruleSystem *RuleSystem.RuleSystem) *RuleSystemRPCInterface {
	rpc := NewRPCInterface(Config.GetServerConfig().RuleSystem.RPCInterface)
	ruleRPC := &RuleSystemRPCInterface{RPCInterface: rpc, ruleSystem: ruleSystem}
	rpc.publishHandler(&RuleSystemRPCHandler{*ruleRPC})
	return ruleRPC
}
