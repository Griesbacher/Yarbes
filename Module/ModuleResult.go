package Module

import (
	"github.com/griesbacher/SystemX/Logging/LogServer"
	"github.com/kdar/factorlog"
	"strings"
	"time"
)

//Result represents the output which is expected from a called module
type Result struct {
	Event       interface{}
	ReturnCode  int
	LogMessages []struct {
		Timestamp string
		Level     string
		Message   string
		Source    string
	}
}

//TimeParseLayout is the format in which the timestamps, within the JSON, are expected
const TimeParseLayout = ""

//DecodeLogMessages converts the LogMessages from the JSON object to LogServer.LogMessages
func (moduleResult Result) DecodeLogMessages() *[]*LogServer.LogMessage {
	result := []*LogServer.LogMessage{}
	for _, message := range moduleResult.LogMessages {
		var level factorlog.Severity
		switch strings.ToLower(message.Level) {
		case "debug":
			level = factorlog.DEBUG
		case "info":
			level = factorlog.INFO
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
		} else {
			timestamp = time.Now()
		}

		result = append(result, &LogServer.LogMessage{
			Timestamp: timestamp,
			Source:    message.Source,
			LogLevel:  level,
			Message:   message.Message,
		})
	}
	return &result
}
