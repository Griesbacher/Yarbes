package RuleFileParser

import (
	"encoding/csv"
	"fmt"
	"github.com/griesbacher/Yarbes/Config"
	"github.com/griesbacher/Yarbes/Event"
	"github.com/griesbacher/Yarbes/Logging"
	"github.com/griesbacher/Yarbes/Module"
	"github.com/griesbacher/Yarbes/RuleSystem/RuleFileParser/ConditionParser"
	"github.com/griesbacher/nagflux/helper"
	"os"
)

//RuleFileParser represents a single rule file
type RuleFileParser struct {
	ruleFile       string
	lines          []RuleLine
	externalModule Module.ExternalModule
	LogClient      *Logging.Client
}

//NewRuleFileParser creates a new RuleFileParser, returns an error if the object is not valid
func NewRuleFileParser(ruleFile string) (*RuleFileParser, error) {
	fileReader, err := os.Open(ruleFile)
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(fileReader)
	r.TrimLeadingSpace = true
	r.Comment = '#'
	r.Comma = ';'
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	lines := []RuleLine{}
	for _, elements := range records {
		command, args, err := parseCommand(elements[2])
		if err != nil {
			return nil, err
		}

		if len(elements[0]) == 0 && len(elements[1]) == 0 {
			if len(lines) == 0 {
				return nil, fmt.Errorf("The first rule can not referance back")
			}
			lastRule := lines[len(lines)-1]
			lines = append(lines, RuleLine{name: lastRule.name,
				condition: lastRule.condition,
				command:   command,
				args:      args,
				flags:     helper.StringToMap(elements[3], ",", "="),
			})
		} else {
			lines = append(lines,
				RuleLine{name: elements[0],
					condition: elements[1],
					command:   command,
					args:      args,
					flags:     helper.StringToMap(elements[3], ",", "="),
				})
		}
	}
	client, err := Logging.NewClient(Config.GetClientConfig().LogServer.RPCInterface)
	if err != nil {
		return nil, err
	}
	return &RuleFileParser{ruleFile: ruleFile, lines: lines, externalModule: *Module.NewExternalModule(), LogClient: client}, nil
}

//EvaluateJSON will be called if a new Event occurred an the rulefile will be executed
func (rule RuleFileParser) EvaluateJSON(event Event.Event) {
	currentEvent := event
	eventMetadata := map[string]interface{}{"executedLines": 0}
	for _, line := range rule.lines {
		fmt.Print(line.name + " ")
		valid, err := line.EvaluateLine(currentEvent, eventMetadata)
		if err != nil {
			if err == ConditionParser.ErrElementNotFound {
				valid = false
			} else {
				rule.LogClient.Warn("EvaluteLine:" + err.Error())
			}
		}

		fmt.Println(valid)
		if valid {
			eventMetadata["executedLines"] = eventMetadata["executedLines"].(int) + 1
			moduleResult, err := rule.externalModule.Call(line.command, line.args, currentEvent.String())
			if err != nil {
				rule.LogClient.Error(err)
			} else {
				if moduleResult != nil {
					rule.LogClient.Debug("Module Result: ", *moduleResult)
					//If the module provides a new Event replace the old one
					if moduleResult.Event != nil {
						var newEvent *Event.Event
						newEvent, err = Event.NewEventFromInterface(moduleResult.Event)
						if err != nil {
							rule.LogClient.Warn("NewEventFromInterface: " + err.Error())
						}
						currentEvent = *newEvent
					}

					messages := moduleResult.DecodeLogMessages()
					if len(*messages) > 0 {
						rule.LogClient.LogMultiple(moduleResult.DecodeLogMessages())
					}
				}

				if line.LastLine() {
					break
				}
			}
		}
	}
}
