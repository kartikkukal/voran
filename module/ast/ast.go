package ast

import (
	"fmt"
)

type Kind uint8

//go:generate stringer -type=Kind
const (
	Add Kind = iota
	Sub
	Mul
	Div
	Mod
	Equal
	Equals
	NotEq
	Less
	LessEq
	Greater
	GreaterEq

	Glue
	Declare
	Branch
	While
	Function
	Parameter
	Identifier
	Sycall
	Call

	Implicit

	Literal

	Invalid
)

type Type uint8

//go:generate stringer -type=Type
const (
	Byte Type = iota
	Short
	Int
	Long
	UByte
	UShort
	UInt
	ULong

	Float
	Double

	String
	None
)

type Node struct {
	Kind  Kind
	Type  Type
	Value int

	RValue bool

	Left  *Node
	Right *Node
}

func New(t Kind) *Node {
	return &Node{
		Kind:   t,
		Type:   None,
		Value:  0,
		RValue: true,
		Left:   nil,
		Right:  nil,
	}
}

func NewWithValue(kind Kind, value int, t Type) *Node {
	return &Node{
		Kind:   kind,
		Value:  value,
		Type:   t,
		RValue: true,
		Left:   nil,
		Right:  nil,
	}
}

func NewWithChildren(t Kind, left, right *Node) *Node {
	return &Node{
		Kind:  t,
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

func (self *Node) Terminal() bool {
	return self.Left == nil && self.Right == nil
}

func (self *Node) debug(indent string, last bool, middle bool) string {

	var output string

	if middle {
		output += fmt.Sprintf("%v├─ %v\n", indent, self.Kind)
	} else {
		output += fmt.Sprintf("%v└─ %v\n", indent, self.Kind)

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
