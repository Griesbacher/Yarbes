package Strings

import (
	"bytes"
	"encoding/json"
)

//IndexOf returns the index of a string in a string slice or -1 if not found
func IndexOf(hay []string, needle string) int {
	for i, a := range hay {
		if a == needle {
			return i
		}
	}
	return -1
}

//Contains returns true if the given string is within the array
func Contains(hay []string, needle string) bool {
	if IndexOf(hay, needle) < 0 {
		return false
	}
	return true
}

//FormatJSON formats the given string in pretty json
func FormatJSON(jsonString string) string {
	var out bytes.Buffer
	if json.Indent(&out, []byte(jsonString), "", "  ") != nil {
		return ""
	}
	return string(out.Bytes())
}

//UnmarshalJSONEvent expects a jsonstring and tries to create a object from it
func UnmarshalJSONEvent(jsonString string) map[string]interface{} {
	var jsonInterface interface{}
	json.Unmarshal([]byte(jsonString), &jsonInterface)
	return jsonInterface.(map[string]interface{})
}
