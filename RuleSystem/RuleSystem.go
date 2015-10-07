package RuleSystem
import (
	"github.com/griesbacher/SystemX/RuleSystem/RuleFileParser"
	"github.com/griesbacher/SystemX/Event"
	"github.com/griesbacher/SystemX/Config"
)


type RuleSystem struct {
	EventQueue chan Event.Event
	workers    []ruleSystemWorker
	quit       chan bool
}

func NewRuleSystem() *RuleSystem {
	eventQueue := make(chan Event.Event, 1000)

	parser := *RuleFileParser.NewRuleFileParser(Config.GetConfig().RuleSystem.Rulefile)

	amountOfWorker := Config.GetConfig().RuleSystem.Worker
	workers := []ruleSystemWorker{}
	for i := 0; i < amountOfWorker; i++ {
		workers = append(workers, ruleSystemWorker{eventQueue:eventQueue, parser:parser, quit:make(chan bool)})
	}
	return &RuleSystem{EventQueue:eventQueue, workers:workers}
}

func (system RuleSystem)Start() {
	for _, worker := range system.workers {
		worker.Start()
	}
}

func (system RuleSystem)Stop() {
	system.quit <- true
	for _, worker := range system.workers {
		worker.Stop()
	}
	<-system.quit
}

type ruleSystemWorker struct {
	eventQueue chan Event.Event
	parser     RuleFileParser.RuleFileParser
	quit       chan bool
}

func (worker ruleSystemWorker)Start() {
	go worker.work()
}

func (worker ruleSystemWorker)Stop() {
	worker.quit <- true
	<-worker.quit
}

func (worker ruleSystemWorker) work() {
	var event Event.Event
	for {
		select {
		case <-worker.quit:
			worker.quit <- true
			return
		case event = <-worker.eventQueue:
			worker.parser.EvaluateJson(event)
		}
	}
}