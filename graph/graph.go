package graph

import (
	"strings"
)

// Node struct is a graph data structure that is non-cyclical, uni-directional, unweighted and map-based
//
// With this approach it is intended for the app to load the entire dictionary into memory, however it will
// allow quick operations against it
//
// Despite being naturally unweighted, the Node struct will contain methods which will analyze the results
// and provide objects (Result) that provide insight (and weight values) on the output for a certain query.
//
// The Node struct will contain the following elements:
//   - charMap map[byte]*Node; this element is a matrix of all the letters in the alphabet, which are linked to
//   a pointer to a Node. This represents the next character in a word. For the English language, this does not
//   require a different format such as a rune map (instead of byte); and this implementation is using bytes for
//   simplicity reasons.
//
//   - char byte; a quick reference to the current character this Node represents. It can be zero if it's the root;
//   however for that check, the parent element is used instead.
//
//   - isEnd bool; is a placeholder delimiter for end-of-words. Just because isEnd == true, may not mean that there
//   are no more nodes in the graph (bi[g] marks the end of the word "big", but the graph can continue to the word
//   "bi[g]ger")
//
//   - parent *Node; points to the node above this one. The root node will have this element set to nil, and to get to
//   the root node, it's only needed to traverse this pointer until it's nil.
type Node struct {
	charMap map[byte]*Node
	char    byte
	isEnd   bool
	parent  *Node
}

// New function will create a new Node pointer, with an already initialized charMap.
func New() *Node {
	return &Node{charMap: map[byte]*Node{}}
}

// getRoot method is private, and is used by (mostly public) methods to get to the top-level node, or root node.
//
// This action is essential to perform new queries, in case the user is accessing a method from a child Node. It
// will traverse the node's parent element (a pointer to the Node above it) until it is nil.
func (n *Node) getRoot() *Node {
	if n.parent != nil {
		return n.parent.getRoot()
	}
	return n
}

// Add method will be variadic, taking any number of strings to populate the graph / dictionary.
//
// It does so by traversing to the root Node, then calling the `rAdd()` private method on each input word, recursively.
// The `rAdd()` call will focus on a word character-by-character, by checking if it already exists in the matrix. If
// it doesn't, it creates a new node for it, with its character and pointing to its parent, until the last character.
//
// If the node is already present, it will recursively call `rAdd()` on that child node, popping one character from
// the beginning of the word until it's stored.
func (n *Node) Add(word ...string) {
	// short-circuit if the input is empty
	if len(word) == 0 {
		return
	}

	// ensure the call is done on the root node
	node := n.getRoot()

	// iterate through all words, calling `rAdd()` on each
	for _, w := range word {
		node.rAdd(w)
	}
}

// rAdd method will recursively add a word to the graph.
//
// It does so by checking if the first letter of the word is present in the graph (starting from the root); creating
// a new pointer of a Node if it isn't, or cascading down to it if it exists.
//
// Cascading down is simply calling the same method but on the appropriate child node while popping the first character
// of the word.
//
// When the word reaches its last character, the node's `isEnd` element is set to true.
func (n *Node) rAdd(word string) {
	// shor-circuit on zero-length input
	if len(word) == 0 {
		return
	}

	// take the first character in the word
	char := word[0]

	// if it doesn't exist in the map; populate it with a new pointer
	if n.charMap[char] == nil {

		n.charMap[char] = &Node{
			charMap: map[byte]*Node{},
			char:    char,
			parent:  n,
		}

		// if this is the last character, set node's isEnd as true
		if len(word) == 1 {
			n.charMap[char].isEnd = true
			return
		}

		// continue until input is empty
		n.charMap[char].rAdd(word[1:])
		return
	}

	// continue until input is empty
	n.charMap[char].rAdd(word[1:])
}

// Byte method returns the node's representative character, in bytes
func (n *Node) Byte() byte {
	return n.char
}

// String method returns the node's representative character.
func (n *Node) String() string {
	return string(n.char)
}

// Rune method returns the node's representative character, as a rune
func (n *Node) Rune() rune {
	return rune(n.char)
}

// Print method will pretty-print a list-like relationship between all elements of this node and
// its children, in a text format. Not advisable to exectute for larger-sized dictionaries.
//
// Such as:
//
//     [c]
//      +-[a]
//         +-[t]
//         +-[r]
//
func (n *Node) Print(sep, brk string) string {
	var sb = &strings.Builder{}

	for b, nd := range n.charMap {

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
