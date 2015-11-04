package Strings
import (
	"bytes"
	"encoding/json"
)

//Contains returns true if the given string is within the array
func Contains(hay []string, needle string) bool {
	for _, a := range hay {
		if a == needle {
			return true
		}
	}
	return false
}

//IndexOf returns the index of a string in a string slice or -1 if not found
func IndexOf(hay []string, needle string) int {
	for i, a := range hay {
		if a == needle {
			return i
		}
	}
	return -1
}

//FormatJSON formats the given string in pretty json
func FormatJSON(jsonString string) string {
	var out bytes.Buffer
	if json.Indent(&out, []byte(jsonString), "", "  ") != nil {
		return ""
	}
	return string(out.Bytes())
}