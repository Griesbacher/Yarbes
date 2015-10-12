package RuleFileParser

import (
	"fmt"
	"github.com/griesbacher/SystemX/Event"
	"github.com/griesbacher/SystemX/LogServer"
	"github.com/griesbacher/SystemX/Module"
	"github.com/griesbacher/nagflux/helper"
	"io/ioutil"
	"strings"
)

type RuleFileParser struct {
	ruleFile       string
	lines          []RuleLine
	externalModule Module.ExternalModule
	logClient      *LogServer.Client
}

func NewRuleFileParser(ruleFile string) (*RuleFileParser, error) {
	fileContent, err := ioutil.ReadFile(ruleFile)
	if err != nil {
		return nil, err
	}

	lines := []RuleLine{}
	for index, line := range strings.SplitAfter(string(fileContent), "\n") {
		line = strings.TrimSpace(line)
		if len(line) == 0 || string(line[0]) == "#" {
			continue
		}
		elements := strings.Split(line, ";")
		if len(elements) != 4 {
			return nil, fmt.Errorf("Number of Elements are not four in line: %d", index)
		}
		if len(elements[0]) == 0 && len(elements[1]) == 0 {
			if len(lines) == 0 {
				return nil, fmt.Errorf("The first rule can not referance back")
			}
			lastRule := lines[len(lines)-1]
			lines = append(lines, RuleLine{name: lastRule.name,
				condition: lastRule.condition,
				command:   elements[2],
				flags:     helper.StringToMap(elements[3], ",", "="),
			})
		} else {
			lines = append(lines,
				RuleLine{name: elements[0],
					condition: elements[1],
					command:   elements[2],
					flags:     helper.StringToMap(elements[3], ",", "="),
				})
		}
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
				if line.LastLine() {
					break
				}
			}
		}
	}
}
