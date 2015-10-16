package Local

import (
	"github.com/kdar/factorlog"
	"os"
	"sync"
)

const logFormat = "%{Date} %{Time} %{Severity}: %{Message}"
const logColors = "%{Color \"magenta\" \"WARN\"}%{Color \"red\" \"CRITICAL\"}"

var singleLogger *factorlog.FactorLog
var mutex = &sync.Mutex{}

//InitLogger creates a factorlog with the given min loglevel
func InitLogger(minSeverity string) {
	mutex.Lock()
	logFormatter := factorlog.NewStdFormatter(logColors + logFormat)
	targetWriter := os.Stdout
	singleLogger = factorlog.New(targetWriter, logFormatter)
	singleLogger.SetMinMaxSeverity(factorlog.StringToSeverity(minSeverity), factorlog.StringToSeverity("PANIC"))
	mutex.Unlock()
}

//GetLogger returns a factorlog an constructs a new if not exists
func GetLogger() *factorlog.FactorLog {
	if singleLogger == nil {
		InitLogger("DEBUG")
	}
	return singleLogger
}
