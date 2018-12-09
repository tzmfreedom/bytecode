package main

import (
	"fmt"
	"os"
	"strconv"
	"github.com/k0kubun/pp"
)

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
	generator := &Generator{}
	root.Accept(generator)
	pp.Println(instructions)
}

var instructions = []string{}

type Generator struct{}

func (v *Generator) VisitInteger(n *IntegerLiteral) int {
	instructions = append(instructions, "push %d", n.Value)
	return 0
}

func (v *Generator) VisitBinaryExpression(n *BinaryExpression) int {
	n.Left.Accept(v)
	n.Right.Accept(v)
	instructions = append(instructions,"pop rdi")
	instructions = append(instructions,"pop rax")
	if n.Op == "*" {
		instructions = append(instructions,"mul rdi")
	} else if n.Op == "/" {
		instructions = append(instructions,"mov rdx, 0")
		instructions = append(instructions, "div rdi")
	} else {
		instructions = append(instructions, fmt.Sprintf("%s rax, rdi", opMap[n.Op]))
	}
	instructions = append(instructions,"push rax")
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
