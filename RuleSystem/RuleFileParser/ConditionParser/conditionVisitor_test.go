package ConditionParser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"testing"
)

func TestCompareBasicLit(t *testing.T) {
	t.Parallel()
	v := conditionVisitor{store: &dataStore{}}
	if _, err := v.compareBasicLit(&ast.BasicLit{Kind: token.COMMENT}, &ast.BasicLit{Kind: token.COMMENT}, &ast.BinaryExpr{Op: token.EQL}); err == nil {
		t.Error("compareBasicLit should not compare these tokens")
	}
	if _, err := v.compareBasicLitNumber(&ast.BasicLit{Value: "a"}, &ast.BasicLit{Value: "1"}, "=="); err == nil {
		t.Error("compareBasicLit should not work with not numbers")
	}
	if _, err := v.compareBasicLitNumber(&ast.BasicLit{Value: "1"}, &ast.BasicLit{Value: "a"}, "=="); err == nil {
		t.Error("compareBasicLit should not work with not numbers")
	}
	if result, err := v.compareBasicLit(&ast.BasicLit{Kind: token.INT, Value: "1"}, &ast.BasicLit{Kind: token.INT, Value: "1"}, &ast.BinaryExpr{Op: token.EQL}); !result {
		t.Error("This should return true")
	} else if err != nil {
		t.Error("compareBasicLit should not return an error", err)
	}
	if result, err := v.compareBasicLit(&ast.BasicLit{Kind: token.INT, Value: "1"}, &ast.BasicLit{Kind: token.INT, Value: "2"}, &ast.BinaryExpr{Op: token.EQL}); result {
		t.Error("This should return false")
	} else if err != nil {
		t.Error("compareBasicLit should not return an error", err)
	}
	if result, err := v.compareBasicLit(&ast.BasicLit{Kind: token.STRING, Value: "1"}, &ast.BasicLit{Kind: token.STRING, Value: "1"}, &ast.BinaryExpr{Op: token.EQL}); !result {
		t.Error("This should return true")
	} else if err != nil {
		t.Error("compareBasicLit should not return an error", err)
	}
	if result, err := v.compareBasicLit(&ast.BasicLit{Kind: token.STRING, Value: "1"}, &ast.BasicLit{Kind: token.STRING, Value: "2"}, &ast.BinaryExpr{Op: token.EQL}); result {
		t.Error("This should return false")
	} else if err != nil {
		t.Error("compareBasicLit should not return an error", err)
	}
}

func TestGenBasicLitFromIndexExpr(t *testing.T) {
	v := conditionVisitor{store: &dataStore{data: map[int]interface{}{1: 1}, stack: []ast.Node{&ast.IndexExpr{Index: &ast.BasicLit{Kind: token.STRING}}}}}
	if _, err := v.genBasicLitFromIndexExpr(&ast.Ident{Name: "_"}); err == nil {
		t.Error("TestGenBasicLitFromIndexExpr expects string but got number")
	}

	v.store.stack = []ast.Node{&ast.IndexExpr{Index: &ast.BasicLit{Kind: token.INT, Value: "a"}}}
	if _, err := v.genBasicLitFromIndexExpr(&ast.Ident{Name: "_"}); err == nil {
		t.Error("TestGenBasicLitFromIndexExpr expects number but got string")
	}

	v = conditionVisitor{store: &dataStore{data: map[int]interface{}{1: 1}, stack: []ast.Node{&ast.IndexExpr{Index: &ast.BasicLit{Kind: token.INT, Value: "1"}}}}}
	if _, err := v.genBasicLitFromIndexExpr(&ast.Ident{Name: "_"}); err != nil {
		t.Error("TestGenBasicLitFromIndexExpr got int and should return int")
	}

	v.store.data = map[float32]interface{}{1.0: 1}
	if _, err := v.genBasicLitFromIndexExpr(&ast.Ident{Name: "_"}); err == nil {
		t.Error("TestGenBasicLitFromIndexExpr got int but no store")
	}
}

var EvaluateBooleanData = []struct {
	input1 bool
	input2 bool
	input3 token.Token
	output bool
	err    error
}{
	{true, true, token.LAND, true, nil},
	{false, true, token.LAND, false, nil},
	{false, false, token.LAND, false, nil},
	{true, false, token.LAND, false, nil},
	{true, true, token.LOR, true, nil},
	{true, false, token.LOR, true, nil},
	{false, true, token.LOR, true, nil},
	{false, false, token.LOR, false, nil},
	{false, false, token.ADD, false, errors.New("not supported")},
}

func TestEvaluateBoolean(t *testing.T) {
	for i, data := range EvaluateBooleanData {
		actual, err := evaluateBoolean(data.input1, data.input2, &ast.BinaryExpr{Op: data.input3})
		if err == nil && data.err == nil {
			if actual != data.output {
				t.Errorf("ParseStringData #%d: expected: %t, actual: %t.", i, data.output, actual)
			}
		} else if err != nil && data.err != nil {
			//fine
		} else {
			t.Errorf("The errors do not match for '%d'. Expexted: %s Got: %s", i, fmt.Sprint(data.err), fmt.Sprint(err))
		}
	}
}
