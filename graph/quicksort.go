package graph

func swap(slice []*Result, a, b int) []*Result {
	slice[a], slice[b] = slice[b], slice[a]
	return slice
}

func quickSortPotential(r []*Result) []*Result {
	if len(r) < 2 {
		return r
	}

	left, right := 0, len(r)-1

	pivot := len(r) / 2

	r = swap(r, pivot, right)

	for i := range r {
		if r[i].potential > r[right].potential {
			r = swap(r, left, i)
			left++
		}
	}

	r = swap(r, left, right)

	quickSortPotential(r[:left])
	quickSortPotential(r[left+1:])

	return r
}

func quickSortWeight(r []*Result) []*Result {
	if len(r) < 2 {
		return r
	}

	left, right := 0, len(r)-1

	pivot := len(r) / 2

	r = swap(r, pivot, right)

	for i := range r {
		if r[i].weight > r[right].weight {
			r = swap(r, left, i)
			left++
		}
	}

	r = swap(r, left, right)

	quickSortWeight(r[:left])
	quickSortWeight(r[left+1:])

	return r
}

func splitByWeight(r []*Result) [][]*Result {
	out := [][]*Result{}
	last := 0
	ref := map[int]int{}

	for _, result := range r {
		if len(out) == 0 {
			inner := []*Result{result}
			out = append(out, inner)
			ref[last] = result.weight
			continue
		}

		if result.weight == out[last][0].weight {
			out[last] = append(out[last], result)
		} else {
			last++
			inner := []*Result{result}
			out = append(out, inner)
			ref[last] = result.weight
		}
	}

	return out
}

func quickSort(r []*Result) []*Result {
	var out [][]*Result
	var res []*Result

	weighSorted := quickSortWeight(r)

	split := splitByWeight(weighSorted)

	for _, s := range split {
		ordered := quickSortPotential(s)
		out = append(out, ordered)
	}

	for _, o := range out {
		res = append(res, o...)
	}

	return res
}
