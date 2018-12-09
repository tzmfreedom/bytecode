package main

import (
	"fmt"
	"os"
	"strconv"
)

//var root = &BinaryExpression{
//	Left: &BinaryExpression{
//		Left: &IntegerLiteral{
//			Value: 25,
//		},
//		Right: &IntegerLiteral{
//			Value: 5,
//		},
//		Op: "/",
//	},
//	Right: &BinaryExpression{
//		Left: &IntegerLiteral{
//			Value: 4,
//		},
//		Right: &IntegerLiteral{
//			Value: 3,
//		},
//		Op: "*",
//	},
//	Op: "+",
//}

var root = &Call{
	Name: "plus",
	Args: []Node{
		&IntegerLiteral{
			Value: 120,
		},
		&IntegerLiteral{
			Value: 23,
		},
	},
}

var plus = &Body{
	Statement: []Node{
		&Return{
			Expression: &BinaryExpression{
				Left:  &Identifier{Name: "i"},
				Right: &Identifier{Name: "j"},
				Op:    "+",
			},
		},
	},
}

func main() {
	g := &Generator{}
	root.Accept(g)
	fmt.Printf(`.intel_syntax noprefix
.global _main, _plus

_main:
`)
	for _, ins := range g.Instructions {
		fmt.Println("  " + ins)
	}
	fmt.Println(`
  ret`)

	createFunction("plus", plus)
}

func createFunction(name string, root Node) {
	g := &Generator{}
	root.Accept(g)
	fmt.Printf(`
_%s:
  push rbp
  mov rbp, rsp
`, name)
	for _, ins := range g.Instructions {
		fmt.Println("  " + ins)
	}
}

type Generator struct {
	Instructions []string
}

func (v *Generator) AddInstruction(src string) {
	v.Instructions = append(v.Instructions, src)
}

func (v *Generator) VisitInteger(n *IntegerLiteral) int {
	v.Instructions = append(v.Instructions, fmt.Sprintf("push %d", n.Value))
	return 0
}

func (v *Generator) VisitBinaryExpression(n *BinaryExpression) int {
	n.Left.Accept(v)
	n.Right.Accept(v)
	v.AddInstruction("pop rdi")
	v.AddInstruction("pop rax")
	if n.Op == "*" {
		v.AddInstruction("mul rdi")
	} else if n.Op == "/" {
		v.AddInstruction("mov rdx, 0")
		v.AddInstruction("div rdi")
	} else {
		v.AddInstruction(fmt.Sprintf("%s rax, rdi", opMap[n.Op]))
	}
	v.AddInstruction("push rax")
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
	// jumpIfIndex
	// v.AddInstruction("jump", varIndex[n.Name])
	for _, stmt := range n.FalseStatement {
		stmt.Accept(v)
	}
	// jumpIndex
	// v.AddInstruction("jump", varIndex[n.Name])

	for _, stmt := range n.TrueStatement {
		stmt.Accept(v)
	}
	return 0
}

func (v *Generator) VisitIdentifier(n *Identifier) int {
	address := 16
	if n.Name == "j" {
		address = 24
	}
	v.AddInstruction(fmt.Sprintf("push [rbp+%d]", address))
	return 0
}

func (v *Generator) VisitReturn(n *Return) int {
	n.Expression.Accept(v)
	v.AddInstruction("pop rax")
	v.AddInstruction("mov rsp, rbp")
	v.AddInstruction("pop rbp")
	v.AddInstruction("ret")
	return 0
}

func (v *Generator) VisitCall(n *Call) int {
	for _, arg := range n.Args {
		arg.Accept(v)
	}
	v.AddInstruction(fmt.Sprintf("call _%s", n.Name))
	v.AddInstruction(fmt.Sprintf("add rsp, %d", len(n.Args)*8))
	return 0
}

var opMap = map[string]string{
	"+": "add",
	"-": "sub",
	"*": "mul",
	"/": "div",
}

func main02() {
	fmt.Printf(`.intel_syntax noprefix
.global _main

_main:
	push 5
	push 2
	pop rdi
	pop rax
	add rax, rdi
	ret`)
}

func main01() {
	num, _ := strconv.Atoi(os.Args[1])
	fmt.Printf(`.intel_syntax noprefix
.global _main

_main:
  mov rax, %d 
  ret
`, num)
}
