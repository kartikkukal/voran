package token

import (
	"fmt"
)

const (
	// Keywords
	Byte   = "byte"
	Short  = "short"
	Int    = "int"
	Long   = "long"
	UByte  = "ubyte"
	UShort = "ushort"
	UInt   = "uint"
	ULong  = "ulong"
	If     = "if"
	Else   = "else"
	While  = "while"

	// Operators
	Equal = "="
	Add   = "+"
	AddEq = "+="
	Sub   = "-"
	SubEq = "-="
	Mul   = "*"
	MulEq = "*="
	Div   = "/"
	DivEq = "/="
	Mod   = "%"
	ModEq = "%="

	Equals    = "=="
	NotEq     = "!="
	Less      = "<"
	LessEq    = "<="
	Greater   = ">"
	GreaterEq = ">="

	// Punctuators
	ParenL    = "("
	ParenR    = ")"
	CurlyL    = "{"
	CurlyR    = "}"
	SquareL   = "["
	SquareR   = "]"
	Semicolon = ";"
	Colon     = ":"
	Comma     = ","
	Period    = "."

	Empty = ""
)

type Type uint8

// Types of tokens
const (
	Identifier = iota
	Keyword
	Operator
	Punctuator

	LiteralStr
	LiteralInt

	EOL
	EOF
)

var (
	Keywords = []string{
		Byte,
		Short,
		Int,
		Long,
		UByte,
		UShort,
		UInt,
		ULong,
		If,
		Else,
		While,
	}

	Operators = []string{
		Equal,
		Add,
		AddEq,
		Sub,
		SubEq,
		Mul,
		MulEq,
		Div,
		DivEq,
		Mod,
		ModEq,
		Equals,
		NotEq,
		Less,
		LessEq,
		Greater,
		GreaterEq,
	}
)

type Token struct {
	Value string
	Type  Type

	line   int
	column int
}

func New(data string, t Type, line, col int) Token {
	return Token{
		Value:  data,
		Type:   t,
		line:   line,
		column: col,
	}
}

func TypeOf(t string) Type {

	for _, v := range Keywords {
		if t == v {
			return Keyword
		}
	}

	for _, v := range Operators {
		if t == v {
			return Operator
		}
	}

	return Punctuator
}

func (self Token) Debug() string {
	return fmt.Sprintf("'%v': %v at line %v, column %v", self.Value, self.Type, self.line, self.column)
}
