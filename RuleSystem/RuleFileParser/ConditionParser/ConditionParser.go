package ConditionParser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
)

//ConditionParser parses a given string which should contain a go like condition, conditions can reference on a given JSON struct
type ConditionParser struct {
	//Prints debug output on std
	Debug bool
}

//ParseStringChannel parses the string and communicates through channels, if there is an error the result is irrelevant
func (p ConditionParser) ParseStringChannel(condition string, jsonData interface{}, output chan bool, errors chan error) {
	result, err := p.ParseString(condition, jsonData)
	for i := 0; i < 2; i++ {
		select {
		case output <- result:
		case errors <- err:
		}
	}
}

//ParseString parses the string and JSON object, if there is an error the result is irrelevant
func (p ConditionParser) ParseString(condition string, jsonData interface{}) (bool, error) {
	data := &dataStore{data: jsonData, stack: []ast.Node{}, result: []bool{}, ignoreNextX: 0}
	tree, err := parser.ParseExpr(condition)
	if err != nil {
		return false, err
	}
	if p.Debug {
		ast.Print(token.NewFileSet(), tree)
	}

	visitor := conditionVisitor{p.Debug, data}
	ast.Walk(visitor, tree)

	return visitor.store.returnResult(), data.err
}

func printNode(node ast.Node, appendix string) {
	switch n := node.(type) {
	case *ast.BasicLit:
		fmt.Print(n.Value, appendix)
	case *ast.BinaryExpr:
		fmt.Print(n.Op, appendix)
	case *ast.ParenExpr:
	case *ast.IndexExpr:
		fmt.Print(n.Index, appendix)
	case *ast.Ident:
		fmt.Print(n.Name, appendix)
	default:
		if n != nil {
			fmt.Print("ERROR - ", reflect.TypeOf(node), ", ")
		}
	}
}
