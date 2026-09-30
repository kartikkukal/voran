package parse

import (
	"errors"
	"language/module/token"
)

var (
	ErrExpectedToken = errors.New("expected token")
)

type Parser struct {
	tokens []token.Token
	index  int

	Constants []string
	Symbols   []string
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

func (self *Parser) registerSymbol(name string) int {

	index := -1

	for i, v := range self.Symbols {
		if name == v {
			index = i
			break
		}
	}

	if index == -1 {
		self.Symbols = append(self.Symbols, name)
		return len(self.Symbols) - 1
	}

	return index
}
