package parse

import (
	n "language/module/node"
	t "language/module/token"
)

var (
	operators = [...]int{
		t.Add,
		t.Subtract,
		t.Multiply,
		t.Divide,
		t.Modulus,
		t.LessThan,
		t.GreaterThan,
		t.LessOrEqual,
		t.GreaterOrEqual,
	}
)

func bindingPower(token t.Token) (int, int) {
	switch token.Kind {
	case t.LessThan:
		return 2, 3
	case t.LessOrEqual:
		return 2, 3
	case t.GreaterThan:
		return 2, 3
	case t.GreaterOrEqual:
		return 2, 3
	case t.Add:
		return 4, 5
	case t.Subtract:
		return 4, 5
	case t.Multiply:
		return 6, 7
	case t.Divide:
		return 6, 7
	case t.Modulus:
		return 6, 7
	default:
		return 0, 1
	}
}

func (self *Parser) matchOperator() (t.Token, bool) {

	token, ok := self.consume()
	if !ok {
		return t.Token{}, false
	}

	found := false

	for _, value := range operators {
		if value == token.Kind {
			found = true
			break
		}
	}

	return token, found
}

func (self *Parser) ParseExpression(rbp int) *n.Node {

	var lhs *n.Node

	token, ok := self.consume()
	if !ok {
		return lhs
	}

	if token.Kind == t.ParenL {
		lhs = self.ParseExpression(0)
	}

	if token.Kind == t.IntConstant ||
		token.Kind == t.Identifier {
		lhs = n.New(token)
	}

	for {
		token, ok = self.matchOperator()
		if !ok {
			return lhs
		}

		if token.Kind == t.ParenR {
			return lhs
		}

		l, r := bindingPower(token)

		if rbp > l {
			return lhs
		}

		rhs := self.ParseExpression(r)
		self.index--

		op := n.NewWithChildren(token, lhs, rhs)

		lhs = op
	}
}
