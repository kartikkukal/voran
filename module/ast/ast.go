package ast

const (
	Add = iota
	Subtract
	Multiply
	Divide
	Modulus

	LessThan
	LessOrEqual
	GreaterThan
	GreaterOrEqual

	Glue
	IfElse
	While

	Inter

	IsEqual
	NotEqual

	DeclareInt

	LiteralInt
	LiteralStr
	Identifier

	Equal
)

type Node struct {
	Type  int
	Value int

	RValue bool

	Left  *Node
	Right *Node
}

func New(t int) *Node {
	return &Node{
		Type:   t,
		Value:  0,
		RValue: true,
		Left:   nil,
		Right:  nil,
	}
}

func NewWithValue(t, value int) *Node {
	return &Node{
		Type:   t,
		Value:  value,
		RValue: true,
		Left:   nil,
		Right:  nil,
	}
}

func NewWithChildren(t int, left, right *Node) *Node {
	return &Node{
		Type:  t,
		Value: 0,
		Left:  left,
		Right: right,
	}
}

func (self *Node) Insert(node *Node) bool {
	if self.Left == nil {
		self.Left = node
		return true
	}

	if self.Right == nil {
		self.Right = node
		return true
	}

	return false
}
