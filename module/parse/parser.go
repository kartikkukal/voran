package parse

import (
	"errors"
	"fmt"
	"language/module/token"
)

var (
	ErrExpectedToken = errors.New("expected token")
)

type Parser struct {
	tokens []token.Token
	index  int

	Symbols []string
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

func (self *Parser) matchToken(kind int) (token.Token, bool) {
	t, ok := self.consume()
	if !ok {
		return token.Token{}, false
	}

	if t.Type != kind {
		fmt.Println("what we found:", t.Type, "expected:", kind)
		fmt.Println(t.Debug())
		return token.Token{}, false
	}

	return t, true
}
