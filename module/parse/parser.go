package parse

import (
	t "language/module/token"
)

type Parser struct {
	tokens []t.Token
	index  int
}

func New(tokens []t.Token) Parser {
	return Parser{
		tokens: tokens,
		index:  0,
	}
}

func (self *Parser) peek() (t.Token, bool) {
	if self.index >= len(self.tokens) {
		return t.Token{}, false
	}

	return self.tokens[self.index], true
}

func (self *Parser) consume() (t.Token, bool) {
	defer func() {
		self.index++
	}()

	return self.peek()
}

func (self Parser) matchToken(kind int) (t.Token, bool) {
	token, ok := self.consume()
	if !ok {
		return t.Token{}, false
	}

	if token.Kind != kind {
		return t.Token{}, false
	}

	return token, true
}
