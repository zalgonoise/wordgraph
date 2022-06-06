package graph

import (
	"errors"
	"time"
)

var (
	ErrNonExistent error = errors.New("word does not exist")
	ErrNoMatches   error = errors.New("no matches found")
	ErrNoRoute     error = errors.New("no route to target")
	ErrSameWord    error = errors.New("origin and target words can't be the same")
)

const (
	maxQueryTime = time.Second * 20
	maxRoutes    = 5
	minAccuracy  = 98
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

	if err != nil {
		return nil, err
	}

	return n.burstRouter(origin, target, r), nil
}

func (n *Node) burstRouter(origin, target string, siblings []*Result) []string {
	done := make(chan struct{})
	res := make(chan []string)
	out := make(chan []string)

	for _, s := range siblings {

		carry := []string{origin, s.w}

		if s.weight >= minAccuracy {

			done <- struct{}{}
			return carry
		}

		go n.findRoute(s.w, target, carry, done, res)

		go n.findBestRoute(res, done, out)
	}

	for {
		select {
		case result := <-out:
			return result
		}
	}
}

func (n *Node) findBestRoute(rCh chan []string, done chan struct{}, out chan []string) {
	size := map[int]int{}
	routes := [][]string{}
	var smallest int

	go func() {
		time.Sleep(maxQueryTime)
		done <- struct{}{}
		if len(routes) == 0 {
			out <- []string{}
			return
		}
		out <- routes[smallest]
		return
	}()

	for {
		select {
		case route := <-rCh:
			if len(routes) > maxRoutes {
				done <- struct{}{}
				out <- routes[smallest]
				return
			}

			if len(size) == 0 {
				smallest = 0
				size[0] = len(route)
				routes = append(routes, route)
			}

			if len(route) > 0 && len(route) < size[smallest] {
				smallest++
				size[smallest] = len(route)
				routes = append(routes, route)
			}

		}
	}
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

	if len(carry) >= len(origin)*3 {
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
			var exists bool

			for _, carryObj := range carry {
				if sibling.w == carryObj {
					exists = true
					break
				}
			}
			if exists {
				continue
			}

			carry = append(carry, sibling.w)

			if sibling.weight >= minAccuracy {
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
