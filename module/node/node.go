package node

import (
	t "language/module/token"
	"strconv"
)

type Node struct {
	Token t.Token
	Left  *Node
	Right *Node
}

func New(token t.Token) *Node {
	return &Node{
		Token: token,
		Left:  nil,
		Right: nil,
	}
}

func NewWithChildren(token t.Token, left, right *Node) *Node {
	return &Node{
		Token: token,
		Left:  left,
		Right: right,
	}
}

func (self Node) Value() (int, bool) {
	if self.Token.Kind != t.IntConstant {
		return 0, false
	}

	i, err := strconv.ParseUint(string(self.Token.Data), 10, 64)
	if err != nil {
		return 0, false
	}

	return int(i), true
}

func (self Node) IsTerminal() bool {
	return self.Left == nil && self.Right == nil
}

/*
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
*/
