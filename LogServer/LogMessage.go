package LogServer

import (
	"github.com/kdar/factorlog"
	"time"
)

type LogMessage struct {
	Timestamp time.Time
	Source    string
	LogLevel  factorlog.Severity
	Message   string
}

func NewLogMessage(source, message string, level factorlog.Severity) *LogMessage {
	return &LogMessage{Timestamp: time.Now(), Source: source, LogLevel: level, Message: message}
}

func NewDebugLogMessage(source, message string) *LogMessage {
	return NewLogMessage(source, message, factorlog.DEBUG)
}
