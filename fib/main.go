package main

import (
	"fmt"
	"os"
	"time"
	// "github.com/k0kubun/pp"

	"strconv"
)

func main() {
	val, _ := strconv.Atoi(os.Args[1])
	var root Node
	root = &Call{
		Name: "fib",
		Arg: &IntegerLiteral{
			Value: val,
		},
	}
	//calculator := &Calculator{}
	//d := lapseTime(func() {
	//	result := root.Accept(calculator)
	//	fmt.Println(result)
	//})
	//fmt.Printf("time: %f\n", d.Seconds())

	g := &Generator{}
	functions["fib"].Accept(g)
	functionCode["fib"] = g.Instructions

	g = &Generator{}
	root.Accept(g)
	//for i, ins := range g.Instructions {
	//	fmt.Println(fmt.Sprintf("%02d: %s", i, ins.ToString()))
	//}
	d := lapseTime(func() {
		result := execute(g.Instructions)
		fmt.Println(result)
	})
	fmt.Printf("time: %f\n", d.Seconds())
}

func lapseTime(f func()) time.Duration {
	start := time.Now()
	f()
	return time.Now().Sub(start)
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

var functionCode = map[string][]ByteCode{}

var varIndex = map[string]int{}

type Generator struct {
	Instructions []ByteCode
}

func (v *Generator) AddInstructions(op int, value int) {
	v.Instructions = append(v.Instructions, ByteCode{
		Op:    op,
		Value: value,
	})
}

func (v *Generator) VisitInteger(n *IntegerLiteral) int {
	v.AddInstructions(push, n.Value)
	return 0
}

func (v *Generator) VisitBinaryExpression(n *BinaryExpression) int {
	n.Left.Accept(v)
	n.Right.Accept(v)
	v.AddInstructions(opMap[n.Op], 0)
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
	jumpIfIndex := len(v.Instructions)
	v.AddInstructions(jumpIf, 0)
	for _, stmt := range n.FalseStatement {
		stmt.Accept(v)
	}
	jumpIndex := len(v.Instructions)
	v.AddInstructions(jump, 0)

	v.Instructions[jumpIfIndex].Value = len(v.Instructions)

	for _, stmt := range n.TrueStatement {
		stmt.Accept(v)
	}
	v.Instructions[jumpIndex].Value = len(v.Instructions)
	return 0
}

func (v *Generator) VisitIdentifier(n *Identifier) int {
	v.AddInstructions(ident, varIndex[n.Name])
	return 0
}

func (v *Generator) VisitReturn(n *Return) int {
	n.Expression.Accept(v)
	v.AddInstructions(ret, 0)
	return 0
}

func (v *Generator) VisitCall(n *Call) int {
	n.Arg.Accept(v)
	v.AddInstructions(call, 0) // TODO: implement
	return 0
}

type ByteCode struct {
	Op     int
	Value  int
	Value2 int
}

func (b *ByteCode) ToString() string {
	op := ""
	switch b.Op {
	case push:
		op = "push"
	case plus:
		op = "plus"
	case minus:
		op = "minus"
	case mul:
		op = "mul"
	case div:
		op = "div"
	case lt:
		op = "lt"
	case ret:
		op = "ret"
	case call:
		op = "call"
	case ident:
		op = "ident"
	case jumpIf:
		op = "jumpIf"
	case jump:
		op = "jump"
	}
	return fmt.Sprintf("%s %d %d", op, b.Value, b.Value2)
}

var opMap = map[string]int{
	"+": plus,
	"-": minus,
	"*": mul,
	"/": div,
	"<": lt,
}
var stack = make([]int, 100)
var sp = 0

const (
	push = iota
	plus
	minus
	mul
	div
	lt
	ret
	call
	ident
	jumpIf
	jump
)

func execute(instructions []ByteCode) int {
	arg := make([]int, 1)
	if sp > 0 {
		arg[0] = stack[sp-1]
	}
	pc := 0
	max := len(instructions)
	for true {
		code := instructions[pc]
		// fmt.Printf("%02d %s\n", pc, code.ToString())
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
		case lt:
			if stack[sp-2] < stack[sp-1] {
				stack[sp-2] = 1
			} else {
				stack[sp-2] = 0
			}
			sp--
		case push:
			stack[sp] = code.Value
			sp++
		case ret:
			sp--
			// fmt.Printf(">> %d => %d\n", arg[0], stack[sp])
			stack[sp-1] = stack[sp]
			return stack[sp-1]
		case call:
			execute(functionCode["fib"])
			// debug()
		case ident:
			stack[sp] = arg[0]
			sp++
		case jumpIf:
			condition := stack[sp-1]
			sp--
			if condition == 1 {
				pc = code.Value
				continue
			}
		case jump:
			pc = code.Value
			continue
		}
		// debug()
		pc++
		if pc >= max {
			break
		}
	}
	return stack[0]
}

func debug() {
	for i, s := range stack {
		if i >= sp {
			fmt.Println()
			return
		}
		fmt.Printf("%d : ", s)
	}
}
