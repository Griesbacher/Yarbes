package Livestatus

import (
	"fmt"
	"strings"
)

type livestatusResultConverter struct {
	query string
	index map[string]int
}

var host = map[string]string{"hostname": "host_name", "time": "time", "service": "current_service_display_name", "plugin_output": "long_plugin_output"}
var service = map[string]string{"hostname": "host_name", "time": "time", "service": "current_service_display_name", "plugin_output": "long_plugin_output"}

func newLivestatusResultConverter(query string) *livestatusResultConverter {
	index := map[string]int{}
	for _, line := range strings.Split(query, "\n") {
		if len(line) > 7 && line[0:8] == "Columns:" {
			for i, column := range strings.Split(line, " ") {
				if i == 0 {
					continue
				}
				index[column] = i - 1
			}
		}
	}
	fmt.Println(index)
	return &livestatusResultConverter{query: query, index: index}
}

func (c livestatusResultConverter) createObject(result []string) map[string]interface{} {
	fmt.Println(result)
	typ := result[c.index["type"]]
	event := map[string]interface{}{}
	var mappingTable map[string]string
	switch typ {
	case "HOST ALERT", "HOST FLAPPING ALERT":
		mappingTable = host
	case "SERVICE ALERT", "SERVICE FLAPPING ALERT":
		mappingTable = service
	default:
		return event
	}

	for k, v := range mappingTable {
		event[k] = result[c.index[v]]
	}
	fmt.Println(event)
	return event
}
