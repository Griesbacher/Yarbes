package LogServer

import (
	"github.com/kdar/factorlog"
	"time"
)

//LogMessage represents a single Message which can be send over network
type LogMessage struct {
	Timestamp time.Time
	Source    string
	LogLevel  factorlog.Severity
	Message   string
}

//NewLogMessage constructs a new LogMessage
func NewLogMessage(source, message string, level factorlog.Severity) *LogMessage {
	return &LogMessage{Timestamp: time.Now(), Source: source, LogLevel: level, Message: message}
}

//NewDebugLogMessage constructs a new LogMessage with level debug
func NewDebugLogMessage(source, message string) *LogMessage {
	return NewLogMessage(source, message, factorlog.DEBUG)
}

//NewInfoLogMessage constructs a new LogMessage with level info
func NewInfoLogMessage(source, message string) *LogMessage {
	return NewLogMessage(source, message, factorlog.INFO)
}

//NewWarnLogMessage constructs a new LogMessage with level warn
func NewWarnLogMessage(source, message string) *LogMessage {
	return NewLogMessage(source, message, factorlog.WARN)
}

//NewErrorLogMessage constructs a new LogMessage with level error
func NewErrorLogMessage(source, message string) *LogMessage {
	return NewLogMessage(source, message, factorlog.ERROR)
}
