package Incoming

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/abbot/go-http-auth"
	"github.com/griesbacher/Yarbes/Config"
	"github.com/griesbacher/Yarbes/Logging/LogServer"
	"github.com/griesbacher/Yarbes/Tools/Influx"
	"github.com/influxdb/influxdb/client/v2"
	"io"
	"net/http"
)

//LogServerHTTPHandler is a HTTP handler which displays LogMessages
type LogServerHTTPHandler struct {
	influxClient client.Client
}

//LogView displays the basic logs
func (handler LogServerHTTPHandler) LogView(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	result, err := Influx.QueryDB(handler.influxClient, fmt.Sprintf("select * from %s", LogServer.TableName), Config.GetServerConfig().LogServer.InfluxDatabase)
	if err != nil {
		panic(err)
	}

	table := `<table style="width:100%"><tr>`
	for _, column := range result[0].Series[0].Columns {
		table += fmt.Sprintf(`<th>%s</th>`, column)
	}
	table += "</tr>"
	for _, row := range result[0].Series[0].Values {
		table += "<tr>"
		for i, field := range row {
			var output = fmt.Sprint(field)
			if i == 1 {
				var out bytes.Buffer
				_ = json.Indent(&out, []byte(output), "", "  ")
				output = fmt.Sprintf(`<pre>%s</pre>`, string(out.Bytes()))
			}
			table += fmt.Sprintf(`<td>%s</td>`, output)
		}
		table += "</tr>"
	}
	table += "</table>"

	io.WriteString(w, "<html><body>")
	io.WriteString(w, table)
	io.WriteString(w, "</body></html>")

}
