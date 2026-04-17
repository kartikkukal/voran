package lexer

import (
	"errors"
	"fmt"
	"slices"

	t "language/module/token"
)

var (
	ErrUnknownCharacter = errors.New("unknown character encountered")
	ErrUnknownToken     = errors.New("unknown token encountered")

	keywords = []string{
		"auto", "break", "case", "char", "const", "continue", "default",
		"do", "double", "else", "enum", "extern", "float", "for", "goto",
		"if", "int", "long", "register", "return", "short", "signed", "sizeof",
		"static", "struct", "switch", "typedef", "union", "unsigned", "void",
		"volatile", "while",
	}

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
	char int
}

func New(source []byte) Lexer {
	return Lexer{
		source: source,
		index:  0,

		line: 1,
		char: 0,
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
		self.char++
	}()

	return self.peek()
}

func (self *Lexer) advance() {
	self.line += 1
	self.char = 0
}

func (self *Lexer) Tokenize() ([]t.Token, error) {
	tokens := make([]t.Token, 0)
	buffer := make([]byte, 0)

	wrapper := func(err error) error {
		return fmt.Errorf("error encountered at line %v:%v: %w", self.line, self.char, err)
	}

	flush := func(kind int) {
		tokens = append(tokens, t.Token{
			Kind: kind,
			Data: buffer,
		})
		buffer = nil
	}

	read := func(characters []byte, kind int) {
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

		if slices.Contains(keywords, string(buffer)) {
			flush(t.Keyword)
			return
		}

		flush(kind)
	}

	delimiter := func(kind int, terminal byte) {
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

		flush(kind)
	}

	for {
		c, ok := self.consume()
		if !ok {
			break
		}

		if slices.Contains(whitespace, c) {
			continue
		}

		if c == newline {
			self.advance()
			tokens = append(tokens, t.Single(t.EOL))

			continue
		}

		if slices.Contains(identifierStart, c) {

			buffer = append(buffer, c)
			read(identifier, t.Identifier)

			continue
		}

		if slices.Contains(operators, c) {

			k, ok := t.IdentifySingle(c)
			if !ok {
				return nil, wrapper(ErrUnknownCharacter)
			}

			next, ok := self.peek()
			if !ok {
				break
			}

			kNext, ok := t.IdentifySingle(next)
			if kNext == t.Equal && ok {

				k, ok := t.IdentifyAugmented(k)
				if !ok {
					return nil, wrapper(ErrUnknownToken)
				}

				tokens = append(tokens, t.Single(k))

				_, ok = self.consume()
				if !ok {
					break
				}

				continue
			}

			tokens = append(tokens, t.Single(k))
			continue
		}

		if slices.Contains(punctuators, c) {

			k, ok := t.IdentifySingle(c)
			if !ok {
				return nil, wrapper(ErrUnknownCharacter)
			}

			tokens = append(tokens, t.Single(k))
			continue
		}

		if slices.Contains(decimal, c) {
			buffer = append(buffer, c)
			read(decimal, t.IntConstant)

			continue
		}

		if c == stringDelimiter {
			delimiter(t.StrConstant, stringDelimiter)
			continue
		}

		if c == charDelimiter {
			delimiter(t.ChrConstant, charDelimiter)
			continue
		}

		return nil, wrapper(ErrUnknownCharacter)
	}

	tokens = append(tokens, t.Single(t.EOF))
	return tokens, nil
}
