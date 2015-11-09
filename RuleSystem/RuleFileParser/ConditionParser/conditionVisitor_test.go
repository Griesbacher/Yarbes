package ConditionParser

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestHandleNode(t *testing.T) {
	t.Parallel()
	v := conditionVisitor{store: &dataStore{result: []bool{}, stack: []ast.Node{&ast.IndexExpr{}, &ast.Ident{}}}}
	if v.handleNode(&ast.BasicLit{}); v.store.err == nil {
		t.Error("This should not happen")
	}
	if v.handleNode(&ast.Comment{}); v.store.err == nil {
		t.Error("This should not happen")
	}
}

func TestCompareBasicLit(t *testing.T) {
	t.Parallel()
	v := conditionVisitor{store: &dataStore{}}
	if v.compareBasicLit(&ast.BasicLit{Kind: token.COMMENT}, &ast.BasicLit{Kind: token.COMMENT}, "=="); v.store.err == nil {
		t.Error("compareBasicLit should not compare these tokens")
	}
	if v.compareBasicLitNumber(&ast.BasicLit{Value: "a"}, &ast.BasicLit{Value: "1"}, "=="); v.store.err == nil {
		t.Error("compareBasicLit should not work with not numbers")
	}
	if v.compareBasicLitNumber(&ast.BasicLit{Value: "1"}, &ast.BasicLit{Value: "a"}, "=="); v.store.err == nil {
		t.Error("compareBasicLit should not work with not numbers")
	}
}

func TestGenBasicLitFromIndexExpr(t *testing.T) {
	v := conditionVisitor{store: &dataStore{data: map[int]interface{}{1: 1}, stack: []ast.Node{&ast.IndexExpr{Index: &ast.BasicLit{Kind: token.STRING}}}}}
	if v.genBasicLitFromIndexExpr(&ast.Ident{Name: "_"}); v.store.err == nil {
		t.Error("TestGenBasicLitFromIndexExpr expects string but got number")
	}

	v.store.stack = []ast.Node{&ast.IndexExpr{Index: &ast.BasicLit{Kind: token.INT, Value: "a"}}}
	if v.genBasicLitFromIndexExpr(&ast.Ident{Name: "_"}); v.store.err == nil {
		t.Error("TestGenBasicLitFromIndexExpr expects number but got string")
	}

	v = conditionVisitor{store: &dataStore{data: map[int]interface{}{1: 1}, stack: []ast.Node{&ast.IndexExpr{Index: &ast.BasicLit{Kind: token.INT, Value: "1"}}}}}
	if v.genBasicLitFromIndexExpr(&ast.Ident{Name: "_"}); v.store.err != nil {
		t.Error("TestGenBasicLitFromIndexExpr got int and should return int")
	}

	v.store.data = map[float32]interface{}{1.0: 1}
	if v.genBasicLitFromIndexExpr(&ast.Ident{Name: "_"}); v.store.err == nil {
		t.Error("TestGenBasicLitFromIndexExpr got int but no store")
	}
}
