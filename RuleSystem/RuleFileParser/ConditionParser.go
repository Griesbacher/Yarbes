package RuleFileParser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type dataStore struct {
	data        interface{}
	stack       []ast.Node
	result      []bool
	ignoreNextX int
	err         error
}

func (d *dataStore) appendToStack(node ast.Node) {
	d.stack = append(d.stack, node)
}

func (d *dataStore) popFromStack() ast.Node {
	var last ast.Node
	if len(d.stack) > 0 {
		last, d.stack = d.stack[len(d.stack)-1], d.stack[:len(d.stack)-1]
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
	switch d.stack[len(d.stack)-1].(*ast.BinaryExpr).Op.String() {
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
	}

	switch lastToken := d.stack[len(d.stack)-1].(type) {
	case *ast.BasicLit:
		if lastToken.Kind == token.ILLEGAL {
			return false
		}
	}
	d.err = errors.New("Result is empty")
	return false
}

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

	visitor := myVisitor{p.Debug, data}
	ast.Walk(visitor, tree)

	return visitor.store.returnResult(), data.err
}

type myVisitor struct {
	debug bool
	store *dataStore
}

const charsToTimInStrings = "\"`"

func (v myVisitor) Visit(node ast.Node) ast.Visitor {
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

func (v myVisitor) compareBasicLit(lit1, lit2 *ast.BasicLit, op string) bool {
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

func (v myVisitor) genBasicLitFromIndexExpr(ident *ast.Ident) *ast.BasicLit {
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

func (v myVisitor) printNodes(node []ast.Node) {
	fmt.Print(">> ")
	for _, v := range v.store.stack {
		printNode(v, ", ")
	}
	fmt.Print("\n")
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
