package ConditionParser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
)

type (
	dataStore struct {
		data          interface{}
		eventMetadata map[string]interface{}
		stack         []ast.Node
		new_stack     []ast.Node
		ignoreNextX   int
		err           error
		maxRecursion  int
		debug         bool
	}

	Lparen struct {
		pos token.Pos
	}

	Rparen struct {
		pos token.Pos
	}
)

func (p Lparen) Pos() token.Pos {
	return p.pos
}
func (p Rparen) Pos() token.Pos {
	return p.pos
}
func (p Lparen) End() token.Pos {
	return p.pos
}
func (p Rparen) End() token.Pos {
	return p.pos
}
func (p Lparen) String() string {
	return "("
}
func (p Rparen) String() string {
	return ")"
}

func NewDataStore(jsonData interface{}, eventMetadata map[string]interface{}, size int, debug bool) *dataStore {
	return &dataStore{
		data:          jsonData,
		eventMetadata: eventMetadata,
		stack:         []ast.Node{},
		ignoreNextX:   0,

		new_stack:    make([]ast.Node, size),
		maxRecursion: 10,
		debug:        debug,
	}
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

func (d *dataStore) returnResult() (bool, error) {
	d.new_stack = cleanupStack(d.new_stack)
	if d.debug {
		fmt.Println("*********")
		printNodes(d.new_stack)
	}
	boolStack, err := d.evaluateBasicLiterals()
	if err != nil {
		return false, err
	}
	result, err := d.evaluateBooleanStack(boolStack)
	if d.debug {
		fmt.Println("*********")
	}
	return result, err
}

func cleanupStack(stack []ast.Node) []ast.Node {
	result := []ast.Node{}
	for _, n := range stack {
		if n != nil {
			if value, ok := n.(*ast.BasicLit); ok && value == nil {
				continue
			}
			result = append(result, n)
		}
	}
	return result
}

func (d *dataStore) evaluateBasicLiterals() ([]interface{}, error) {
	switch len(d.new_stack) {
	case 1:
		lit, litOk := d.new_stack[0].(*ast.BasicLit)
		if litOk {
			if lit == nil {
				return []interface{}{}, errors.New("Got just one element nil literal on the stack ")
			} else {
				return []interface{}{true}, nil
			}
		} else {
			return []interface{}{}, errors.New("Got just one element on the stack which is no literal")
		}
	case 0, 2:
		return []interface{}{}, fmt.Errorf("The stack got %d elements, which is not allowed", len(d.new_stack))
	}

	result := []interface{}{}
	dummyVisitor := conditionVisitor{}
	for i := 0; i < len(d.new_stack)-2; i++ {
		left := d.new_stack[i]
		op := d.new_stack[i+1]
		right := d.new_stack[i+2]

		leftLit, leftOk := left.(*ast.BasicLit)
		opExpr, opOk := op.(*ast.BinaryExpr)
		rightLit, rightOk := right.(*ast.BasicLit)

		if leftOk && rightOk && opOk {
			if b, err := dummyVisitor.compareBasicLit(leftLit, rightLit, opExpr); err == nil {
				result = append(result, b)
				i += 2
			}
		} else {
			result = append(result, left)
		}
	}
	_, lastOk := d.new_stack[len(d.new_stack)-1].(*Rparen)
	_, lastBeforeOk := d.new_stack[len(d.new_stack)-2].(*Rparen)
	if lastOk {
		result = append(result, d.new_stack[len(d.new_stack)-1])
	}
	if lastBeforeOk {
		result = append(result, d.new_stack[len(d.new_stack)-2])
	}

	return result, nil
}
func cleanupBoolStack(stack []interface{}) []interface{} {
	result := []interface{}{}
	for _, n := range stack {
		if n != nil {
			result = append(result, n)
		}
	}
	return result
}
func printBoolStack(stack []interface{}) {
	fmt.Print("--> ")
	for _, v := range stack {
		switch v.(type) {
		case bool:
			fmt.Print(v, ", ")
		case nil:
			fmt.Print("nil, ")
		default:
			printNode(v.(ast.Node), ", ")
		}
	}
	fmt.Println()
}

func (d dataStore) evaluateBooleanStack(stack []interface{}) (bool, error) {
	if len(stack) == 0 {
		return false, errors.New("The booleanStack is empty")
	}
	return d.evaluateBooleanStackRecursive(stack, 0, d.maxRecursion)
}

func (d dataStore) evaluateBooleanStackRecursive(stack []interface{}, level, maxLevel int) (bool, error) {
	level++
	if level > maxLevel {
		return false, errors.New("The max recusion limit has been reached")
	}
	if d.debug {
		printBoolStack(stack)
	}
	cleanUpTime := false
	didSomethingHappen := false
	//search for bool operator bool
	for i := 0; i < len(stack)-2; i++ {
		left := stack[i]
		op := stack[i+1]
		right := stack[i+2]

		leftBool, leftOk := left.(bool)
		opExpr, opOk := op.(*ast.BinaryExpr)
		rightBool, rightOk := right.(bool)
		if leftOk && rightOk && opOk {
			if b, err := evaluateBoolean(leftBool, rightBool, opExpr); err == nil {
				stack[i] = b
				stack[i+1] = nil
				stack[i+2] = nil
				i += 2
				cleanUpTime, didSomethingHappen = true, true
			}
		}
	}

	if cleanUpTime {
		stack = cleanupBoolStack(stack)
		cleanUpTime = false
	}

	// search for ( bool )
	for i := 0; i < len(stack)-2; i++ {
		left := stack[i]
		boolean := stack[i+1]
		right := stack[i+2]

		_, leftOk := left.(*Lparen)
		booleanValue, booleanOk := boolean.(bool)
		_, rightOk := right.(*Rparen)
		if leftOk && rightOk && booleanOk {
			stack[i] = booleanValue
			stack[i+1] = nil
			stack[i+2] = nil
			i += 2
			cleanUpTime, didSomethingHappen = true, true
		}
	}
	if cleanUpTime {
		stack = cleanupBoolStack(stack)
		cleanUpTime = false
	}
	if len(stack) != 1 {
		if !didSomethingHappen {
			return false, errors.New("Could not replace anything on the booleanStack")
		}
		return d.evaluateBooleanStackRecursive(stack, level, maxLevel)
	} else {
		if result, ok := stack[0].(bool); ok {
			return result, nil
		} else {
			return ok, fmt.Errorf("The last value is not a bool: %s", fmt.Sprint(stack[0]))
		}
	}
}
