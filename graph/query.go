package graph

import (
	"errors"
	"time"
)

var (
	ErrNonExistent error = errors.New("word does not exist")                       // default error when a word does not exist in the dictionary
	ErrNoMatches   error = errors.New("no matches found")                          // default error when no matches are found for the query
	ErrNoRoute     error = errors.New("no route to target")                        // default error when no routes are found
	ErrSameWord    error = errors.New("origin and target words can't be the same") // default error when providing the same origin / target words
)

const (
	maxQueryTime = time.Second * 20 // timer ceiling for a FindRoute() operation
	maxRoutes    = 5                // maximum number of accumulated routes before halting the query
	minAccuracy  = 98               // minimum accuracy threshold to validate as a match
)

// Find method will look up the dictionary for string w, and return true if it exists
//
// It does so by getting to the root node, and iterating through each index recursively, with the
// private `rFind()` method; which will recursively look up if each character is initialized in the pointer map.
//
// Once it reaches the last character, it expects the node to also mark the end of a word, returning a boolean
// based on this match.
func (n *Node) Find(word string) bool {
	node := n.getRoot()

	return node.rFind(word)
}

// find method is called recursively, to ensure the input word exists in the dictionary.
//
// The input word is recursively shortened on each call as the graph is traversed, until the word is
// only one character long. Then, its `isEnd` element should be true.
func (n *Node) rFind(word string) bool {

	// input can't be empty
	if len(word) == 0 {
		return false
	}

	// take the first character
	char := word[0]

	// if this pointer is not initialized, it's not a word in this dictionary
	if n.charMap[char] == nil {
		return false
	}

	// if the isEnd element for this character is true and this is the last one left,
	// return true as there is a match
	if n.charMap[char].isEnd && len(word) == 1 {
		return true
	}

	// otherwise, there are more than one characters in the word, continue recursively
	return n.charMap[char].rFind(word[1:])
}

// GetNodes method will take in an input word and return a slice of pointers to Nodes, for
// each character of that word.
//
// It does so by recursivelly calling the `rGetNodes()` method, to populate the output slice
func (n *Node) GetNodes(word string) []*Node {
	// short-circuit if the word does not exist or is empty
	if len(word) == 0 || !n.Find(word) {
		return []*Node{}
	}

	return n.rGetNodes(word)
}

// rGetNodes method will traverse the graph recursively, and populating a slice of pointers to Nodes
// as it grabs it from the node's charMap. Then, the output is appended a new call to this method,
// taking off the first character in the word, continuously (until it's empty).
func (n *Node) rGetNodes(word string) []*Node {
	var out = []*Node{}

	// short-circuit if / when the word is / becomes zero-length
	if len(word) == 0 {
		return out
	}

	// take the first character in the word
	var char = word[0]

	// store its Node pointer in the output slice
	out = append(out, n.charMap[char])

	// recursively call this method to populate the slice with its child nodes.
	out = append(out, n.charMap[char].rGetNodes(word[1:])...)

	return out
}

// Siblings method will try to scramble the origin word's letters (one at a time), trying to find
// a real word. This is done by creating a new word, then matching if it exists in the dictionary.
//
// This call will be by default applied to the root, as its `Fuzz()` call will work with the word's
// corresponding nodes.
func (n *Node) Siblings(origin string) ([]string, error) {
	// work on the root for a complete word look-up
	node := n.getRoot()

	// return an error if the word does not exist
	if !node.Find(origin) {
		return nil, ErrNonExistent
	}

	// fuzz the words letters, checking if they are in fact words; returning a slice of all
	// one-step combinations
	return node.Fuzz(origin)
}

// TargetSiblings method will perform a call similar to `Siblings()`, but it will rank the results with
// metrics which may help achieving a quicker, better route.
//
// This is done with a `WeighedFuzz()` call, which builds a profile on each result, giving it a weight
// (a 0-100 score on the number of matched characters) and a potential (number of siblings it has).
//
// These results are sorted with a simple quicksort technique that orders by weight and by potential,
// accordingly. This ensures that a `FindRoutes()` call will prioritize the most "efficient" words.
func (n *Node) TargetSiblings(origin, target string) ([]*Result, error) {
	// work on the root for a complete word look-up
	node := n.getRoot()

	// return an error if the word does not exist
	if !node.Find(origin) {
		return nil, ErrNonExistent
	}

	// fuzz the words letters, checking if they are in fact words; returning a slice of all
	// one-step combinations; while building a profile on their relationship with the target word
	weighed, err := node.WeighedFuzz(origin, target)

	if err != nil {
		return nil, err
	}

	// return a sorted list of results, from most relevant to the least.
	return quickSort(weighed), nil
}
