package RuleFileParser
import (
	"github.com/griesbacher/SystemX/Event"
)

type RuleLine struct {
	name      string
	condition string
	command   string
	last      bool
	parser    ConditionParser
}

func NewRuleLine(name, condition, command string, last bool) RuleLine {
	return RuleLine{name, condition, command, last, ConditionParser{}}
}

func (line RuleLine)EvaluateLine(event Event.Event) (bool, error) {
	return line.parser.ParseString(line.condition, event.GetDataAsInterface())
}