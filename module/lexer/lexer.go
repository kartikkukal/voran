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

	line int
	col  int
}

func New(source []byte) Lexer {
	return Lexer{
		source: source,
		index:  0,

		line: 1,
		col:  0,
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
		self.col++
	}()

	return self.peek()
}

func (self *Lexer) advance() {
	self.line += 1
	self.col = 0
}

func (self *Lexer) Tokenize() ([]token.Token, error) {
	tokens := make([]token.Token, 0)
	buffer := make([]byte, 0)

	wrapper := func(err error) error {
		return fmt.Errorf("error encountered at line %v:%v: %w", self.line, self.col, err)
	}

	flush := func(kind, line, col int) {

		value := string(buffer)
		t := token.New(value, kind, line, col)

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
		col := self.col

		if slices.Contains(whitespace, c) {
			continue
		}

		if c == newline {

			flush(token.EOL, self.line, self.col)
			self.advance()

			continue
		}

		if slices.Contains(identifierStart, c) {

			buffer = append(buffer, c)
			read(identifier)

			value := string(buffer)

			t, ok := token.Keyword(value, line, col)
			if ok {
				tokens = append(tokens, t)
				buffer = nil
				continue
			}

			t = token.New(value, token.Identifier, line, col)
			tokens = append(tokens, t)
			buffer = nil

			continue
		}

		if slices.Contains(operators, c) {

			t, ok := token.Identify(c, self.line, self.col)
			if !ok {
				return nil, wrapper(ErrUnknownCharacter)
			}

			n, ok := self.peek()
			if !ok {
				break
			}

			kNext, ok := token.Identify(n, self.line, self.col)
			if ok && kNext.Type == token.Equal {

				t, ok := token.Augmented(t.Type, self.line, self.col)
				if !ok {
					return nil, wrapper(ErrUnknownToken)
				}

				fmt.Println("Augmented", t)

				tokens = append(tokens, t)

				_, ok = self.consume()
				if !ok {
					break
				}

				continue
			}

			tokens = append(tokens, t)
			continue
		}

		if slices.Contains(punctuators, c) {

			t, ok := token.Identify(c, self.line, self.col)
			if !ok {
				return nil, wrapper(ErrUnknownCharacter)
			}

			tokens = append(tokens, t)
			continue
		}

		if slices.Contains(decimal, c) {

			buffer = append(buffer, c)
			read(decimal)

			flush(token.LiteralInt, line, col)

			continue
		}

		if c == stringDelimiter {

			delimiter(stringDelimiter)
			flush(token.LiteralStr, line, col)

			continue
		}

		if c == charDelimiter {

			delimiter(charDelimiter)
			flush(token.LiteralChar, line, col)

			continue
		}

		return nil, wrapper(ErrUnknownCharacter)
	}

	flush(token.EOL, self.line, self.col)
	flush(token.EOF, self.line, self.col)
	return tokens, nil
}
