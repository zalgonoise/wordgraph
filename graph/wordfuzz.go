package graph

func (n *Node) Fuzz(o string) ([]string, error) {
	matches := []string{}
	nodes := n.rGetNodes(o)

	for idx, n := range nodes {
		matches = append(matches, fuzz(o, idx, n)...)
		matches = append(matches, expanded(o, nodes)...)
		matches = append(matches, reduced(o, nodes)...)
		matches = trimDuplicates(matches)

	}

	if len(matches) == 0 {
		return nil, ErrNoMatches
	}

	return matches, nil
}

func (n *Node) WeighedFuzz(o, t string) ([]*Result, error) {
	matches := []string{}
	nodes := n.rGetNodes(o)

	for idx, n := range nodes {
		matches = append(matches, fuzz(o, idx, n)...)
		matches = append(matches, expanded(o, nodes)...)
		matches = append(matches, reduced(o, nodes)...)
		matches = trimDuplicates(matches)

	}

	if len(matches) == 0 {
		return nil, ErrNoMatches
	}

	out := []*Result{}

	for _, m := range matches {
		siblings, err := n.Siblings(m)

		if err != nil {
			continue
		}

		out = append(out, newResult(t, m, siblings))
	}

	return out, nil
}

func fuzz(o string, idx int, n *Node) []string {
	matches := []string{}

	// get the parent
	parent := n.parent

	// scramble all keys
	for k, v := range parent.charMap {

		// ignore the same key for this run
		if k == o[idx] {
			continue
		}

		// change the character
		new := []byte(o)
		new[idx] = v.Byte()

		// look it up
		if n.getRoot().Find(string(new)) {
			matches = append(matches, string(new))
		}
	}

	return matches
}

func expanded(o string, n []*Node) []string {
	var out []string

	for _, node := range n[len(n)-1].charMap {
		if node.isEnd {
			new := []byte(o)
			new = append(new, node.char)

			out = append(out, string(new))
		}
	}

	return out
}

func reduced(o string, n []*Node) []string {
	var out []string

	if len(n) > 2 {
		for _, node := range n[len(n)-2].charMap {
			if node.isEnd {
				new := []byte(o[:len(o)-1])
				out = append(out, string(new))
			}
		}
	}

	return out
}

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
