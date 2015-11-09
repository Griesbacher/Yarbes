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
var ErrElementNotFound = errors.New("Element not found")

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
		fmt.Print(">")
		printNode(node, " - ")
		fmt.Println(reflect.TypeOf(node))
		v.printNodes(v.store.stack)
		fmt.Println(v.store.result)
		fmt.Println(v.store.ignoreNextX)
	}
	if len(v.store.stack) == 0 {
		v.store.appendToStack(node)
	} else {
		v.handleNode(node)
	}

	return ast.Visitor(v)
}

func (v conditionVisitor) handleNode(node ast.Node) {
	switch n := node.(type) { //Current element
	case *ast.BasicLit:
		switch nstack := v.store.stack[len(v.store.stack)-1].(type) { //Last element
		case *ast.BasicLit:
			v.store.popFromStack()
			op := v.store.popFromStack().(*ast.BinaryExpr).Op.String()
			v.store.appendToResult(v.compareBasicLit(nstack, n, op))
		case *ast.IndexExpr, *ast.Ident:
			v.store.err = errors.New("should not happen")
		default:
			v.store.appendToStack(node)
		}
	case *ast.BinaryExpr, *ast.IndexExpr:
		v.store.appendToStack(node)
	case *ast.Ident:
		if indexExpr := v.genBasicLitFromIndexExpr(n); indexExpr != nil {
			if len(v.store.stack) > 0 {
				/*switch nstack := v.store.stack[len(v.store.stack) - 1].(type) { //Last element
				case *ast.BasicLit:
					v.store.popFromStack()
					op := v.store.popFromStack().(*ast.BinaryExpr).Op.String()
					v.store.appendToResult(v.compareBasicLit(nstack, indexExpr, op))
				default:
					v.store.appendToStack(indexExpr)
				}*/ //Useless?!
				v.store.appendToStack(indexExpr)
			} else {
				v.store.appendToResult(true) //Found index
			}
		} else {
			v.store.appendToStack(&ast.BasicLit{ValuePos: token.NoPos, Kind: token.ILLEGAL, Value: "X"})
		}
	case *ast.ParenExpr, nil:
	//Allowed but not used
	default:
		//Not allowed
		//panic("Token not allowed!")
		v.store.err = errors.New("Token not allowed!")
	}
}

func (v conditionVisitor) compareBasicLit(lit1, lit2 *ast.BasicLit, op string) bool {
	if lit1.Kind != lit2.Kind {
		v.store.err = fmt.Errorf("Types don't match: %s != %s. Values: %s, %s", lit1.Kind, lit2.Kind, lit1.Value, lit2.Value)
		return false
	}
	switch lit1.Kind {
	case token.STRING:
		return v.compareBasicLitString(lit1, lit2, op)
	case token.INT, token.FLOAT:
		return v.compareBasicLitNumber(lit1, lit2, op)
	default:
		v.store.err = errors.New("An unkown token appeard")
		return false
	}
}

func (v conditionVisitor) compareBasicLitString(lit1, lit2 *ast.BasicLit, op string) bool {
	value1 := strings.Trim(lit1.Value, charsToTimInStrings)
	value2 := strings.Trim(lit2.Value, charsToTimInStrings)
	switch op {
	case "==":
		return value1 == value2
	case "!=":
		return value1 != value2
	case "&^":
		value2 = strings.Replace(value2, "\\\\", "\\", -1)
		matched, err := regexp.MatchString(value2, value1)
		if err != nil {
			v.store.err = err
			return false
		}
		return matched
	default:
		v.store.err = errors.New("used unsupported operator")
		return false
	}
}

func (v conditionVisitor) compareBasicLitNumber(lit1, lit2 *ast.BasicLit, op string) bool {
	value1, err := strconv.ParseFloat(lit1.Value, 32)
	if err != nil {
		v.store.err = err
		return false
	}
	value2, err := strconv.ParseFloat(lit2.Value, 32)
	if err != nil {
		v.store.err = err
		return false
	}
	switch op {
	case "==":
		return value1 == value2
	case "!=":
		return value1 != value2
	case ">=":
		return value1 >= value2
	case "<=":
		return value1 <= value2
	case ">":
		return value1 > value2
	case "<":
		return value1 < value2
	default:
		v.store.err = errors.New("used unsupported operator")
		return false
	}
}

func (v conditionVisitor) genBasicLitFromIndexExpr(ident *ast.Ident) *ast.BasicLit {
	currentLevel, err := v.searchForData(ident.Name)
	if err != nil {
		v.store.err = err
		return nil
	}
	switch value := currentLevel.(type) {
	case string:
		return &ast.BasicLit{ValuePos: token.NoPos, Kind: token.STRING, Value: "\"" + value + "\""}
	case int, float32, float64:
		asString := fmt.Sprint(value)
		if strings.Contains(asString, ".") {
			return &ast.BasicLit{ValuePos: token.NoPos, Kind: token.FLOAT, Value: asString}
		}

		return &ast.BasicLit{ValuePos: token.NoPos, Kind: token.INT, Value: asString}
	case nil:
		v.store.err = ErrElementNotFound
		return nil
	default:
		v.store.err = fmt.Errorf("No suitable type found... %s", reflect.TypeOf(currentLevel))
		return nil
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

func (v conditionVisitor) printNodes(node []ast.Node) {
	fmt.Print(">> ")
	for _, v := range v.store.stack {
		printNode(v, ", ")
	}
	fmt.Print("\n")
}
