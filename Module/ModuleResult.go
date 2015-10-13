package Module

import (
	"github.com/griesbacher/SystemX/Logging/LogServer"
	"strings"
	"github.com/kdar/factorlog"
	"time"
)

const TimeParseLayout = ""

type ModuleResult struct {
	Event       interface{}
	ReturnCode  int
	LogMessages []struct {
		Timestamp string
		Level     string
		Message   string
		Source    string
	}
}

func (moduleResult ModuleResult) DecodeLogMessages() []*LogServer.LogMessage {
	result := []*LogServer.LogMessage{}
	for _, message := range moduleResult.LogMessages {
		var level factorlog.Severity
		switch strings.ToLower(message.Level) {
		case "debug":
			level = factorlog.DEBUG
		case "warn":
			level = factorlog.WARN
		case "error":
			level = factorlog.ERROR
		default:
			level = factorlog.NONE
		}
		var timestamp time.Time

		if newTime, err := time.Parse(TimeParseLayout, message.Timestamp); err == nil {
			timestamp = newTime
		}else {
			timestamp = time.Now()
		}

		result = append(result, &LogServer.LogMessage{
			Timestamp:timestamp,
			Source:message.Source,
			LogLevel:level,
			Message:message.Message,
		})
	}
	return result
}