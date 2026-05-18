package parse

import (
	"errors"
	"language/module/ast"
	"language/module/token"
	"slices"
	"strconv"
)

var (
	precedenceTable = map[token.Value]int{
		token.Equals:         -1,
		token.NotEqual:       2,
		token.Equal:          2,
		token.LessThan:       2,
		token.LessOrEqual:    2,
		token.GreaterThan:    2,
		token.GreaterOrEqual: 2,
		token.Add:            3,
		token.Subtract:       3,
		token.Multiply:       3,
		token.Divide:         3,
		token.Modulus:        3,
	}

	ErrUndeclaredIdentifier = errors.New("identifier has not been declared before")
	ErrExpectedOperator     = errors.New("expected operator")
	ErrUnidentifiedOperator = errors.New("unidentified operator")
)

func bindingPower(token token.Token) (int, int) {
	power, ok := precedenceTable[token.Value]
	if !ok {
		return 0, 0
	}

	if power < 0 {
		return -(power*2 - 1), -(power * 2)
	} else {
		return (power * 2), (power*2 + 1)
	}
}

func astNodeType(t token.Value) (ast.Type, bool) {
	var node ast.Type

	switch t {
	case token.Add:
		node = ast.Add
	case token.Subtract:
		node = ast.Subtract
	case token.Multiply:
		node = ast.Multiply
	case token.Divide:
		node = ast.Divide
	case token.Modulus:
		node = ast.Modulus
	case token.Equal:
		node = ast.Equal
	case token.NotEqual:
		node = ast.NotEqual
	case token.LessThan:
		node = ast.LessThan
	case token.LessOrEqual:
		node = ast.LessOrEqual
	case token.GreaterThan:
		node = ast.GreaterThan
	case token.GreaterOrEqual:
		node = ast.GreaterOrEqual
	case token.Equals:
		node = ast.Equals
	default:
		return ast.Type(""), false
	}

	return node, true
}

func (self *Parser) findSymbol(name string) (int, bool) {
	symbol := slices.Index(self.Symbols, name)
	if symbol == -1 {
		return 0, false
	}
	return symbol, true
}

func (self *Parser) ParseExpression(rbp int) (*ast.Node, error) {

	var lhs *ast.Node

	t, ok := self.consume()
	if !ok {
		return nil, ErrExpectedToken
	}

	if t.Type == token.EOL {
		return lhs, nil
	}

	if t.Value == token.ParenL {
		var err error

		lhs, err = self.ParseExpression(0)
		if err != nil {
			return nil, err
		}
	}

	if t.Type == token.Integer {
		value, err := strconv.ParseUint(string(t.Value), 10, 64)
		if err != nil {
			return nil, err
		}

		lhs = ast.NewWithValue(ast.Integer, int(value))
	}

	if t.Type == token.Identifier {
		index, ok := self.findSymbol(string(t.Value))
		if !ok {
			return nil, ErrUndeclaredIdentifier
		}

		lhs = ast.NewWithValue(ast.Identifier, index)
	}

	for {
		t, ok = self.matchType(token.Operator)
		if !ok {
			if t.Type == token.EOL || t.Value == token.CurlyL {
				return lhs, nil
			}

			return nil, ErrExpectedOperator
		}

		if t.Value == token.ParenR {
			return lhs, nil
		}

		l, r := bindingPower(t)

		if rbp > l {
			return lhs, nil
		}

		rhs, err := self.ParseExpression(r)
		if err != nil {
			return nil, err
		}

		self.index--

		node, ok := astNodeType(t.Value)
		if !ok {
			return nil, ErrUnidentifiedOperator
		}

		op := ast.NewWithChildren(node, lhs, rhs)

		if node != ast.Equals {
			op.RValue = true

		} else {
			op.Left.RValue = false

			swap := op.Left
			op.Left = op.Right
			op.Right = swap
		}

		lhs = op
	}
}
