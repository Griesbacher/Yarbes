package RuleFileParser

import (
	"github.com/griesbacher/SystemX/Event"
	"strconv"
)

type RuleLine struct {
	name      string
	condition string
	command   string
	flags     map[string]string
	parser    ConditionParser
}

func (line RuleLine) EvaluateLine(event Event.Event) (bool, error) {
	return line.parser.ParseString(line.condition, event.GetDataAsInterface())
}

func (line RuleLine) LastLine() bool {
	if value, ok := line.flags["last"]; ok {
		result, _ := strconv.ParseBool(value)
		return result
	}
	return false
}
