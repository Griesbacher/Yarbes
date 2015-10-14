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
		Severity  string
		Message   string
		Source    string
	}
}

//TimeParseLayout is the format in which the timestamps, within the JSON, are expected RFC3339
const TimeParseLayout = "2015-10-14T08:26:26+02:00"

//DecodeLogMessages converts the LogMessages from the JSON object to LogServer.LogMessages
func (moduleResult Result) DecodeLogMessages() *[]*LogServer.LogMessage {
	result := []*LogServer.LogMessage{}
	for _, message := range moduleResult.LogMessages {
		level := factorlog.StringToSeverity(strings.ToUpper(message.Severity))
		var timestamp time.Time

		if newTime, err := time.Parse(TimeParseLayout, message.Timestamp); err == nil {
			timestamp = newTime
		} else {
			timestamp = time.Now()
		}

		result = append(result, &LogServer.LogMessage{
			Timestamp: timestamp,
			Source:    message.Source,
			Severity:  level,
			Message:   message.Message,
		})
	}
	return &result
}
