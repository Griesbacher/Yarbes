package Incoming

import (
	"github.com/griesbacher/Yarbes/Config"
	"net/http"
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
	http.HandleFunc("/logs", logHandler.LogView)
	return ruleHTTP
}
