package graph

// Fuzz method will take an input word and alter its characters while trying to find real words doing so
// (by exploring populated pointers in the dictionary).
//
// All existing words will be aggregated into a slice and returned
func (n *Node) Fuzz(word string) ([]string, error) {
	matches := []string{}

	// get all nodes for each character in the input word
	nodes := n.rGetNodes(word)

	// for each character's nodes, run the fuzzer func
	for idx, n := range nodes {
		matches = append(matches, fuzz(word, idx, n)...)    // fuzz the word's contents
		matches = append(matches, expanded(word, nodes)...) // explore if word can be expanded (+1 characters)
		matches = append(matches, reduced(word, nodes)...)  // explore if word can be reduced (-1 characters)
		matches = trimDuplicates(matches)                   // trim duplicates

	}

	if len(matches) == 0 {
		return nil, ErrNoMatches
	}

	return matches, nil
}

func (n *Node) WeighedFuzz(word, target string) ([]*Result, error) {
	// fuzz the input word
	m, err := n.Fuzz(word)

	if err != nil {
		return nil, err
	}

	out := []*Result{}

	// for each match, fetch the siblings, and use the target, match and siblings
	// to create a new Result entry
	for _, match := range m {
		siblings, err := n.Siblings(match)

		if err != nil {
			continue
		}

		out = append(out, newResult(target, match, siblings))
	}

	return out, nil
}

// fuzz function will take an origin word, an index, and a target Node, to explore possible existing words
// simply by altering a character in the sequence.
//
// This is done by exploring all non-nil pointers of the parent (except for the character being swapped),
// and then calling the Find() method to ensure it is a real word.
func fuzz(word string, idx int, node *Node) []string {
	matches := []string{}

	// get the parent
	parent := node.parent

	// scramble all keys
	for key, val := range parent.charMap {

		// ignore the same key for this run
		if key == word[idx] {
			continue
		}

		// change the character
		new := []byte(word)
		new[idx] = val.Byte()

		// look it up
		if node.Find(string(new)) {
			matches = append(matches, string(new))
		}
	}

	return matches
}

// expanded function will be similar to fuzz, but instead of looking up the parent, it will look up the provided
// pointer for non-nil characters, in nodes that also have set `isEnd` to true.
func expanded(word string, nodes []*Node) []string {
	var out []string

	for _, node := range nodes[len(nodes)-1].charMap {
		if node.isEnd {
			new := []byte(word)
			new = append(new, node.char)

			out = append(out, string(new))
		}
	}

	return out
}

// reduced function will be similar to expanded, but the other way around. This is done by ensuring that the
// nodes are at least two elements long, as it will look up the charMap of the next-to-last node.
//
// If this next-to-last node has `isEnd` set to true, it can be reduced and counts as a match.
func reduced(word string, nodes []*Node) []string {
	var out []string

	if len(nodes) >= 2 {
		for _, node := range nodes[len(nodes)-2].charMap {
			if node.isEnd {
				new := []byte(word[:len(word)-1])
				out = append(out, string(new))
			}
		}
	}

	return out
}

// trimDuplicates function will leverage a simple, memory-efficient data structure (map of something to an
// empty struct), to ensure that no repeated entries are returned.
//
// This is done by populating a map with each key, and checking if that same key has been initialized before
// on a subsequent word (with the bool / OK value taken from maps)
func trimDuplicates(slice []string) []string {
	keys := map[string]struct{}{}
	out := []string{}

	for _, word := range slice {
		if _, ok := keys[word]; ok {
			continue
		}
		keys[word] = struct{}{}

		out = append(out, word)
	}

	return out
}
