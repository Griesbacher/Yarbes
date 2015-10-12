package RuleFileParser

import (
	"fmt"
	"github.com/griesbacher/SystemX/Event"
	"github.com/griesbacher/SystemX/LogServer"
	"github.com/griesbacher/SystemX/Module"
	"io/ioutil"
	"strconv"
	"strings"
)

type RuleFileParser struct {
	ruleFile       string
	lines          []RuleLine
	externalModule Module.ExternalModule
	logClient      LogServer.Client
}

func NewRuleFileParser(ruleFile string) (*RuleFileParser, error) {
	fileContent, err := ioutil.ReadFile(ruleFile)
	if err != nil {
		return nil, err
	}

	lines := []RuleLine{}
	for index, line := range strings.SplitAfter(string(fileContent), "\n") {
		line = strings.TrimSpace(line)
		elements := strings.Split(line, ";")
		if len(elements) != 4 {
			return nil, fmt.Errorf("Number of Elements are not four in line: %d", index)
		}
		last, err := strconv.ParseBool(elements[3])
		if err != nil {
			return nil, fmt.Errorf("Could not parse bool in line: %d", index)
		}
		lines = append(lines,
			RuleLine{name: elements[0],
				condition: elements[1],
				command:   elements[2],
				last:      last})
	}
	client, err := LogServer.NewClient()
	if err != nil {
		return nil, err
	}
	return &RuleFileParser{ruleFile: ruleFile, lines: lines, externalModule: *Module.GetExternalModule(), logClient: client}, nil
}

func (rule RuleFileParser) EvaluateJSON(event Event.Event) {
	currentEvent := event
	for _, line := range rule.lines {
		fmt.Print(line.name + " ")
		valid, err := line.EvaluateLine(currentEvent)
		if err != nil {
			rule.logClient.Debug(err.Error())
		}

		if valid {
			fmt.Println(valid)
			newEvent, err := rule.externalModule.Call(line.command, currentEvent)
			if err != nil {
				rule.logClient.Debug(err.Error())
			} else {
				fmt.Println(newEvent)
				currentEvent = *newEvent
				if line.last {
					break
				}
			}
		}
	}
}
