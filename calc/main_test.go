package main

import (
	"fmt"
	"testing"

	"github.com/k0kubun/pp"
)

var root = &BinaryExpression{
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

var calculator = &Calculator{}

func init() {
	calculator := &Calculator{}
	result := runIntepreter(root, calculator)
	fmt.Println(result)
	generator := &Generator{}
	root.Accept(generator)
	pp.Println(instructions)
	fmt.Println(execute())
}

func BenchmarkInterpreter(b *testing.B) {
	for n := 0; n < b.N; n++ {
		root.Accept(calculator)
	}
}

func BenchmarkVM(b *testing.B) {
	for n := 0; n < b.N; n++ {
		sp = 0
		execute()
	}
}
