package parse

import (
	"errors"
	"language/module/ast"
	"language/module/token"
	"strconv"
)

var (
	ErrUnidentifiedKeyword = errors.New("unidentified keyword")
)

func (self *Parser) registerSymbol(name string) int {
	self.Symbols = append(self.Symbols, name)
	return len(self.Symbols) - 1
}

func (self *Parser) ParseStatement() (*ast.Node, error) {

	t, ok := self.peek()
	if !ok {
		return nil, ErrExpectedToken
	}

	if t.Type == token.Keyword {
		self.index += 1
	} else {
		return self.ParseExpression(0)
	}

	switch t.Value {
	case token.I32:

		identifier, ok := self.matchType(token.Identifier)
		if !ok {
			return nil, ErrExpectedToken
		}

		_, ok = self.matchType(token.EOL)
		if !ok {
			return nil, ErrExpectedToken
		}

		node := ast.New(ast.Declare)
		index := self.registerSymbol(string(identifier.Value))

		node.Value = index

		return node, nil

	case token.If:

		condition, err := self.ParseExpression(0)
		if err != nil {
			return nil, err
		}

		first, err := self.ParseStatements()
		if err != nil {
			return nil, err
		}

		node := ast.New(ast.If)
		node.Left = condition
		node.Right = ast.New(ast.Glue)

		node.Right.Insert(first)

		t, ok := self.peek()
		if ok && t.Value == token.Else {

			self.index++

			t, ok := self.consume()
			if !ok {
				return nil, ErrExpectedToken
			}

			if t.Value != token.CurlyL {
				return nil, ErrExpectedToken
			}

			second, err := self.ParseStatements()
			if err != nil {
				return nil, err
			}

			node.Right.Insert(second)
		}

		return node, nil

	case token.While:

		condition, err := self.ParseExpression(0)
		if err != nil {
			return nil, err
		}

		first, err := self.ParseStatements()
		if err != nil {
			return nil, err
		}

		node := ast.New(ast.While)

		node.Left = condition
		node.Right = first

		return node, nil

	case token.Interrupt:

		value, ok := self.matchType(token.Integer)
		if !ok {
			return nil, ErrExpectedToken
		}

		_, ok = self.matchType(token.EOL)
		if !ok {
			return nil, ErrExpectedToken
		}

		node := ast.New(ast.Interrupt)

		interrupt, err := strconv.ParseUint(string(value.Value), 10, 64)
		if err != nil {
			return nil, err
		}

		node.Value = int(interrupt)

		return node, nil
	}

	return nil, ErrUnidentifiedKeyword
}

func (self *Parser) ParseStatements() (*ast.Node, error) {

	nodes := ast.New(ast.Glue)

	root := nodes

	for {
		t, ok := self.peek()
		if !ok {
			break
		}

		if t.Type == token.EOF ||
			t.Value == token.CurlyR {
			self.index++
			break
		}

		if t.Type == token.EOL {
			self.index++
			continue
		}

		node, err := self.ParseStatement()

		if err != nil {
			return nil, err
		}

		ok = nodes.Insert(node)
		if !ok {
			swap := nodes.Right
			nodes.Right = ast.New(ast.Glue)
			nodes = nodes.Right

			nodes.Insert(swap)
			nodes.Insert(node)
		}
	}

	return root, nil
}
