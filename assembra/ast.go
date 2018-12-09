package main

//type Body struct {
//	Statement []Node
//}
//
//func (n *Body) Accept(v Visitor) int {
//	return v.VisitBody(n)
//}
//
//type If struct {
//	Condition      *BinaryExpression
//	TrueStatement  []Node
//	FalseStatement []Node
//}
//
//func (n *If) Accept(v Visitor) int {
//	return v.VisitIf(n)
//}
//
//type Return struct {
//	Expression Node
//}
//
//func (n *Return) Accept(v Visitor) int {
//	return v.VisitReturn(n)
//}
//
//type Identifier struct {
//	Name string
//}
//
//func (n *Identifier) Accept(v Visitor) int {
//	return v.VisitIdentifier(n)
//}

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

//type Call struct {
//	Name string
//	Arg  Node
//}
//
//func (n *Call) Accept(v Visitor) int {
//	return v.VisitCall(n)
//}

type Visitor interface {
	// VisitBody(*Body) int
	// VisitIf(*If) int
	// VisitIdentifier(*Identifier) int
	// VisitReturn(*Return) int
	// VisitCall(*Call) int
	VisitInteger(*IntegerLiteral) int
	VisitBinaryExpression(*BinaryExpression) int
}

type Node interface {
	Accept(v Visitor) int
}
