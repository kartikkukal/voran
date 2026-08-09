package ast

import (
	"fmt"
)

type Type string

const (
	Add            Type = "+"
	Subtract       Type = "-"
	Multiply       Type = "*"
	Divide         Type = "/"
	Modulus        Type = "%"
	Equals         Type = "="
	Equal          Type = "=="
	NotEqual       Type = "!="
	LessThan       Type = "<"
	LessOrEqual    Type = "<="
	GreaterThan    Type = ">"
	GreaterOrEqual Type = ">="

	Glue      Type = "glue"
	Declare   Type = "declare"
	If        Type = "if"
	While     Type = "while"
	Function  Type = "function"
	Interrupt Type = "interrupt"
	Call      Type = "call"

	Integer    Type = "number"
	String     Type = "string"
	Identifier Type = "identifier"
)

type Node struct {
	Type  Type
	Value int

	RValue bool

	Left  *Node
	Right *Node
}

func New(t Type) *Node {
	return &Node{
		Type:   t,
		Value:  0,
		RValue: true,
		Left:   nil,
		Right:  nil,
	}
}

func NewWithValue(t Type, value int) *Node {
	return &Node{
		Type:   t,
		Value:  value,
		RValue: true,
		Left:   nil,
		Right:  nil,
	}
}

func NewWithChildren(t Type, left, right *Node) *Node {
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

func (self *Node) debug(indent string, last bool, middle bool) string {

	var output string

	if middle {
		output += fmt.Sprintf("%v├─ %v\n", indent, self.Type)
	} else {
		output += fmt.Sprintf("%v└─ %v\n", indent, self.Type)

	}

	if last {
		indent += "   "
	} else {
		indent += "│  "
	}

	if self.Left != nil {
		left := self.Left.debug(indent, self.Right == nil, self.Right != nil)
		output += fmt.Sprintf("%v", left)
	}
	if self.Right != nil {
		right := self.Right.debug(indent, self.Left != nil, self.Left == nil)
		output += fmt.Sprintf("%v", right)
	}

	return output
}

func (self *Node) Debug() string {
	return self.debug("", true, false)
}
