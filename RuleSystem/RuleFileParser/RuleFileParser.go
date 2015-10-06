package RuleFileParser
import (
	"io/ioutil"
	"fmt"
	"strings"
	"strconv"
	"github.com/griesbacher/SystemX/Module"
	"github.com/griesbacher/SystemX/Event"
)

type RuleFileParser struct {
	ruleFile       string
	lines          []RuleLine
	externalModule Module.ExternalModule
}

func NewRuleFileParser(ruleFile string) *RuleFileParser {
	fileContent, err := ioutil.ReadFile(ruleFile)
	if err != nil {
		panic(err)
	}

	lines := []RuleLine{}
	for index, line := range strings.SplitAfter(string(fileContent), "\n") {
		line = strings.TrimSpace(line)
		elements := strings.Split(line, ";")
		if len(elements) != 4 {
			panic(fmt.Sprintf("Error in Line: &d", index))
		}
		last, err := strconv.ParseBool(elements[3])
		if err != nil {
			panic(err)
		}
		lines = append(lines,
			RuleLine{name:elements[0],
				condition:elements[1],
				command:elements[2],
				last:last})
	}
	return &RuleFileParser{ruleFile: ruleFile, lines:lines, externalModule:*Module.GetExternalModule()}
}

func (rule RuleFileParser)EvaluateJson(event Event.Event) {
	currentEvent := event
	for _, line := range rule.lines {
		fmt.Print(line.name+" ")
		valid, err := line.EvaluateLine(currentEvent)
		if err != nil {
			panic(err)
		}

		if valid {
			fmt.Println(valid)
			newEvent, err := rule.externalModule.Call(line.command, currentEvent)
			if err != nil {
				panic(err)
			}else {
				fmt.Println(newEvent)
				currentEvent = *newEvent
				if line.last {
					break
				}
			}
		}
	}
}
