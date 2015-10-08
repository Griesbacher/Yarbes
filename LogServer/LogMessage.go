package LogServer
import (
	"time"
	"github.com/kdar/factorlog"
)

type LogMessage struct {
	Timestamp time.Time
	Source    string
	LogLevel  factorlog.Severity
	Message   string
}

func NewLogMessage(source, message string) *LogMessage {
	return &LogMessage{Timestamp:time.Now(), Source:source, LogLevel:factorlog.DEBUG, Message:message}
}