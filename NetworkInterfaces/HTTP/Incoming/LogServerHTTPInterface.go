package Incoming

import (
	"github.com/griesbacher/Yarbes/Config"
)

//LogServerHTTPInterface is HTTP interface which offers log access
type LogServerHTTPInterface struct {
	*HTTPInterface
	logHandler LogServerHTTPHandler
}

//NewLogServerHTTPInterface creates a new LogServerHTTPInterface
func NewLogServerHTTPInterface() *LogServerHTTPInterface {
	httpI := NewHTTPInterface(Config.GetServerConfig().LogServer.HTTPInterface)
	logHandler := LogServerHTTPHandler{}
	ruleHTTP := &LogServerHTTPInterface{HTTPInterface: httpI, logHandler: logHandler}
	ruleHTTP.HTTPInterface.PublishHandler("/logs", logHandler.LogView)
	return ruleHTTP
}
