package RuleSystem

import (
	"github.com/griesbacher/SystemX/Config"
	"github.com/griesbacher/SystemX/Event"
	"github.com/griesbacher/SystemX/RuleSystem/RuleFileParser"
)

//RuleSystem is a deamonlike struct which recives events and executes modules depending on the rules
type RuleSystem struct {
	EventQueue chan Event.Event
	workers    []ruleSystemWorker
	quit       chan bool
	isRunning  bool
}

//NewRuleSystem is the constructor
func NewRuleSystem() *RuleSystem {
	eventQueue := make(chan Event.Event, 1000)

	parser, err := RuleFileParser.NewRuleFileParser(Config.GetServerConfig().RuleSystem.Rulefile)
	if err != nil {
		panic(err)
	}
	amountOfWorker := Config.GetServerConfig().RuleSystem.Worker
	workers := []ruleSystemWorker{}
	for i := 0; i < amountOfWorker; i++ {
		workers = append(workers, ruleSystemWorker{eventQueue: eventQueue, parser: *parser, quit: make(chan bool), isRunning: false})
	}
	return &RuleSystem{EventQueue: eventQueue, workers: workers, isRunning: false}
}

//Start starts the RuleSystem with its workers
func (system *RuleSystem) Start() {
	if !system.isRunning {
		for _, worker := range system.workers {
			worker.Start()
		}
		system.isRunning = true
	}
}

//Stop stops the RuleSystem with its workers
func (system RuleSystem) Stop() {
	for _, worker := range system.workers {
		worker.Stop()
	}
	system.isRunning = false
}

//IsRunning returns true if the daemon is running
func (system RuleSystem) IsRunning() bool {
	return system.isRunning
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
