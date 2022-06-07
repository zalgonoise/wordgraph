package graph

import "time"

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
