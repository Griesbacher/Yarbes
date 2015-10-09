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

func NewLogMessage(source, message string) *LogMessage {
	return &LogMessage{Timestamp: time.Now(), Source: source, LogLevel: factorlog.DEBUG, Message: message}
}
