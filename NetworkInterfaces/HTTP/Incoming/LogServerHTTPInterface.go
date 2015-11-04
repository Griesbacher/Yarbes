package Incoming

import (
	"github.com/griesbacher/Yarbes/Config"
	"github.com/influxdb/influxdb/client/v2"
)

//LogServerHTTPInterface is HTTP interface which offers log access
type LogServerHTTPInterface struct {
	*HTTPInterface
}

//NewLogServerHTTPInterface creates a new LogServerHTTPInterface
func NewLogServerHTTPInterface(influxClient client.Client) *LogServerHTTPInterface {
	httpI := NewHTTPInterface(Config.GetServerConfig().LogServer.HTTPInterface)
	logHandler := LogServerHTTPHandler{influxClient}
	ruleHTTP := &LogServerHTTPInterface{HTTPInterface: httpI}
	ruleHTTP.HTTPInterface.PublishHandler("/logs", logHandler.LogView)
	ruleHTTP.HTTPInterface.PublishHandler("/resend", logHandler.ResendEvent)
	return ruleHTTP
}
