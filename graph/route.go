package graph

import (
	"time"
)

// FindRoute method will take an origin and target words, and return the most efficient path from one word
// to the other, with single-character changes.
func (n *Node) FindRoute(origin, target string) ([]string, error) {
	// if the origin is the same as the target, no route needs to be found
	if origin == target {
		return nil, ErrSameWord
	}

	// if target doesn't exist, return an error
	if !n.Find(target) {
		return nil, ErrNonExistent
	}

	// get weighed results for the origin word's siblings
	r, err := n.TargetSiblings(origin, target)

	if err != nil {
		return nil, err
	}

	// call burstRoutes() to fire-off goroutines
	return n.burstRoutes(origin, target, r), nil
}

// burstRoutes method will handle the channels and comms necessary for performing this query while
// leveraging goroutines.
//
// This method will serve as a router for generating the first goroutines, and also as a results receiver
// to return the best match when it gets this.
func (n *Node) burstRoutes(origin, target string, siblings []*Result) []string {
	done := make(chan struct{}) // done channel to signal a closure action
	res := make(chan []string)  // res channel to communicate results
	out := make(chan []string)  // out channel to communicate the final output

	// iterate through all siblings
	for _, s := range siblings {

		// initialize a carry slice with the origin and the sibling
		carry := []string{origin, s.word}

		// if there is a match already, send done signal to done channel and return results
		if s.weight >= minAccuracy {
			return carry
		}

		// kick off findRoute() and findBestRoute()
		go n.findRoute(s.word, target, carry, done, res)
	}
	go n.findBestRoute(res, out, done)

	// once all goroutines are kicked-off, wait for a results message from the output channel
	for {
		select {
		case result := <-out:
			return result
		}
	}
}

func getShortest(input [][]string) []string {
	if len(input) == 0 {
		return []string{}
	}

	var smallest int = 0
	var smallIdx int = 0

	for idx, v := range input {
		if idx == 0 {
			continue
		}

		if len(v) < smallest {
			smallest = len(v)
			smallIdx = idx
		}
	}
	return input[smallIdx]

}

// findBestRoute method takes in a results and output channel (chan []string), and a done channel (chan struct{}),
// to serve as a listener to these channels.
//
// This method will set a time limit for this operation, after which is sends off the done signal to halt all queries,
// and will listen to the results channel, accumulating the received routes
//
// As it receives new routes, it places it in the slice as per their length, reserving the smallest one to be returned.
//
// Once the set maxRoutes value is achieved in its routes slice, it send the smallest to the output channel after sending
// the done signal.
func (n *Node) findBestRoute(rCh, out chan []string, done chan struct{}) {
	routes := [][]string{} // initialize a slice of slices to store the routes
	var entries int        // keep track of the number of received entries

	for {
		select {
		// safeguard to exit after a period of time
		case <-time.After(maxQueryTime):
			out <- getShortest(routes)
			return
		case <-done:
			out <- getShortest(routes)
			return
		case route := <-rCh:

			entries++
			go func(entries *int) {
				copy := *entries
				for {
					select {
					case <-time.After(maxNoResponseTime):
						if *entries == copy {
							done <- struct{}{}
						}
					}
				}
			}(&entries)

			routes = append(routes, route)

			// if length of routes exceeds the set maximum, exit by sending
			// the done signal and pushing the smallest slice to the output channel
			if len(routes) > maxRoutes {
				done <- struct{}{}
			}
		}
	}
}

// findRoute method is a recursive call to keep looking up new routes (by exploring new words in the
// same sequence).
//
// it takes in the origin string and the target string for reference. The carry slice will represent
// the current route, as it is called recursively. the done and results channel will serve as points of
// communication between the goroutines and other methods.
//
//
func (n *Node) findRoute(
	origin string, target string,
	carry []string,
	done chan struct{},
	rCh chan []string,
) {
	// kick-off a goroutine as a controller, to close this call and its children when the
	// done signal is called
	go func() {
		select {
		case <-done:
			return
		}
	}()

	// hard-limit for the carry length -- cannot exceed 3x the size of the word
	// so, "cat" to "dog" cannot take 9 or more operations to complete
	if len(carry) >= len(origin)*3 {
		return
	}

	// get weighted results
	r, err := n.TargetSiblings(origin, target)

	if err != nil {
		return
	}

	// cycle through each sibling
	for _, sibling := range r {
		var exists bool

		// check if the current sibling is already present as one of the items in the carry routes
		for _, carryObj := range carry {
			if sibling.word == carryObj {
				exists = true
				break
			}
		}

		// skip it if so
		if exists {
			continue
		}

		// append this word to the routes list
		carry = append(carry, sibling.word)

		// check if its weight counts passes as a match, if so send this carry slice to
		// the results channel, and return
		if sibling.weight >= minAccuracy {
			rCh <- carry
			return
		}

		// otherwise, keep exploring the siblings in a new goroutine, with this sibling's word
		// as the origin instead
		go n.findRoute(sibling.word, target, carry, done, rCh)

	}

}
