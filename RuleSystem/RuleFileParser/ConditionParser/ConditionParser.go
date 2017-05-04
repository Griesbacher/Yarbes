package ConditionParser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

//ConditionParser parses a given string which should contain a go like condition, conditions can reference on a given JSON struct
type ConditionParser struct {
	//Prints debug output on std
	Debug bool
}

//ParseStringChannel parses the string and communicates through channels, if there is an error the result is irrelevant
func (p ConditionParser) ParseStringChannel(condition string, jsonData interface{}, eventMetadata map[string]interface{}, output chan bool, errors chan error) {
	result, err := p.ParseString(condition, jsonData, eventMetadata)
	select {
	case output <- result:
		errors <- err
	case errors <- err:
		output <- result
	}
}

//ParseString parses the string and JSON object, if there is an error the result is irrelevant
func (p ConditionParser) ParseString(condition string, jsonData interface{}, eventMetadata map[string]interface{}) (bool, error) {
	condition = strings.Replace(condition, " ", "", -1)
	dataStore := NewDataStore(jsonData, eventMetadata, len(condition)+1, p.Debug)
	tree, err := parser.ParseExpr(condition)
	if err != nil {
		return false, err
	}
	if p.Debug {
		ast.Print(token.NewFileSet(), tree)
	}

	visitor := conditionVisitor{p.Debug, dataStore}
	ast.Walk(visitor, tree)

	if visitor.store.err != nil {
		return false, visitor.store.err
	}

	return visitor.store.returnResult()
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
	case *Lparen, *Rparen:
		fmt.Print(n, appendix)
	default:
		if n != nil {
			fmt.Print("ERROR - ", reflect.TypeOf(node), ", ")
		} else {
			fmt.Print("nil", appendix)
		}
	}
}
