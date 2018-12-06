package main

import (
	"fmt"
	"os"
	// "github.com/k0kubun/pp"
	"strconv"
)

func main() {
	val, _ := strconv.Atoi(os.Args[1])
	root := &Call{
		Name: "fib",
		Arg: &IntegerLiteral{
			Value: val,
		},
	}
	calculator := &Calculator{}
	result := root.Accept(calculator)
	fmt.Println(result)
	// generator := &Generator{}
	// root.Accept(generator)
	// pp.Println(instructions)
	// fmt.Println(execute())
}

var env = map[string]int{}
var functions = map[string]Node{
	"fib": &Body{
		Statement: []Node{
			&If{
				Condition: &BinaryExpression{
					Left: &Identifier{
						Name: "i",
					},
					Right: &IntegerLiteral{
						Value: 2,
					},
					Op: "<",
				},
				TrueStatement: []Node{
					&Return{
						Expression: &Identifier{
							Name: "i",
						},
					},
				},
				FalseStatement: []Node{
					&Return{
						Expression: &BinaryExpression{
							Left: &Call{
								Name: "fib",
								Arg: &BinaryExpression{
									Left: &Identifier{
										Name: "i",
									},
									Right: &IntegerLiteral{
										Value: 1,
									},
									Op: "-",
								},
							},
							Right: &Call{
								Name: "fib",
								Arg: &BinaryExpression{
									Left: &Identifier{
										Name: "i",
									},
									Right: &IntegerLiteral{
										Value: 2,
									},
									Op: "-",
								},
							},
							Op: "+",
						},
					},
				},
			},
		},
	},
}

type Body struct {
	Statement []Node
}

func (n *Body) Accept(v Visitor) int {
	return v.VisitBody(n)
}

type If struct {
	Condition      *BinaryExpression
	TrueStatement  []Node
	FalseStatement []Node
}

func (n *If) Accept(v Visitor) int {
	return v.VisitIf(n)
}

type Return struct {
	Expression Node
}

func (n *Return) Accept(v Visitor) int {
	return v.VisitReturn(n)
}

type Identifier struct {
	Name string
}

func (n *Identifier) Accept(v Visitor) int {
	return v.VisitIdentifier(n)
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

type Call struct {
	Name string
	Arg  Node
}

func (n *Call) Accept(v Visitor) int {
	return v.VisitCall(n)
}

type Visitor interface {
	VisitBody(*Body) int
	VisitIf(*If) int
	VisitIdentifier(*Identifier) int
	VisitReturn(*Return) int
	VisitCall(*Call) int
	VisitInteger(*IntegerLiteral) int
	VisitBinaryExpression(*BinaryExpression) int
}

type Node interface {
	Accept(v Visitor) int
}

type Calculator struct{}

func (v *Calculator) VisitBody(n *Body) int {
	var r int
	for _, stmt := range n.Statement {
		r = stmt.Accept(v)
	}
	return r
}

func (v *Calculator) VisitIf(n *If) int {
	var r int
	if n.Condition.Accept(v) == 1 {
		for _, stmt := range n.TrueStatement {
			r = stmt.Accept(v)
		}
		return r
	}
	for _, stmt := range n.FalseStatement {
		r = stmt.Accept(v)
	}
	return r
}

func (v *Calculator) VisitIdentifier(n *Identifier) int {
	return env[n.Name]
}

func (v *Calculator) VisitReturn(n *Return) int {
	return n.Expression.Accept(v)
}

func (v *Calculator) VisitCall(n *Call) int {
	newEnv := map[string]int{}
	newEnv["i"] = n.Arg.Accept(v)
	pre := env
	env = newEnv
	value := functions[n.Name].Accept(v)
	env = pre
	return value
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
	case "<":
		if l < r {
			return 1
		}
	}
	return 0
}

func (v *Calculator) VisitInteger(n *IntegerLiteral) int {
	return n.Value
}

var instructions = []ByteCode{}
var varIndex = map[string]int{}

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

func (v *Generator) VisitBody(n *Body) int {
	for _, stmt := range n.Statement {
		stmt.Accept(v)
	}
	return 0
}

func (v *Generator) VisitIf(n *If) int {
	n.Condition.Accept(v)
	ifPos := len(instructions)
	for _, stmt := range n.TrueStatement {
		stmt.Accept(v)
	}
	elsePos := len(instructions)
	for _, stmt := range n.FalseStatement {
		stmt.Accept(v)
	}
	instructions = append(instructions, ByteCode{
		Op:     ifOp,
		Value:  ifPos,
		Value2: elsePos,
	})
	return 0
}

func (v *Generator) VisitIdentifier(n *Identifier) int {
	instructions = append(instructions, ByteCode{
		Op: ident,
		Value: varIndex[n.Name],
	})
	return 0
}

func (v *Generator) VisitReturn(n *Return) int {
	n.Expression.Accept(v)
	instructions = append(instructions, ByteCode{
		Op:    ret,
	})
	return 0
}

func (v *Generator) VisitCall(n *Call) int {
	n.Arg.Accept(v)
	instructions = append(instructions, ByteCode{
		Op: call,
		Value: 0, // TODO: implement
	})
	return 0
}

type ByteCode struct {
	Op     int
	Value  int
	Value2 int
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
	ret
	call
	ident
	ifOp
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
