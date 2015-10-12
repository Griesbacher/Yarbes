package RuleFileParser

import (
	"encoding/json"
	"testing"
)

var ParseStringData = []struct {
	input  string
	output bool
}{
	{"1==1", true},
	{"1==2", false},
	{"1!=2", true},
	{"1!=1", false},
	{"1>=1", true},
	{"2>=1", true},
	{"1>=2", false},
	{"1<=1", true},
	{"1<=2", true},
	{"2<=1", false},
	{"2>1", true},
	{"2>2", false},
	{"1<2", true},
	{"2<2", false},
	{`"a" == "a"`, true},
	{`"a" == "b"`, false},
	{`"a" != "b"`, true},
	{`"a" != "a"`, false},
	{`"abba" &^ "a.+a"`, true},
	{`data["k1"] &^ "v\\d"`, true},
	{`data["k1"] == "v1"`, true},
	{`data["k2"] == 10`, true},
	{`data["k3"][0] == "v4"`, true},
	{`data["k3"][1] == 12.5`, true},
	{`data["k3"][2]["k11"] == "v11"`, true},
	{`_["k3"][2]["k22"] == "v22"`, true},
	{`data["k1"]`, true},
	{`data["zzz"]`, false},
}

func TestParseString(t *testing.T) {
	t.Parallel()
	b := []byte(`{
   "k1" : "v1",
   "k2" : 10,
   "k3" : ["v4",12.3,{"k11" : "v11", "k22" : "v22"}]
	}`)
	var jsonData interface{}
	err := json.Unmarshal(b, &jsonData)
	if err != nil {
		panic(err)
	}

	parser := ConditionParser{}
	for _, data := range ParseStringData {
		actual, err := parser.ParseString(data.input, jsonData)
		if actual != data.output && err != nil {
			t.Errorf("ParseStringData(%s): expected: %t, actual: %t. Err: %s", data.input, data.output, actual, err)
		}
	}

}
func TestParseStringChannel(t *testing.T) {
	t.Parallel()
	b := []byte(`{
   "k1" : "v1",
   "k2" : 10,
   "k3" : ["v4",12.3,{"k11" : "v11", "k22" : "v22"}]
	}`)
	var jsonData interface{}
	err := json.Unmarshal(b, &jsonData)
	if err != nil {
		panic(err)
	}

	parser := ConditionParser{}
	var outputData []chan bool
	var outputErr []chan error
	for _, data := range ParseStringData {
		c := make(chan bool, 1)
		e := make(chan error, 1)
		outputData = append(outputData, c)
		outputErr = append(outputErr, e)
		go parser.ParseStringChannel(data.input, jsonData, c, e)
	}
	for i, data := range ParseStringData {
		actual := <-outputData[i]
		err := <-outputErr[i]
		if actual != data.output && err != nil {
			t.Errorf("ParseStringData(%s): expected: %t, actual: %t. Err: %s", data.input, data.output, actual, err)
		}
	}
}
