package ConditionParser

import (
	"go/ast"
	"go/token"
	"testing"
)

//TestReturnResult is a dummy test because the PrintNode is just for debugging
func TestReturnResult(t *testing.T) {
	t.Parallel()
	d := dataStore{}
	if result, err := d.returnResult(); result && err == nil {
		t.Error("Result should be false on empty list")
	}
	d = dataStore{stack: []ast.Node{&ast.BasicLit{Kind: token.COMMENT}}}
	if !didThisPanicOrError(d.returnResult) {
		t.Error("returnResult should panic on false Asttype")
	}
	d = dataStore{stack: []ast.Node{&ast.Comment{}}}
	if !didThisPanicOrError(d.returnResult) {
		t.Error("returnResult should panic on false Tokentype")
	}
}

func didThisPanicOrError(f func() (bool, error)) (result bool) {
	defer func() {
		if rec := recover(); rec != nil {
			result = true
		}
	}()
	_, err := f()
	if err != nil {
		return true
	}
	return false
}
