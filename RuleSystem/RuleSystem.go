package RuleSystem

import (
	"github.com/griesbacher/Yarbes/Config"
	"github.com/griesbacher/Yarbes/Event"
	"github.com/griesbacher/Yarbes/RuleSystem/RuleFileParser"
	"sync"
	"time"
)

//RuleSystem is a deamonlike struct which recives events and executes modules depending on the rules
type RuleSystem struct {
	EventQueue    chan Event.Event
	workers       []ruleSystemWorker
	quit          chan bool
	isRunning     bool
	delayedEvents []*Event.DelayedEvent
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

//AddDelayedEvent creates a DelayedEvent and starts the countdown
func (system *RuleSystem) AddDelayedEvent(event *Event.Event, delay time.Duration) {
	delayedEvent := Event.NewDelayedEvent(event, delay, system.EventQueue)
	system.delayedEvents = append(system.delayedEvents, delayedEvent)
	delayedEvent.Start()
	system.clearDelayedEvents()
}

//GetDelayedEvent returns the list of DelayedEvents which are still waiting
func (system RuleSystem) GetDelayedEvent() []*Event.DelayedEvent {
	system.clearDelayedEvents()
	return system.delayedEvents
}

var mutex = &sync.Mutex{}

func (system *RuleSystem) clearDelayedEvents() {
	mutex.Lock()
	old := system.delayedEvents
	result := []*Event.DelayedEvent{}
	for _, event := range old {
		if event.IsWaiting() {
			result = append(result, event)
		}
	}
	system.delayedEvents = result
	mutex.Unlock()
}
