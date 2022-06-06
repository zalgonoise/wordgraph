package graph

import (
	"errors"
	"fmt"
)

var (
	ErrNonExistent error = errors.New("word does not exist")
	ErrNoMatches   error = errors.New("no matches found")
	ErrNoRoute     error = errors.New("no route to target")
	ErrSameWord    error = errors.New("origin and target words can't be the same")
)

func (n *Node) Find(w string) bool {
	node := n.getRoot()

	return node.find(w)
}

func (n *Node) find(w string) bool {
	if len(w) == 0 {
		return false
	}

	b := w[0]

	if n.c[b] == nil {
		return false
	}

	if n.c[b].isEnd && len(w) == 1 {
		return true
	}

	return n.c[b].find(w[1:])
}

func (n *Node) GetNodes(w string) []*Node {
	if !n.Find(w) {
		return []*Node{}
	}

	return n.getNodes(w)
}

func (n *Node) getNodes(w string) []*Node {
	var out = []*Node{}

	if len(w) == 0 {
		return out
	}

	var b = w[0]

	out = append(out, n.c[b])

	out = append(out, n.c[b].getNodes(w[1:])...)

	return out
}

func (n *Node) Siblings(origin string) ([]string, error) {
	node := n.getRoot()

	if !node.Find(origin) {
		return nil, ErrNonExistent
	}

	return node.Fuzz(origin)
}

func (n *Node) TargetSiblings(origin, target string) ([]*Result, error) {
	node := n.getRoot()

	if !node.Find(origin) {
		return nil, ErrNonExistent
	}

	weighed, err := node.WeighedFuzz(origin, target)

	if err != nil {
		return nil, err
	}

	return quickSort(weighed), nil
}

func (n *Node) FindRoute(origin, target string) ([]string, error) {
	if origin == target {
		return nil, ErrSameWord
	}

	if !n.Find(target) {
		return nil, ErrNonExistent
	}

	r, err := n.TargetSiblings(origin, target)

	fmt.Println(r, "\n\n")

	if err != nil {
		return nil, err
	}

	return n.burstRouter(origin, target, r), nil
}

func (n *Node) burstRouter(origin, target string, siblings []*Result) []string {
	done := make(chan struct{})
	res := make(chan []string)

	out := [][]string{}

	for _, s := range siblings {
		fmt.Println(s.w)

		carry := []string{s.w}

		if s.weight == 100 {
			fmt.Println("done\n\n")

			done <- struct{}{}
			return carry
		}

		go n.findRoute(s.w, target, carry, done, res)
	}

	for {
		select {

		case route := <-res:
			fmt.Println("\n--\nroute:", route)

			if len(out) < 50 {
				fmt.Println("appended; list: ", out)
				out = append(out, route)
			} else {
				return n.findBestRoute(out)
			}
		}
	}
}

func (n *Node) findBestRoute(routes [][]string) []string {
	size := map[int]int{}
	var smallest int

	for idx, r := range routes {
		if idx == 0 {
			smallest = idx
			size[idx] = len(r)
			continue
		}

		if len(r) > 0 && len(r) < size[smallest] {
			smallest = idx
			size[idx] = len(r)
		}
	}

	return routes[smallest]
}

func (n *Node) findRoute(
	origin string, target string,
	carry []string,
	done <-chan struct{},
	res chan []string,
) {
	var innerDone = make(chan struct{})

	go func() {

		select {
		case <-done:
			innerDone <- struct{}{}
			return
		}

	}()

	if len(carry) >= len(origin)*4 {
		return
	}

	r, err := n.TargetSiblings(origin, target)

	if err != nil {
		return
	}

	for _, sibling := range r {
		select {
		case <-innerDone:
			return
		default:
			carry = append(carry, sibling.w)

			if sibling.weight == 100 {
				res <- carry
				return
			}

			go n.findRoute(sibling.w, target, carry, done, res)
		}
	}

}

func fuzz(o string, idx int, n *Node) []string {
	matches := []string{}

	// get the parent
	parent := n.parent

	// scramble all keys
	for k, v := range parent.c {

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

	for _, node := range n[len(n)-1].c {
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
		for _, node := range n[len(n)-2].c {
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

func (n *Node) Fuzz(o string) ([]string, error) {
	matches := []string{}
	nodes := n.getNodes(o)

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
	nodes := n.getNodes(o)

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

// func (n *Node) siblings(o string, s int) ([]string, error) {

// }
