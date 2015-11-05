package ConditionParser

import (
	"errors"
	"go/ast"
	"go/token"
)

type dataStore struct {
	data          interface{}
	eventMetadata map[string]interface{}
	stack         []ast.Node
	result        []bool
	ignoreNextX   int
	err           error
}

func (d *dataStore) appendToStack(node ast.Node) {
	d.stack = append(d.stack, node)
}

func (d *dataStore) popFromStack() ast.Node {
	var last ast.Node
	if len(d.stack) > 0 {
		last, d.stack = d.stack[len(d.stack) - 1], d.stack[:len(d.stack) - 1]
	}
	return last
}

func (d *dataStore) appendToResult(result bool) {
	d.result = append(d.result, result)
	if len(d.result) == 2 {
		d.evaluateResultQueue()
	}
}

func (d *dataStore) evaluateResultQueue() {
	switch d.stack[len(d.stack) - 1].(*ast.BinaryExpr).Op.String() {
	case "&&":
		d.result = []bool{d.result[0] && d.result[1]}
	case "||":
		d.result = []bool{d.result[0] || d.result[1]}
	}
	d.popFromStack()
}

func (d *dataStore) returnResult() bool {
	if len(d.result) == 1 {
		return d.result[0]
	}else if len(d.stack) > 0 {
		switch lastToken := d.stack[len(d.stack) - 1].(type) {
		case *ast.BasicLit:
			if lastToken.Kind == token.ILLEGAL {
				return false
			}
			panic("should not happen")
		default:
			panic("should not happen")
		}
	}
	d.err = errors.New("Result/Stack is empty")
	return false
}
