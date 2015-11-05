package ConditionParser

import (
	"testing"
	"go/ast"
	"go/token"
)

func TestHandleNode(t *testing.T) {
	t.Parallel()
	v := conditionVisitor{store:&dataStore{result:[]bool{}, stack:[]ast.Node{&ast.IndexExpr{}, &ast.Ident{}}}}
	if v.handleNode(&ast.BasicLit{}); v.store.err == nil {
		t.Error("This should not happen")
	}
	if v.handleNode(&ast.Comment{}); v.store.err == nil {
		t.Error("This should not happen")
	}
}

func TestCompareBasicLit(t *testing.T) {
	t.Parallel()
	v := conditionVisitor{store:&dataStore{}}
	if v.compareBasicLit(&ast.BasicLit{Kind:token.COMMENT}, &ast.BasicLit{Kind:token.COMMENT}, "=="); v.store.err == nil {
		t.Error("compareBasicLit should not compare these tokens")
	}
}