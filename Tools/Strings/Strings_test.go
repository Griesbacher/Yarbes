package Strings

import (
	"reflect"
	"testing"
)

var SearchData = []struct {
	hay    []string
	needle string
	output int
}{
	{[]string{"a", "b", "c"}, "a", 0},
	{[]string{"a", "b", "c"}, "b", 1},
	{[]string{"a", "b", "c"}, "c", 2},
	{[]string{"a", "b", "c"}, "d", -1},
}

func TestContains(t *testing.T) {
	t.Parallel()
	for _, data := range SearchData {
		out := Contains(data.hay, data.needle)
		if out && data.output < 0 {
			t.Errorf("Slice: %s, StringToSearch: %s, Expected: %d, Got: %t", data.hay, data.needle, data.output, out)
		}
	}
}

func TestFormatJSON(t *testing.T) {
	t.Parallel()
	if FormatJSON(`"a":"b"`) != "" {
		t.Errorf("Expected empty string because JSON is not valid: %s", `"a":"b"`)
	}
	expected := `{
  "a": "b"
}`
	if FormatJSON(`{"a":"b"}`) != expected {
		t.Errorf("Expected string because JSON is valid: %s, got: %s", `{"a":"b"}`, FormatJSON(`{"a":"b"}`))
	}
}

func TestUnmarshalJSONEvent(t *testing.T) {
	t.Parallel()
	result := UnmarshalJSONEvent(`{"a":"b"}`)
	expected := map[string]interface{}{"a": "b"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected an object, because the string is valid json. Got: %s, Expected: %s", result, expected)
	}
}
