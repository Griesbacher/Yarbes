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
		switch n := node.(type) { //Current element
		case *ast.BasicLit:
			switch nstack := v.store.stack[len(v.store.stack)-1].(type) { //Last element
			case *ast.BasicLit:
				v.store.popFromStack()
				op := v.store.popFromStack().(*ast.BinaryExpr).Op.String()
				v.store.appendToResult(v.compareBasicLit(nstack, n, op))
			case *ast.IndexExpr, *ast.Ident:
				panic("should not happen")
			default:
				v.store.appendToStack(node)
			}
		case *ast.BinaryExpr, *ast.IndexExpr:
			v.store.appendToStack(node)
		case *ast.Ident:
			if indexExpr := v.genBasicLitFromIndexExpr(n); indexExpr != nil {
				if len(v.store.stack) > 0 {
					switch nstack := v.store.stack[len(v.store.stack)-1].(type) { //Last element
					case *ast.BasicLit:
						v.store.popFromStack()
						op := v.store.popFromStack().(*ast.BinaryExpr).Op.String()
						v.store.appendToResult(v.compareBasicLit(nstack, indexExpr, op))
					default:
						v.store.appendToStack(indexExpr)
					}
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

	return ast.Visitor(v)
}

func (v conditionVisitor) compareBasicLit(lit1, lit2 *ast.BasicLit, op string) bool {
	if lit1.Kind != lit2.Kind {
		v.store.err = fmt.Errorf("Types don't match: %s != %s. Values: %s, %s", lit1.Kind, lit2.Kind, lit1.Value, lit2.Value)
		return false
	}

	switch lit1.Kind {
	case token.STRING:
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
			//panic("used unsupported operator!")
		}
	case token.INT, token.FLOAT:
		value1, _ := strconv.Atoi(lit1.Value)
		value2, _ := strconv.Atoi(lit2.Value)
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
			//panic("used unsupported operator!")
		}
	default:
		v.store.err = errors.New("An unkown token appeard")
		return false
		//panic("An unkown token appeard")
	}
}

func (v conditionVisitor) genBasicLitFromIndexExpr(ident *ast.Ident) *ast.BasicLit {
	var currentLevel interface{}
	currentLevel = v.store.data
	if ident.Name != "" { //TODO: Namen für datenstrucktur überlegen
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
						v.store.err = errors.New("Got string but it's no map")
						return nil
					}
				case token.INT, token.FLOAT: //Should never happen due to json convention
					key, err := strconv.Atoi(lit.Value)
					if err != nil {
						v.store.err = fmt.Errorf("Could not cast string: %s to int", lit.Value)
						return nil
					}
					switch d := currentLevel.(type) {
					case map[int]interface{}:
						currentLevel = d[key]
					case []interface{}:
						currentLevel = d[key]
					default:
						v.store.err = errors.New("Got int but it's no map nor list")
						return nil
					}
				}
			default:
				break
			}
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
	} else {
		v.store.err = fmt.Errorf("Given datastructure name is wrong. Given: %s Expected", ident.Name)
		return nil
	}
}

func (v conditionVisitor) printNodes(node []ast.Node) {
	fmt.Print(">> ")
	for _, v := range v.store.stack {
		printNode(v, ", ")
	}
	fmt.Print("\n")
}
