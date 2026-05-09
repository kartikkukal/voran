package parse

import (
	"errors"
	"fmt"
	"language/module/ast"
	"language/module/token"
)

var (
	keywords = [...]int{
		token.I8,
		token.I16,
		token.I32,
		token.I64,
		token.U8,
		token.U16,
		token.U32,
		token.U64,
		token.If,
		token.While,
	}

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

	found := false

	for _, value := range keywords {
		if value == t.Type {
			found = true
			break
		}
	}
	if found {
		self.index += 1
	} else {
		return self.ParseExpression(0)
	}

	switch t.Type {
	case token.I32:

		identifier, ok := self.matchToken(token.Identifier)
		if !ok {
			return nil, ErrExpectedToken
		}

		_, ok = self.matchToken(token.EOL)
		if !ok {
			return nil, ErrExpectedToken
		}

		node := ast.New(ast.DeclareInt)
		index := self.registerSymbol(identifier.Data)

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

		t, _ := self.peek()

		fmt.Println("Found if", t.Debug())

		node := ast.New(ast.IfElse)

		node.Left = condition
		node.Right = first

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

	case token.Inter:

		return ast.New(ast.Inter), nil
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
			t.Type == token.CurlyR {
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
