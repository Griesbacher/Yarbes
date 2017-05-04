package ConditionParser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type conditionVisitor struct {
	debug bool
	store *dataStore
}

const charsToTimInStrings = "\"`"

//ErrElementNotFound will be returned if the key is not found in the event
var ErrElementNotFound = errors.New("Key could not be found in the event/metadata")

func (v conditionVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil || v.store.err != nil {
		return ast.Visitor(nil)
	}
	if v.store.ignoreNextX > 0 {
		v.store.ignoreNextX--
		return ast.Visitor(nil)
	}
	if v.debug {
		fmt.Println("-------------------------")
		printNode(node, " - ")
		fmt.Println(reflect.TypeOf(node))
		printNodes(v.store.stack)
		printNodes(v.store.new_stack)
	}
	v.handleNode(node)

	return ast.Visitor(v)
}

func (v conditionVisitor) handleNode(node ast.Node) {
	switch n := node.(type) { //Current element
	case *ast.BasicLit:
		v.store.new_stack[n.ValuePos] = n
	case *ast.BinaryExpr:
		v.store.appendToStack(node)
		v.store.new_stack[n.OpPos] = n
	case *ast.IndexExpr:
		v.store.appendToStack(node)
	case *ast.Ident:
		v.store.new_stack[n.Pos()] = n
		if indexExpr, err := v.genBasicLitFromIndexExpr(n); err == nil && indexExpr != nil {
			v.store.new_stack[n.Pos()] = indexExpr
		} else {
			v.store.err = err
			v.store.appendToStack(&ast.BasicLit{ValuePos: token.NoPos, Kind: token.ILLEGAL, Value: "X"})
		}
	case *ast.ParenExpr:
		v.store.new_stack[n.Lparen] = &Lparen{pos: n.Lparen}
		v.store.new_stack[n.Rparen] = &Rparen{pos: n.Rparen}
	case nil:
	default:
		//Not allowed
		v.store.err = errors.New("Token not allowed!")
	}
}

func (v conditionVisitor) compareBasicLit(lit1, lit2 *ast.BasicLit, op *ast.BinaryExpr) (bool, error) {
	if lit1.Kind != lit2.Kind {
		return false, fmt.Errorf("Types don't match: %s != %s. Values: %s, %s", lit1.Kind, lit2.Kind, lit1.Value, lit2.Value)
	}
	operator := op.Op.String()
	switch lit1.Kind {
	case token.STRING:
		return v.compareBasicLitString(lit1, lit2, operator)
	case token.INT, token.FLOAT:
		return v.compareBasicLitNumber(lit1, lit2, operator)
	default:
		return false, errors.New("An unkown token appeard")
	}
}

func (v conditionVisitor) compareBasicLitString(lit1, lit2 *ast.BasicLit, op string) (bool, error) {
	value1 := strings.Trim(lit1.Value, charsToTimInStrings)
	value2 := strings.Trim(lit2.Value, charsToTimInStrings)
	switch op {
	case "==":
		return value1 == value2, nil
	case "!=":
		return value1 != value2, nil
	case "&^":
		value2 = strings.Replace(value2, "\\\\", "\\", -1)
		matched, err := regexp.MatchString(value2, value1)
		return matched, err
	default:
		return false, errors.New("used unsupported operator")
	}
}

func (v conditionVisitor) compareBasicLitNumber(lit1, lit2 *ast.BasicLit, op string) (bool, error) {
	value1, err := strconv.ParseFloat(lit1.Value, 32)
	if err != nil {
		return false, err
	}
	value2, err := strconv.ParseFloat(lit2.Value, 32)
	if err != nil {
		return false, err
	}
	switch op {
	case "==":
		return value1 == value2, nil
	case "!=":
		return value1 != value2, nil
	case ">=":
		return value1 >= value2, nil
	case "<=":
		return value1 <= value2, nil
	case ">":
		return value1 > value2, nil
	case "<":
		return value1 < value2, nil
	default:
		return false, errors.New("used unsupported operator")
	}
}

func (v conditionVisitor) genBasicLitFromIndexExpr(ident *ast.Ident) (*ast.BasicLit, error) {
	currentLevel, err := v.searchForData(ident.Name)
	if err != nil {
		return nil, err
	}
	switch value := currentLevel.(type) {
	case string:
		return &ast.BasicLit{ValuePos: token.NoPos, Kind: token.STRING, Value: "\"" + value + "\""}, nil
	case int, float32, float64:
		asString := fmt.Sprint(value)
		if strings.Contains(asString, ".") {
			return &ast.BasicLit{ValuePos: token.NoPos, Kind: token.FLOAT, Value: asString}, nil
		}

		return &ast.BasicLit{ValuePos: token.NoPos, Kind: token.INT, Value: asString}, nil
	case nil:
		return nil, ErrElementNotFound
	default:
		return nil, fmt.Errorf("No suitable type found... %s", reflect.TypeOf(currentLevel))
	}
}

func (v conditionVisitor) searchForData(dataType string) (interface{}, error) {
	var currentLevel interface{}
	switch dataType {
	case "_":
		currentLevel = v.store.data
	case "e":
		currentLevel = v.store.eventMetadata
	default:
		return nil, fmt.Errorf("Given datastructure name is wrong. Given: %s Expected", dataType)
	}
	stackSize := len(v.store.stack) - 1
	for i := stackSize; i >= 0; i-- {
		switch s := v.store.stack[i].(type) {
		case *ast.IndexExpr:
			v.store.ignoreNextX++
			v.store.popFromStack()
			lit := s.Index.(*ast.BasicLit)
			switch lit.Kind {
			case token.STRING:
				key := strings.Trim(lit.Value, charsToTimInStrings)
				switch d := currentLevel.(type) {
				case map[string]interface{}:
					currentLevel = d[key]
				default:
					return nil, errors.New("Got string but it's no map")
				}
			case token.INT, token.FLOAT: //Should never happen due to json convention
				key, err := strconv.Atoi(lit.Value)
				if err != nil {
					return nil, fmt.Errorf("Could not cast string: %s to int", lit.Value)
				}
				switch d := currentLevel.(type) {
				case map[int]interface{}:
					currentLevel = d[key]
				case []interface{}:
					currentLevel = d[key]
				}
			}
		default:
			break
		}
	}
	return currentLevel, nil
}

func evaluateBoolean(a, b bool, op *ast.BinaryExpr) (result bool, err error) {
	switch op.Op.String() {
	case "&&":
		result = a && b
	case "||":
		result = a || b
	default:
		err = fmt.Errorf("this operant is not supported: %s", op.Op.String())
	}
	return
}

func printNodes(node []ast.Node) {
	fmt.Print("stack: ")
	for _, v := range node {
		printNode(v, ", ")
	}
	fmt.Print("\n")
}
