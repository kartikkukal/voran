package token

import "fmt"

type Value string

const (
	// Keywords
	I8        Value = "i8"
	I16       Value = "i16"
	I32       Value = "i32"
	I64       Value = "i64"
	U8        Value = "u8"
	U16       Value = "u16"
	U32       Value = "u32"
	U64       Value = "u64"
	If        Value = "if"
	Else      Value = "else"
	While     Value = "while"
	Interrupt Value = "interrupt"

	// Operators
	Equals        Value = "="
	Add           Value = "+"
	AddEqual      Value = "+="
	Subtract      Value = "-"
	SubtractEqual Value = "-="
	Multiply      Value = "*"
	MultiplyEqual Value = "*="
	Divide        Value = "/"
	DivideEqual   Value = "/="
	Modulus       Value = "%"
	ModulusEqual  Value = "%="

	Equal          Value = "=="
	NotEqual       Value = "!="
	LessThan       Value = "<"
	LessOrEqual    Value = "<="
	GreaterThan    Value = ">"
	GreaterOrEqual Value = ">="

	// Punctuators
	ParenL    Value = "("
	ParenR    Value = ")"
	CurlyL    Value = "{"
	CurlyR    Value = "}"
	SquareL   Value = "["
	SquareR   Value = "]"
	Semicolon Value = ";"
	Colon     Value = ":"
	Comma     Value = ","
	Period    Value = "."

	Empty Value = ""
)

type Type string

// Types of tokens
const (
	Identifier Type = "identifier"
	Keyword    Type = "keyword"
	Operator   Type = "operator"
	Punctuator Type = "punctuator"
	String     Type = "string"
	Integer    Type = "integer"

	EOL Type = "eol"
	EOF Type = "eof"
)

type Token struct {
	Value Value
	Type  Type

	line   int
	column int
}

func New(data Value, t Type, l, c int) Token {
	return Token{
		Value:  data,
		Type:   t,
		line:   l,
		column: c,
	}
}

func (self Token) Debug() string {
	return fmt.Sprintf("'%v': %v at line %v, column %v", self.Value, self.Type, self.line, self.column)
}
