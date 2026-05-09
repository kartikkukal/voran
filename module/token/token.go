package token

import "fmt"

const (
	// Keywords
	I8 = iota
	I16
	I32
	I64
	U8
	U16
	U32
	U64

	If
	Else
	ElseIf

	While

	Inter

	Identifier

	LiteralStr
	LiteralChar
	LiteralInt

	// Operators
	Equal
	Add
	AddEqual
	Subtract
	SubtractEqual
	Multiply
	MultiplyEqual
	Divide
	DivideEqual
	Modulus
	ModulusEqual
	Not

	IsEqual
	NotEqual
	LessThan
	LessOrEqual
	GreaterThan
	GreaterOrEqual

	// Punctuators
	ParenL
	ParenR
	CurlyL
	CurlyR
	SquareL
	SquareR
	Semicolon
	Colon
	Comma
	Period

	EOL
	EOF
)

type Token struct {
	Data string

	Type int
	line int
	col  int
}

func New(data string, t, line, col int) Token {
	return Token{
		Data: data,
		Type: t,
		line: line,
		col:  col,
	}
}

func (self Token) Debug() string {
	return fmt.Sprintf("At line %v, column %v", self.line, self.col)
}

func (self Token) Complete() string {
	return fmt.Sprintf("Type: %v, Data: %v, Line: %v, Col: %v", self.Type, self.Data, self.line, self.col)
}

func Keyword(keyword string, line, col int) (Token, bool) {

	var t int

	switch keyword {
	case "u8":
		t = U8
	case "u16":
		t = U16
	case "u32":
		t = U32
	case "u64":
		t = U64
	case "i8":
		t = I8
	case "i16":
		t = I16
	case "i32":
		t = I32
	case "i64":
		t = I64
	case "if":
		t = If
	case "else":
		t = Else
	case "while":
		t = While
	case "inter":
		t = Inter
	default:
		return Token{}, false
	}

	return Token{
		Data: "",

		Type: t,
		line: line,
		col:  col,
	}, true
}

func Identify(c byte, line, col int) (Token, bool) {

	var t int

	switch c {
	case '(':
		t = ParenL
	case ')':
		t = ParenR
	case '{':
		t = CurlyL
	case '}':
		t = CurlyR
	case '[':
		t = SquareL
	case ']':
		t = SquareR
	case '=':
		t = Equal
	case '+':
		t = Add
	case '-':
		t = Subtract
	case '*':
		t = Multiply
	case '/':
		t = Divide
	case '%':
		t = Modulus
	case '!':
		t = Not
	case '<':
		t = LessThan
	case '>':
		t = GreaterThan
	case ';':
		t = Semicolon
	case ':':
		t = Colon
	case ',':
		t = Comma
	case '.':
		t = Period
	default:
		return Token{}, false
	}

	return Token{
		Data: "",

		Type: t,
		line: line,
		col:  col,
	}, true
}

func Augmented(c, line, col int) (Token, bool) {
	var t int

	switch c {
	case Equal:
		t = IsEqual
	case LessThan:
		t = LessOrEqual
	case GreaterThan:
		t = GreaterOrEqual
	case Not:
		t = NotEqual
	case Add:
		t = AddEqual
	case Subtract:
		t = SubtractEqual
	case Multiply:
		t = MultiplyEqual
	case Divide:
		t = DivideEqual
	case Modulus:
		t = ModulusEqual
	default:
		return Token{}, false
	}

	return Token{
		Data: "",

		Type: t,
		line: line,
		col:  col,
	}, true
}
