package RuleFileParser

import (
	"fmt"
	"github.com/griesbacher/SystemX/Event"
	"github.com/griesbacher/SystemX/RuleSystem/RuleFileParser/ConditionParser"
	"regexp"
	"strconv"
	"strings"
)

//RuleLine represents a single rule in a Rulefile
type RuleLine struct {
	name      string
	condition string
	command   string
	args      string
	flags     map[string]string
	parser    ConditionParser.ConditionParser
}

var commandLayout = regexp.MustCompile("(.*?)\\((.*?)\\)")

//EvaluateLine returns if state of the condition and an error if the result is not valid
func (line RuleLine) EvaluateLine(event Event.Event) (bool, error) {
	return line.parser.ParseString(line.condition, event.GetDataAsInterface())
}

//LastLine returns true if the current line is the last line to check
func (line RuleLine) LastLine() bool {
	if value, ok := line.flags["last"]; ok {
		result, _ := strconv.ParseBool(value)
		return result
	}
	return false
}

func parseCommand(commandField string) (string, string, error) {
	if strings.ContainsAny(commandField, "()") {
		hits := commandLayout.FindStringSubmatch(commandField)
		if len(hits) == 3 {
			return hits[1], hits[2], nil
		}
		return "", "", fmt.Errorf("Commandfield is malformated: %s", commandField)
	}
	return commandField, "", nil
}
