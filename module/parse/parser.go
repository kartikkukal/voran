package parse

import (
	"errors"
	"language/module/ast"
	"language/module/token"
)

var (
	ErrExpectedToken = errors.New("expected token")
)

type Symbol struct {
	Name string
	Type ast.Type
}

type Parser struct {
	tokens []token.Token
	index  int

	Constants []string
	Symbols   []Symbol
}

func New(tokens []token.Token) Parser {
	return Parser{
		tokens: tokens,
		index:  0,
	}
}

func (self *Parser) peek() (token.Token, bool) {
	if self.index >= len(self.tokens) {
		return token.Token{}, false
	}

	return self.tokens[self.index], true
}

func (self *Parser) consume() (token.Token, bool) {
	defer func() {
		self.index += 1
	}()

	return self.peek()
}

func (self *Parser) matchToken(data string) (token.Token, bool) {
	t, ok := self.consume()
	if !ok {
		return token.Token{}, false
	}

	if t.Value != data {
		return t, false
	}

	return t, true
}

func (self *Parser) matchType(kind token.Type) (token.Token, bool) {
	t, ok := self.consume()
	if !ok {
		return token.Token{}, false
	}

	if t.Type != kind {
		return t, false
	}

	return t, true
}

func (self *Parser) registerString(data string) int {

	index := -1

	for i, v := range self.Constants {
		if data == v {
			index = i
			break
		}
	}

	if index == -1 {
		self.Constants = append(self.Constants, data)
		return len(self.Constants) - 1
	}

	return index
}

// Adds a symbol with its name and type and returns the index
func (self *Parser) registerSymbol(name string, t ast.Type) int {

	for _, v := range self.Symbols {
		if name == v.Name {
			// FIXME: Don't panic here
			panic("symbol already declared")
		}
	}

	self.Symbols = append(self.Symbols, Symbol{
		Name: name,
		Type: t,
	})
	return len(self.Symbols) - 1
}

// Find a symbol, if it exists and return its index and type
func (self *Parser) findSymbol(name string) (int, ast.Type) {

	for i, v := range self.Symbols {
		if name == v.Name {
			return i, v.Type
		}
	}

	// FIXME: Don't panic here
	panic("symbol not found")
}
