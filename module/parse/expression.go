package parse

import (
	"errors"
	"language/module/ast"
	"language/module/token"
	"slices"
	"strconv"
)

var (
	operators = [...]int{
		token.Add,
		token.Subtract,
		token.Multiply,
		token.Divide,
		token.Modulus,
		token.IsEqual,
		token.NotEqual,
		token.LessThan,
		token.GreaterThan,
		token.LessOrEqual,
		token.GreaterOrEqual,
		token.Equal,
		token.ParenL,
		token.ParenR,
	}

	precedenceTable = map[int]int{
		token.Equal:          -1,
		token.NotEqual:       2,
		token.IsEqual:        2,
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
	power, ok := precedenceTable[token.Type]
	if !ok {
		return 0, 0
	}

	if power < 0 {
		return -(power*2 - 1), -(power * 2)
	} else {
		return (power * 2), (power*2 + 1)
	}
}

func (self *Parser) matchOperator() (token.Token, bool) {

	t, ok := self.consume()
	if !ok {
		return token.Token{}, false
	}

	found := false

	for _, value := range operators {
		if value == t.Type {
			found = true
			break
		}
	}

	return t, found
}

func astNodeType(t int) (int, bool) {
	var node int

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
	case token.IsEqual:
		node = ast.IsEqual
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
	case token.Equal:
		node = ast.Equal
	default:
		return 0, false
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

	if t.Type == token.ParenL {
		var err error

		lhs, err = self.ParseExpression(0)
		if err != nil {
			return nil, err
		}
	}

	if t.Type == token.LiteralInt {
		value, err := strconv.ParseUint(t.Data, 10, 64)
		if err != nil {
			return nil, err
		}

		lhs = ast.NewWithValue(ast.LiteralInt, int(value))
	}

	if t.Type == token.Identifier {
		index, ok := self.findSymbol(t.Data)
		if !ok {
			return nil, ErrUndeclaredIdentifier
		}

		lhs = ast.NewWithValue(ast.Identifier, index)
	}

	for {
		t, ok = self.matchOperator()
		if !ok {
			if t.Type == token.EOL || t.Type == token.CurlyL {
				return lhs, nil
			}

			return nil, ErrExpectedOperator
		}

		if t.Type == token.ParenR {
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

		node, ok := astNodeType(t.Type)
		if !ok {
			return nil, ErrUnidentifiedOperator
		}

		op := ast.NewWithChildren(node, lhs, rhs)

		if node != ast.Equal {
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
