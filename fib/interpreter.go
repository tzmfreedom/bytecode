package main

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
