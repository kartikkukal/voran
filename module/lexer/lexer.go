package lexer

import (
	"errors"
	"fmt"
	"language/module/token"
	"slices"
)

var (
	ErrUnknownCharacter = errors.New("unknown character encountered")
	ErrUnknownToken     = errors.New("unknown token encountered")

	alphabet = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	decimal  = []byte("1234567890")

	identifierStart = append(alphabet, '_')
	identifier      = append(identifierStart, decimal...)

	punctuators     = []byte(".,:;[]{}()")
	operators       = []byte("=+-*/%&|^!<>")
	whitespace      = []byte(" \t\r")
	stringDelimiter = byte('"')
	charDelimiter   = byte('\'')
	newline         = byte('\n')
	comment         = byte('/')
)

type Lexer struct {
	source []byte
	index  int

	line   int
	column int
}

func New(source []byte) Lexer {
	return Lexer{
		source: source,
		index:  0,

		line:   1,
		column: 0,
	}
}

func (self *Lexer) peek() (byte, bool) {

	if self.index >= len(self.source) {
		return 0, false
	}

	return self.source[self.index], true
}

func (self *Lexer) consume() (byte, bool) {

	defer func() {
		self.index++
		self.column++
	}()

	return self.peek()
}

func (self *Lexer) advance() {
	self.line += 1
	self.column = 0
}

func (self *Lexer) Tokenize() ([]token.Token, error) {
	tokens := make([]token.Token, 0)
	buffer := make([]byte, 0)

	wrapper := func(err error) error {
		return fmt.Errorf("error encountered at line %v:%v: %w", self.line, self.column, err)
	}

	flush := func(data string, kind token.Type, line, column int) {

		t := token.New(data, kind, line, column)

		tokens = append(tokens, t)
		buffer = nil
	}

	read := func(characters []byte) {
		for {
			c, ok := self.peek()
			if !ok {
				break
			}

			if !slices.Contains(characters, c) {
				break
			}

			c, ok = self.consume()
			if !ok {
				break
			}

			buffer = append(buffer, c)
		}
	}

	delimiter := func(terminal byte) {
		for {
			c, ok := self.consume()
			if !ok {
				break
			}

			if c == terminal {
				break
			}

			buffer = append(buffer, c)
		}
	}

	for {
		c, ok := self.consume()
		if !ok {
			break
		}

		line := self.line
		column := self.column

		if slices.Contains(whitespace, c) {
			continue
		}

		if c == newline {

			flush("", token.EOL, line, column)
			self.advance()

			continue
		}

		if slices.Contains(identifierStart, c) {

			buffer = append(buffer, c)
			read(identifier)

			value := string(buffer)

			if token.TypeOf(value) == token.Keyword {
				flush(value, token.Keyword, line, column)
				continue
			}

			flush(value, token.Identifier, line, column)
			continue
		}

		if slices.Contains(operators, c) {

			buffer = append(buffer, c)
			read(operators)

			value := string(buffer)

			if token.TypeOf(value) != token.Operator {
				return nil, wrapper(ErrUnknownToken)
			}

			flush(value, token.Operator, line, column)
			continue
		}

		if slices.Contains(punctuators, c) {

			flush(string(c), token.Punctuator, line, column)
			continue
		}

		if slices.Contains(decimal, c) {

			buffer = append(buffer, c)
			read(decimal)

			value := string(buffer)

			flush(value, token.LiteralInt, line, column)
			continue
		}

		if c == stringDelimiter {

			delimiter(stringDelimiter)

			value := string(buffer)

			flush(value, token.LiteralStr, line, column)
			continue
		}

		if c == charDelimiter {

			delimiter(charDelimiter)

			value := string(buffer)

			flush(value, token.LiteralStr, line, column)
			continue
		}

		return nil, wrapper(ErrUnknownCharacter)
	}

	flush("", token.EOL, self.line, self.column)
	flush("", token.EOF, self.line, self.column)
	return tokens, nil
}
