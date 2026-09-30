package parse

import (
	"errors"
	"language/module/ast"
	"language/module/token"
	"strconv"
)

var (
	precedenceTable = map[string]int{
		token.Equal:     -1,
		token.NotEq:     2,
		token.Equals:    2,
		token.Less:      2,
		token.LessEq:    2,
		token.Greater:   2,
		token.GreaterEq: 2,
		token.Add:       3,
		token.Sub:       3,
		token.Mul:       3,
		token.Div:       3,
		token.Mod:       3,
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

func astNodeType(t string) (ast.Kind, bool) {
	var node ast.Kind

	switch t {
	case token.Add:
		node = ast.Add
	case token.Sub:
		node = ast.Sub
	case token.Mul:
		node = ast.Mul
	case token.Div:
		node = ast.Div
	case token.Mod:
		node = ast.Mod
	case token.Equals:
		node = ast.Equals
	case token.NotEq:
		node = ast.NotEq
	case token.Less:
		node = ast.Less
	case token.LessEq:
		node = ast.LessEq
	case token.Greater:
		node = ast.Greater
	case token.GreaterEq:
		node = ast.GreaterEq
	case token.Equal:
		node = ast.Equal
	default:
		return ast.Invalid, false
	}

	return node, true
}

func (self *Parser) parseArguments() (*ast.Node, error) {

	master := ast.New(ast.Glue)
	root := master

	for {
		t, ok := self.peek()
		if !ok {
			return nil, ErrExpectedToken
		}

		if t.Value == token.ParenR {
			self.index++
			return root, nil
		}

		node, err := self.ParseExpression(0)
		if err != nil {
			return nil, err
		}

		ok = master.Insert(node)
		if !ok {
			swap := master.Right
			master.Right = ast.New(ast.Glue)
			master = master.Right

			master.Insert(swap)
			master.Insert(node)
		}
	}
}

func (self *Parser) ParseExpression(rbp int) (*ast.Node, error) {

	var lhs *ast.Node

	t, ok := self.consume()
	if !ok {
		return nil, ErrExpectedToken
	}

	if t.Type == token.EOL && t.Value == token.Comma {
		return lhs, nil
	}

	if t.Value == token.ParenL {
		var err error

		lhs, err = self.ParseExpression(0)
		if err != nil {
			return nil, err
		}
	}

	if t.Type == token.LiteralInt {
		value, err := strconv.ParseUint(string(t.Value), 10, 64)
		if err != nil {
			return nil, err
		}

		lhs = ast.NewWithValue(ast.Literal, int(value), ast.Int)
	}

	if t.Type == token.Identifier {
		index := self.registerSymbol(string(t.Value))

		next, ok := self.peek()
		if ok && next.Value == token.ParenL {
			lhs = ast.NewWithValue(ast.Call, index, ast.None)

			self.index++

			args, err := self.parseArguments()
			if err != nil {
				return nil, err
			}

			lhs.Left = args
		} else {
			lhs = ast.NewWithValue(ast.Identifier, index, ast.None)
		}
	}

	for {
		t, ok = self.matchType(token.Operator)
		if !ok {
			if t.Type == token.EOL || t.Value == token.CurlyL || t.Value == token.Comma {
				return lhs, nil
			}

			if t.Value == token.ParenR {
				self.index--
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
