package main

import (
	"fmt"

	"github.com/k0kubun/pp"
)

func main() {
	root := &BinaryExpression{
		Left: &BinaryExpression{
			Left: &IntegerLiteral{
				Value: 10,
			},
			Right: &IntegerLiteral{
				Value: 5,
			},
			Op: "/",
		},
		Right: &BinaryExpression{
			Left: &IntegerLiteral{
				Value: 4,
			},
			Right: &IntegerLiteral{
				Value: 3,
			},
			Op: "*",
		},
		Op: "+",
	}
	calculator := &Calculator{}
	result := runIntepreter(root, calculator)
	fmt.Println(result)
	generator := &Generator{}
	root.Accept(generator)
	pp.Println(instructions)
	fmt.Println(execute())
}

func runIntepreter(root Node, calculator Visitor) int {
	return root.Accept(calculator)
}

type Body struct {
	Statement []Statement
}

type Statement struct {
}

type IntegerLiteral struct {
	Value int
}

func (n *IntegerLiteral) Accept(v Visitor) int {
	return v.VisitInteger(n)
}

type BinaryExpression struct {
	Left  Node
	Right Node
	Op    string
}

func (n *BinaryExpression) Accept(v Visitor) int {
	return v.VisitBinaryExpression(n)
}

type Visitor interface {
	VisitInteger(*IntegerLiteral) int
	VisitBinaryExpression(*BinaryExpression) int
}

type Node interface {
	Accept(v Visitor) int
}

type Calculator struct{}

func (v *Calculator) VisitInteger(n *IntegerLiteral) int {
	return n.Value
}

func (v *Calculator) VisitBinaryExpression(n *BinaryExpression) int {
	l := n.Left.Accept(v)
	r := n.Right.Accept(v)
	switch n.Op {
	case "+":
		return l + r
	case "-":
		return l - r
	case "*":
		return l * r
	case "/":
		return l / r
	}
	return 0
}

var instructions = []ByteCode{}

type Generator struct{}

func (v *Generator) VisitInteger(n *IntegerLiteral) int {
	instructions = append(instructions, ByteCode{
		Op:    push,
		Value: n.Value,
	})
	return 0
}

func (v *Generator) VisitBinaryExpression(n *BinaryExpression) int {
	n.Left.Accept(v)
	n.Right.Accept(v)
	instructions = append(instructions, ByteCode{
		Op: opMap[n.Op],
	})
	return 0
}

type ByteCode struct {
	Op    int
	Value int
}

var opMap = map[string]int{
	"+": plus,
	"-": minus,
	"*": mul,
	"/": div,
}
var stack = make([]int, 10)
var sp = 0

const (
	push = iota
	plus
	minus
	mul
	div
)

func execute() int {
	for _, code := range instructions {
		switch code.Op {
		case plus:
			stack[sp-2] += stack[sp-1]
			sp--
		case minus:
			stack[sp-2] -= stack[sp-1]
			sp--
		case mul:
			stack[sp-2] *= stack[sp-1]
			sp--
		case div:
			stack[sp-2] /= stack[sp-1]
			sp--
		case push:
			stack[sp] = code.Value
			sp++
		}
	}
	return stack[0]
}
