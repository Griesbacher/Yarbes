package RuleSystem

import (
	"github.com/griesbacher/SystemX/Config"
	"github.com/griesbacher/SystemX/Event"
	"github.com/griesbacher/SystemX/RuleSystem/RuleFileParser"
)

type RuleSystem struct {
	EventQueue chan Event.Event
	workers    []ruleSystemWorker
	quit       chan bool
}

func NewRuleSystem() *RuleSystem {
	eventQueue := make(chan Event.Event, 1000)

	parser, err := *RuleFileParser.NewRuleFileParser(Config.GetServerConfig().RuleSystem.Rulefile)
	if err != nil {
		panic(err)
	}
	amountOfWorker := Config.GetServerConfig().RuleSystem.Worker
	workers := []ruleSystemWorker{}
	for i := 0; i < amountOfWorker; i++ {
		workers = append(workers, ruleSystemWorker{eventQueue: eventQueue, parser: parser, quit: make(chan bool), isRunning: false})
	}
	return &RuleSystem{EventQueue: eventQueue, workers: workers}
}

func (system RuleSystem) Start() {
	for _, worker := range system.workers {
		worker.Start()
	}
}

func (system RuleSystem) Stop() {
	for _, worker := range system.workers {
		worker.Stop()
	}
}

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
