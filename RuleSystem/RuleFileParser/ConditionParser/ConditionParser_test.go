package ConditionParser

import (
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"os"
	"testing"
)

//go tool cover -html=cover.out

var ParseStringData = []struct {
	input  string
	output bool
	err    error
}{
	{"1==1", true, nil},
	{"1==2", false, nil},
	{"1!=2", true, nil},
	{"1!=1", false, nil},
	{"1>=1", true, nil},
	{"2>=1", true, nil},
	{"1>=2", false, nil},
	{"1<=1", true, nil},
	{"1<=2", true, nil},
	{"2<=1", false, nil},
	{"2>1", true, nil},
	{"2>2", false, nil},
	{"1<2", true, nil},
	{"2<2", false, nil},
	{`"a" == "a"`, true, nil},
	{`"a" == "b"`, false, nil},
	{`"a" != "b"`, true, nil},
	{`"a" != "a"`, false, nil},
	{`"abba" &^ "a.+a"`, true, nil},
	{`_["k1"] &^ "v\\d"`, true, nil},
	{`_["k1"] == "v1"`, true, nil},
	{`_["k2"] == 10`, true, nil},
	{`_["k3"][0] == "v4"`, true, nil},
	{`_["k3"][1] == 12.3`, true, nil},
	{`_["k3"][2]["k11"] == "v11"`, true, nil},
	{`_["k3"][2]["k22"] == "v22"`, true, nil},
	{`_["k1"]`, true, nil},
	{`_["zzz"]`, false, errors.New("Element not found")},
	{"`a` == `a`", true, nil},
	{"`a` == `b`", false, nil},
	{"`a` != `b`", true, nil},
	{"`a` != `a`", false, nil},
	{"`abba` &^ `a.+a`", true, nil},
	{"_[`k1`] &^ `v\\d`", true, nil},
	{"_[`k1`] == `v1`", true, nil},
	{"_[`k2`] == 10", true, nil},
	{"_[`k3`][0] == `v4`", true, nil},
	{"_[`k3`][1] == 12.3", true, nil},
	{"_[`k3`][2][`k11`] == `v11`", true, nil},
	{"_[`k3`][2][`k22`] == `v22`", true, nil},
	{"_[`k1`]", true, nil},
	{"_[`zzz`]", false, errors.New("Element not found")},
	{"e[`executedLines`] == 0", true, nil},
	{`_["k2"] == 10 && _["k2"] == 10`, true, nil},
	{`_["k2"] == 10 || _["k2"] == 11`, true, nil},
	{`(1==1)`, true, nil},
	{`(1==1) && (2==2)`, true, nil},
	{`((1==2) || (1==2 && 1==1))`, false, nil},
	{`10 == "10"`, false, errors.New("string and int compare")},
	{`"10" &^ ")10"`, false, errors.New("not a valid regex")},
	{`1,1 == 1`, false, errors.New("not a valid float")},
	{`1 == 1,1`, false, errors.New("not a valid float")},
	{`->`, false, errors.New("valid go but not allowed")},
	{`"a" < "a"`, false, errors.New("string operator not allowed")},
	{`1 &^ 1`, false, errors.New("number operator not allowed")},
	{`_["1"] == 1`, false, errors.New("invalid key")},
	{`foo["k1"] == 10`, false, errors.New("datastructurename is not allowed")},
	{" (4 == 3 && 10 < 21) || (4 == 2) || (4 == 1 && 10 > 8) ", false, nil},
	{"( (_[`__weekday`] == 3 && _[`__hour`] < 21)  || (_[`__weekday`] == 2) || (_[`__weekday`] == 1 && _[`__hour`] > 8)  )", false, nil},
}

var b = []byte(`{
   "k1" : "v1",
   "k2" : 10,
   "k3" : ["v4",12.3,{"k11" : "v11", "k22" : "v22"}],
   "k4" : true,
   "__weekday": 4,
   "__hour": 10
}`)

func TestParseString(t *testing.T) {
	var jsonData interface{}
	err := json.Unmarshal(b, &jsonData)
	if err != nil {
		panic(err)
	}
	currentMetaData := map[string]interface{}{"executedLines": 0, "42": true}

	parser := ConditionParser{}
	for _, data := range ParseStringData {
		actual, err := parser.ParseString(data.input, jsonData, currentMetaData)
		if err == nil && data.err == nil {
			if actual != data.output {
				t.Errorf("ParseStringData(%s): expected: %t, actual: %t.", data.input, data.output, actual)
				ConditionParser{Debug: true}.ParseString(data.input, jsonData, currentMetaData)
			}
		} else if err != nil && data.err != nil {
			//fine
		} else {
			t.Errorf("The errors do not match for '%s'. Expexted: %s Got: %s", data.input, fmt.Sprint(data.err), fmt.Sprint(err))
		}
	}
}

func TestParseStringChannel(t *testing.T) {
	_, w, _ := os.Pipe()
	os.Stdout = w
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
		if err == nil && data.err == nil {
			if actual != data.output {
				t.Errorf("ParseStringData(%s): expected: %t, actual: %t.", data.input, data.output, actual)
				ConditionParser{Debug: true}.ParseString(data.input, jsonData, currentMetaData)
			}
		} else if err != nil && data.err != nil {
			//fine
		} else {
			t.Errorf("The errors do not match for %s. Expexted: %s Got: %s", data.input, fmt.Sprint(data.err), fmt.Sprint(err))
		}
	}
}

func TestPrintNode(t *testing.T) {
	t.Parallel()
	oldStdout := os.Stdout
	_, writeFile, _ := os.Pipe()
	os.Stdout = writeFile
	nodes := []ast.Node{&ast.BasicLit{}, &ast.BinaryExpr{}, &ast.ParenExpr{}, &ast.IndexExpr{}, &ast.Ident{}, &ast.Comment{}, &Lparen{}, &Rparen{}, nil}
	for _, n := range nodes {
		printNode(n, "")
	}
	writeFile.Close()
	os.Stdout = oldStdout
}

func TestPrintAst(t *testing.T) {
	t.Parallel()
	oldStdout := os.Stdout
	_, writeFile, _ := os.Pipe()
	os.Stdout = writeFile

	parser := ConditionParser{Debug: true}
	var jsonData interface{}
	err := json.Unmarshal(b, &jsonData)
	if err != nil {
		panic(err)
	}
	currentMetaData := map[string]interface{}{"executedLines": 0, "42": true}

	parser.ParseString("1==1", jsonData, currentMetaData)

	writeFile.Close()
	os.Stdout = oldStdout
}

func TestBug(t *testing.T) {
	var jsonData interface{}
	err := json.Unmarshal(b, &jsonData)
	if err != nil {
		panic(err)
	}
	fmt.Println(ConditionParser{Debug: true}.ParseString(`1 == 1`, jsonData, map[string]interface{}{}))
}
