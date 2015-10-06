package main

import (
	"github.com/griesbacher/SystemX/Event"
	"github.com/griesbacher/SystemX/RuleSystem"
	"time"
)

func main() {
	b := []byte(`{
   "k1" : "v1",
   "k2" : 10,
   "k3" : ["v4",12.3,{"k11" : "v11", "k22" : "v22"}]
}`)


	event, err := Event.NewEvent(b)
	if err != nil {
		panic(err)
	}

	ruleSystem := RuleSystem.NewRuleSystem()
	ruleSystem.Start()
	ruleSystem.EventQueue <- *event

	time.Sleep(time.Duration(5)*time.Second)
}
