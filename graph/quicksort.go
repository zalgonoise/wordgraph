package graph

// getterFunc is a simplified type to represent an abstract getter function that
// will return a certain element's value, from a Result
type getterFunc = func(r *Result) int

// metricValue is a type to represent the options (different metrics) for creating
// getterFunc
type metricValue int

const (
	weightV    metricValue = 0 // placeholder enum entry for weight value
	potentialV metricValue = 1 // placeholder enum entry for potentual value
)

// metricsValues is a map connecting different metricValue to their corresponding
// getterFunc.
//
// The enum values above are mapped to their own getterFunc, so that fetching
// these functions is fast, linear, non-complex, and straight-forward (returns nil
// if it doesn't exist).
var metricValues = map[metricValue]getterFunc{
	weightV: func(r *Result) int {
		return r.weight
	},
	potentialV: func(r *Result) int {
		return r.potential
	},
}

// f method is a private method to return a metricValue's getterFunc, from the metricValues
// map. Returns nil if non-existent.
func (m metricValue) f() getterFunc {
	return metricValues[m]
}

// swap function will return the input slice, after swapping its indexes a and b
func swap(slice []*Result, a, b int) []*Result {
	slice[a], slice[b] = slice[b], slice[a]
	return slice
}

// quickSortValue function will apply a simple quicksort algorithm to order the items in a
// list of results, with the comparison being done from the value retrieved with the
// getterFunc.
//
// So, while the Result struct contains a weight and a potential value, this logic is applied
// for either with the appropriate getterFunc.
func quickSortValue(r []*Result, f getterFunc) []*Result {

	// if getterFunc is nil, no comparison can be made.
	// return original slice for safety
	if f == nil {
		return r
	}

	// short-circuit if slice is zero or one items long
	if len(r) < 2 {
		return r
	}

	// set boundaries
	left, right := 0, len(r)-1

	// set pivot to the approximate middle of the slice
	pivot := len(r) / 2

	// swap right edge and pivot
	r = swap(r, pivot, right)

	// swap elements bigger than the pivot to the left
	// (for a descending sort)
	for i := range r {
		if f(r[i]) > f(r[right]) {
			r = swap(r, left, i)
			left++
		}
	}

	// swap pivot and biggest element
	r = swap(r, left, right)

	// rinse and repeat to the left and to the right
	quickSortValue(r[:left], f)
	quickSortValue(r[left+1:], f)

	return r
}

// splitByWeight function will break-down the (weight-ordered) list of results, and separating them
// into different slices (range of weights is smaller than range of potential); returning a slice of slices
// of Result pointers, as ordered blocks, separated by different weights.
func splitByWeight(r []*Result) [][]*Result {
	out := [][]*Result{} // initialize output
	last := 0            // initialize first (incrementing) index
	ref := map[int]int{} // initialize index-to-weight map

	// iterate through all results
	for _, result := range r {

		// if this is the first entry, append it to the output list,
		// taking note of its weight in the reference map
		if len(out) == 0 {
			inner := []*Result{result}
			out = append(out, inner)
			ref[last] = result.weight
			continue
		}

		if result.weight == out[last][0].weight {
			// if the result's weight matches the last entry's weight, append it to the
			// last slice
			out[last] = append(out[last], result)

		} else {
			// otherwise, increment last, add a new slice to the output slice of slices,
			// append this slice to it, and add this as a new reference in the map
			last++
			inner := []*Result{result}
			out = append(out, inner)
			ref[last] = result.weight
		}
	}

	return out
}

// quickSort function will order the input list of results by weight, then by potential. This will
// return a list that prioritizes those results which have the most weight, and within those, the
// ones that have the most potential are listed first.
//
// This is achieved by quicksorting by weight, splitting into different slices (per weight), and sorting
// those slices by potential. Lastly, all results are aggregated and return as a slice of Result pointers.
func quickSort(r []*Result) []*Result {
	var out [][]*Result
	var res []*Result

	// quicksort results by weight
	weighSorted := quickSortValue(r, weightV.f())

	// split results by weight
	split := splitByWeight(weighSorted)

	// for slice, quicksort by potential, then append it to the output
	for _, s := range split {
		ordered := quickSortValue(s, potentialV.f())
		out = append(out, ordered)
	}

	// merge output
	for _, o := range out {
		res = append(res, o...)
	}

	return res
}
