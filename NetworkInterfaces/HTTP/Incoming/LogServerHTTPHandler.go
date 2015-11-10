package Incoming

import (
	"fmt"
	"github.com/abbot/go-http-auth"
	"github.com/griesbacher/Yarbes/Config"
	"github.com/griesbacher/Yarbes/Logging/LogServer"
	"github.com/griesbacher/Yarbes/NetworkInterfaces/RPC/Outgoing"
	"github.com/griesbacher/Yarbes/Tools/Influx"
	"github.com/griesbacher/Yarbes/Tools/Strings"
	"github.com/influxdb/influxdb/client/v2"
	"io"
	"net/http"
	"net/url"
)

//LogServerHTTPHandler is a HTTP handler which displays LogMessages
type LogServerHTTPHandler struct {
	influxClient client.Client
}

//GETEVENTKEY is the key within the get params for the event
const GETEVENTKEY = "event"

//GETADDRESSKEY is the key within the get params for the server address
const GETADDRESSKEY = "address"

//LogView displays the basic logs
func (handler LogServerHTTPHandler) LogView(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	result, err := Influx.QueryDB(handler.influxClient, fmt.Sprintf("select time, event, msg, serveritry, source from %s", LogServer.TableName), Config.GetServerConfig().LogServer.InfluxDatabase)
	if err != nil {
		panic(err)
	}

	table := `<html><body><table style="width:100%"><tr>`
	for _, column := range result[0].Series[0].Columns {
		table += fmt.Sprintf(`<th>%s</th>`, column)
	}
	table += "</tr>"
	for _, row := range result[0].Series[0].Values {
		table += "<tr>"
		for i, field := range row {
			var output = fmt.Sprint(field)
			if i == 1 {
				u := r.URL
				u.Path = "/resend"
				parameters := url.Values{}
				parameters.Add(GETEVENTKEY, output)
				u.RawQuery = parameters.Encode()
				output = fmt.Sprintf(`<pre><a href="%s">%s</a></pre>`, u, Strings.FormatJSON(output))
			}
			table += fmt.Sprintf(`<td>%s</td>`, output)
		}
		table += "</tr>"
	}
	table += "</table></body></html>"

	io.WriteString(w, table)
}

//ResendEvent sends the selected event to a Rulesystem
func (handler LogServerHTTPHandler) ResendEvent(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	event := r.FormValue(GETEVENTKEY)
	address := r.FormValue(GETADDRESSKEY)
	page := ""
	if address == "" {
		page = fmt.Sprintf(
			`<form action="%s">
			Where should the event sent be to?
			<input name="%s">
			<input type="hidden" name='%s' value='%s'>
			<input type="submit" value="send">
		</form>
		<pre>
		%s
		</pre>`, r.URL, GETADDRESSKEY, GETEVENTKEY, event, Strings.FormatJSON(event))
	} else {
		eventRPC := Outgoing.NewRPCInterface(address)
		err := eventRPC.Connect()
		if err != nil {
			page = err.Error()
		} else {
			err = eventRPC.CreateEvent([]byte(event))
			if err != nil {
				page = err.Error()
			} else {
				page = "Event was sent"
			}
			eventRPC.Disconnect()
		}

	}
	io.WriteString(w, "<html><body>")
	io.WriteString(w, page)
	io.WriteString(w, "</body></html>")
}
