package parse

import (
	"errors"
	"language/module/ast"
	"language/module/token"
	"strconv"
)

var (
	ErrInvalidParameters   = errors.New("invalid keyword in parameter list")
	ErrUnidentifiedKeyword = errors.New("unidentified keyword")
)

func (self *Parser) ParseParameters() (*ast.Node, error) {

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

		node, err := self.ParseStatement()
		if err != nil {
			return nil, err
		}

		if node.Type != ast.Declare {
			return nil, ErrInvalidParameters
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

func (self *Parser) ParseDeclaration() (*ast.Node, error) {

	identifier, ok := self.matchType(token.Identifier)
	if !ok {
		return nil, ErrExpectedToken
	}

	index := self.registerSymbol(string(identifier.Value))

	node := ast.New(ast.Declare)
	node.Value = index

	next, ok := self.consume()

	if next.Type == token.EOL ||
		next.Value == token.Comma {
		return node, nil
	}

	if next.Value == token.ParenR {
		self.index--
		return node, nil
	}

	if next.Value == token.Equals {
		self.index = self.index - 2

		initializer, err := self.ParseExpression(0)
		if err != nil {
			return nil, err
		}

		node.Left = initializer

		return node, nil
	}

	if next.Value != token.ParenL {
		return nil, ErrExpectedToken
	}

	parameters, err := self.ParseParameters()
	if err != nil {
		return nil, err
	}

	body, err := self.ParseStatements()
	if err != nil {
		return nil, err
	}

	node = ast.New(ast.Function)

	node.Value = index
	node.Left = parameters
	node.Right = body

	return node, nil
}

/*
func (self *Parser) ParseFunction() (*ast.Node, error) {

}
*/
func (self *Parser) ParseBranch() (*ast.Node, error) {

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
}

func (self *Parser) ParseWhile() (*ast.Node, error) {
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
}

func (self *Parser) ParseInterrupt() (*ast.Node, error) {
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
		return self.ParseDeclaration()

	case token.If:
		return self.ParseBranch()

	case token.While:
		return self.ParseWhile()

	case token.Interrupt:
		return self.ParseInterrupt()
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
