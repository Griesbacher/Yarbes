package ConditionParser

import (
	"encoding/json"
	"testing"
	"os"
	"errors"
	"go/ast"
)

var ParseStringData = []struct {
	input  string
	output bool
	err    error
}{
	{"1==1", true, nil },
	{"1==2", false, nil },
	{"1!=2", true, nil },
	{"1!=1", false, nil },
	{"1>=1", true, nil },
	{"2>=1", true, nil },
	{"1>=2", false, nil },
	{"1<=1", true, nil },
	{"1<=2", true, nil },
	{"2<=1", false, nil },
	{"2>1", true, nil },
	{"2>2", false, nil },
	{"1<2", true, nil },
	{"2<2", false, nil },
	{`"a" == "a"`, true, nil },
	{`"a" == "b"`, false, nil },
	{`"a" != "b"`, true, nil },
	{`"a" != "a"`, false, nil },
	{`"abba" &^ "a.+a"`, true, nil },
	{`_["k1"] &^ "v\\d"`, true, nil },
	{`_["k1"] == "v1"`, true, nil },
	{`_["k2"] == 10`, true, nil },
	{`_["k3"][0] == "v4"`, true, nil },
	{`_["k3"][1] == 12.5`, true, nil },
	{`_["k3"][2]["k11"] == "v11"`, true, nil },
	{`_["k3"][2]["k22"] == "v22"`, true, nil },
	{`_["k1"]`, true, nil },
	{`_["zzz"]`, false, nil },
	{"`a` == `a`", true, nil },
	{"`a` == `b`", false, nil },
	{"`a` != `b`", true, nil },
	{"`a` != `a`", false, nil },
	{"`abba` &^ `a.+a`", true, nil },
	{"_[`k1`] &^ `v\\d`", true, nil },
	{"_[`k1`] == `v1`", true, nil },
	{"_[`k2`] == 10", true, nil },
	{"_[`k3`][0] == `v4`", true, nil },
	{"_[`k3`][1] == 12.5", true, nil },
	{"_[`k3`][2][`k11`] == `v11`", true, nil },
	{"_[`k3`][2][`k22`] == `v22`", true, nil },
	{"_[`k1`]", true, nil },
	{"_[`zzz`]", false, nil },
	{"e[`executedLines`] == 0", true, nil },
	//{"e[`42`] == true", true, nil },
	{`_["k2"] == 10 && _["k2"] == 10`, true, nil },
	{`_["k2"] == 10 || _["k2"] == 11`, true, nil },
	{`10 == "10"`, false, errors.New("string and int compare") },
	{`10 &^ "\y"`, false, errors.New("not a valid regex") },
	{`1,1 == 1`, false, errors.New("not a valid float") },
	{`1 == 1,1`, false, errors.New("not a valid float") },
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
	currentMetaData := map[string]interface{}{"executedLines": 0, "42" : true}

	parser := ConditionParser{}
	for _, data := range ParseStringData {
		actual, err := parser.ParseString(data.input, jsonData, currentMetaData)
		if actual != data.output && (err != nil && data.err == nil) {
			t.Errorf("ParseStringData(%s): expected: %t, actual: %t. Err: %s", data.input, data.output, actual, err)
		}
	}
}
func TestParseStringChannel(t *testing.T) {
	t.Parallel()
	_, w, _ := os.Pipe()
	os.Stdout = w
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
	currentMetaData := map[string]interface{}{"executedLines": 0}

	parser := ConditionParser{}
	var outputData []chan bool
	var outputErr []chan error
	for _, data := range ParseStringData {
		c := make(chan bool, 1)
		e := make(chan error, 1)
		outputData = append(outputData, c)
		outputErr = append(outputErr, e)
		go parser.ParseStringChannel(data.input, jsonData, currentMetaData, c, e)
	}
	for i, data := range ParseStringData {
		actual := <-outputData[i]
		err := <-outputErr[i]
		if actual != data.output && (err != nil && data.err == nil) {
			t.Errorf("ParseStringData(%s): expected: %t, actual: %t. Err: %s", data.input, data.output, actual, err)
		}
	}
}

//TestPrintNode is a dummy test because the PrintNode is just for debugging
func TestPrintNode(t *testing.T) {
	t.Parallel()
	oldStdout := os.Stdout
	_, writeFile, _ := os.Pipe()
	os.Stdout = writeFile
	nodes := []ast.Node{&ast.BasicLit{}, &ast.BinaryExpr{}, &ast.ParenExpr{}, &ast.IndexExpr{}, &ast.Ident{}, &ast.Comment{}, nil}
	for _, n := range nodes {
		printNode(n, "")
	}
	writeFile.Close()
	os.Stdout = oldStdout
}