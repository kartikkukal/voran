package parse

import (
	"errors"
	"language/module/ast"
	"language/module/token"
)

var (
	ErrInvalidParameters   = errors.New("invalid keyword in parameter list")
	ErrUnidentifiedKeyword = errors.New("unidentified keyword")

	ErrInvalidToken = errors.New("invalid token")

	ErrExpectedType = errors.New("expected type")
)

func (self *Parser) ParseParameters() (*ast.Node, error) {

	parseSingle := func() (*ast.Node, error) {
		name, ok := self.matchType(token.Identifier)
		if !ok {
			return nil, ErrExpectedToken
		}

		next, ok := self.consume()
		if !ok {
			return nil, ErrExpectedToken
		}

		astType, err := self.ParseType(next.Value)
		if err != nil {
			return nil, err
		}

		index := self.registerSymbol(name.Value, astType)

		return ast.NewWithValue(ast.Parameter, index, astType), nil
	}

	_, ok := self.matchToken(token.ParenL)
	if !ok {
		return nil, ErrExpectedToken
	}

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

		node, err := parseSingle()
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

// Parses a declaration of a variable or function of a given type.
//
// Returns either a single node with variable name and type
// or a function node with left tree containing parameters
// and right tree containing the body of the function.
func (self *Parser) ParseDeclaration() (*ast.Node, error) {

	identifier, ok := self.matchType(token.Identifier)
	if !ok {
		return nil, ErrExpectedToken
	}

	next, ok := self.consume()

	// For declaration of the kind:
	// var name = <expression>
	if next.Value == token.Equal {
		initializer, err := self.ParseExpression(0)
		if err != nil {
			return nil, err
		}

		astType, err := self.CoerceTypeCast(initializer)
		if err != nil {
			return nil, err
		}

		index := self.registerSymbol(identifier.Value, astType)
		node := ast.NewWithValue(ast.Declare, index, astType)

		node.Left = initializer
		return node, nil
	}

	// For declaration of the kind:
	// var name string = <expression>
	astType, err := self.ParseType(next.Value)
	if err != nil {
		return nil, err
	}

	index := self.registerSymbol(identifier.Value, astType)
	node := ast.NewWithValue(ast.Declare, index, astType)

	// Check if next token is EOL or '='
	next, ok = self.consume()

	if next.Value != token.Equal {
		if next.Type == token.EOL {
			return node, nil
		}

		return nil, ErrInvalidToken
	}

	initializer, err := self.ParseExpression(0)
	if err != nil {
		return nil, err
	}

	_, err = self.CoerceTypeCast(initializer)
	if err != nil {
		return nil, err
	}

	node.Left = initializer

	return node, nil
}

func (self *Parser) ParseFunction() (*ast.Node, error) {
	identifier, ok := self.matchType(token.Identifier)
	if !ok {
		return nil, ErrExpectedToken
	}

	// Type and index of identifier will be set later
	node := ast.New(ast.Function)
	node.Type = ast.None

	parameters, err := self.ParseParameters()
	if err != nil {
		return nil, err
	}

	node.Left = parameters

	next, ok := self.peek()
	if !ok {
		return nil, ErrExpectedToken
	}

	if next.Value != token.CurlyL {
		self.index++

		astType, err := self.ParseType(next.Value)
		if err != nil {
			return nil, err
		}

		node.Type = astType
	}

	// Set type and identifier index
	index := self.registerSymbol(identifier.Value, node.Type)
	node.Value = index

	body, err := self.ParseStatements()
	if err != nil {
		return nil, err
	}

	node.Right = body

	return node, nil
}

func (self *Parser) ParseBranch() (*ast.Node, error) {

	condition, err := self.ParseExpression(0)
	if err != nil {
		return nil, err
	}

	first, err := self.ParseStatements()
	if err != nil {
		return nil, err
	}

	node := ast.New(ast.Branch)
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

/*
func (self *Parser) ParseSyscall() (*ast.Node, error) {

}*/

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
	case token.Var:
		return self.ParseDeclaration()

	case token.Func:
		return self.ParseFunction()

	case token.If:
		return self.ParseBranch()

	case token.While:
		return self.ParseWhile()

		/*case token.Syscall:
		return self.ParseSyscall()*/
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
