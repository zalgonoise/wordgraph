package graph

import (
	"strings"
)

const allChars = "abcdefghijklmnopqrstuvwxyz"

type Node struct {
	c      map[byte]*Node
	char   byte
	isEnd  bool
	parent *Node
}

func New() *Node {
	return &Node{c: map[byte]*Node{}}
}

func (n *Node) getRoot() *Node {
	if n.parent != nil {
		return n.parent.getRoot()
	}
	return n
}

func (n *Node) Insert(w ...string) {
	if len(w) == 0 {
		return
	}

	node := n.getRoot()

	for _, word := range w {
		node.insert(word)
	}
}

func (n *Node) insert(w string) {
	if len(w) == 0 {
		return
	}

	b := w[0]

	if n.c[b] == nil {
		n.c[b] = &Node{
			c:      map[byte]*Node{},
			char:   b,
			parent: n,
		}
		if len(w) == 1 {
			n.c[b].isEnd = true
		}

		n.c[b].insert(w[1:])
		return
	}
	n.c[b].insert(w[1:])
}

func (n *Node) Byte() byte {
	return n.char
}

func (n *Node) String() string {
	return string(n.char)
}

// func (n *Node) Rune() rune {
// 	return rune(n.char)
// }

func (n *Node) Print(sep, brk string) string {
	var sb = &strings.Builder{}

	for b, nd := range n.c {

		sb.WriteString(sep)
		sb.WriteString(brk)
		sb.WriteString("[")
		sb.WriteByte(b)
		sb.WriteString("]")
		if nd.isEnd {
			sb.WriteString(" ;")
		}
		sb.WriteString("\n")

		if brk == "" {
			sb.WriteString(nd.Print("", " +-"))
		} else {
			sb.WriteString(nd.Print(sep+"   ", " +-"))
		}
	}
	return sb.String()
}
