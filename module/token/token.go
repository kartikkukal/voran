package token

import "fmt"

const (
	Keyword = iota
	Identifier
	StrConstant
	ChrConstant
	IntConstant

	ParenL
	ParenR
	CurlyL
	CurlyR
	SquareL
	SquareR
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
	LessThan
	GreaterThan
	Semicolon
	Colon
	Comma
	Period

	Equality
	NotEqual
	LessOrEqual
	GreaterOrEqual

	EOL
	EOF
)

type Token struct {
	Kind int
	Data []byte
}

func Single(kind int) Token {
	return Token{
		Kind: kind,
		Data: nil,
	}
}

func (self Token) Debug() {
	fmt.Print(self.Kind)
	fmt.Println(":", string(self.Data))
}

func IdentifySingle(c byte) (int, bool) {
	switch c {
	case '(':
		return ParenL, true
	case ')':
		return ParenR, true
	case '{':
		return CurlyL, true
	case '}':
		return CurlyR, true
	case '[':
		return SquareL, true
	case ']':
		return SquareR, true
	case '=':
		return Equal, true
	case '+':
		return Add, true
	case '-':
		return Subtract, true
	case '*':
		return Multiply, true
	case '/':
		return Divide, true
	case '%':
		return Modulus, true
	case '!':
		return Not, true
	case '<':
		return LessThan, true
	case '>':
		return GreaterThan, true
	case ';':
		return Semicolon, true
	case ':':
		return Colon, true
	case ',':
		return Comma, true
	case '.':
		return Period, true
	default:
		return 0, false
	}
}

func IdentifyAugmented(kind int) (int, bool) {
	switch kind {
	case Equal:
		return Equality, true
	case LessThan:
		return LessOrEqual, true
	case GreaterThan:
		return GreaterOrEqual, true
	case Not:
		return NotEqual, true
	case Add:
		return AddEqual, true
	case Subtract:
		return SubtractEqual, true
	case Multiply:
		return MultiplyEqual, true
	case Divide:
		return DivideEqual, true
	case Modulus:
		return ModulusEqual, true
	default:
		return 0, false
	}
}
