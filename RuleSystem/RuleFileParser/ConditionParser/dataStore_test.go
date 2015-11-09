package ConditionParser

import (
	"go/ast"
	"go/token"
	"testing"
)

//TestReturnResult is a dummy test because the PrintNode is just for debugging
func TestReturnResult(t *testing.T) {
	t.Parallel()
	d := dataStore{result: []bool{}}
	if d.returnResult() {
		t.Error("Result should be false on empty list")
	}
	d = dataStore{result: []bool{}, stack: []ast.Node{&ast.BasicLit{Kind: token.COMMENT}}}
	if !didThisPanic(d.returnResult) {
		t.Error("returnResult should panic on false Asttype")
	}
	d = dataStore{result: []bool{}, stack: []ast.Node{&ast.Comment{}}}
	if !didThisPanic(d.returnResult) {
		t.Error("returnResult should panic on false Tokentype")
	}
}

func didThisPanic(f func() bool) (result bool) {
	defer func() {
		if rec := recover(); rec != nil {
			result = true
		}
	}()
	f()
	return false
}
