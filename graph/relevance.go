package graph

import (
	"strconv"
	"strings"
)

type Result struct {
	w         string
	siblings  []string
	weight    int
	potential int
}

func (r *Result) SetWeight(target string) {
	var weight int = 0
	var weightUnit = 100 / len(target)

	for i := 0; i < len(target); i++ {
		if len(r.w) >= len(target) && r.w[i] == target[i] {
			weight = weight + weightUnit
		}
	}

	r.weight = weight
}

func (r *Result) SetPotential(target string, matches []string) {
	r.potential = len(matches)
	r.siblings = matches
}

func (r *Result) GetWeight() int {
	return r.weight
}

func (r *Result) GetPotential() int {
	return r.potential
}

func (r *Result) String() string {
	var sb = &strings.Builder{}

	sb.WriteString("{[")
	sb.WriteString(r.w)
	sb.WriteString("] weight: ")
	sb.WriteString(strconv.Itoa(r.weight))
	sb.WriteString(" potential: ")
	sb.WriteString(strconv.Itoa(r.potential))
	sb.WriteString("}")

	return sb.String()
}

func newResult(target, gen string, matches []string) *Result {
	result := &Result{
		w: gen,
	}

	result.SetWeight(target)
	result.SetPotential(target, matches)

	return result
}

func order(in ...*Result) []*Result {
	if len(in) == 0 {
		return nil
	}

	var out = []*Result{}

	for _, r := range in {

		if len(out) == 0 {
			out = append(out, r)
			continue
		}

		if r.weight <= out[len(out)-1].weight {
			out = append(out, r)
		} else {
			new := []*Result{r}
			new = append(new, out...)
			out = new
		}
	}

	return out
}
