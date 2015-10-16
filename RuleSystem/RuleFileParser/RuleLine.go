package RuleFileParser

import (
	"github.com/griesbacher/SystemX/Event"
	"github.com/griesbacher/SystemX/RuleSystem/RuleFileParser/ConditionParser"
	"strconv"
)

//RuleLine represents a single rule in a Rulefile
type RuleLine struct {
	name      string
	condition string
	command   string
	flags     map[string]string
	parser    ConditionParser.ConditionParser
}

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
