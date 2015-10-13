package Local

import (
	"github.com/kdar/factorlog"
	"os"
	"io"
)

const logFormat = "%{Date} %{Time} %{Severity}: %{SafeMessage}"
const logColors = "%{Color \"white\" \"DEBUG\"}%{Color \"magenta\" \"WARN\"}%{Color \"red\" \"CRITICAL\"}"

var singleLogger *factorlog.FactorLog = nil

func InitLogger(minSeverity string) {
	var logFormatter factorlog.Formatter
	var targetWriter io.Writer
	var err error
	//logFormatter = factorlog.NewStdFormatter(logColors + logFormat)
	logFormatter = factorlog.NewStdFormatter( logFormat)
	targetWriter = os.Stdout

	if err != nil {
		panic(err)
	}
	singleLogger = factorlog.New(targetWriter, logFormatter)
	singleLogger.SetMinMaxSeverity(factorlog.StringToSeverity(minSeverity), factorlog.StringToSeverity("PANIC"))
}

func GetLogger() *factorlog.FactorLog {
	if singleLogger == nil {
		InitLogger("DEBUG")
	}
	return singleLogger
}