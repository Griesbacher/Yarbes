package RuleSystem

import (
	"github.com/griesbacher/Yarbes/Event"
	"github.com/griesbacher/Yarbes/RuleSystem/RuleFileParser"
)

type ruleSystemWorker struct {
	eventQueue chan Event.Event
	parser     RuleFileParser.RuleFileParser
	quit       chan bool
	isRunning  bool
}

func (worker ruleSystemWorker) Start() {
	if !worker.isRunning {
		go worker.work()
	}
}

func (worker ruleSystemWorker) Stop() {
	worker.quit <- true
	<-worker.quit
	worker.parser.LogClient.Disconnect()
}

func (worker *ruleSystemWorker) work() {
	worker.isRunning = true
	var event Event.Event
	for {
		select {
		case <-worker.quit:
			worker.quit <- true
			return
		case event = <-worker.eventQueue:
			worker.parser.EvaluateJSON(event)
		}
	}
}
