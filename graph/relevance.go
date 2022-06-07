package graph

import (
	"strconv"
	"strings"
)

// Result struct represent a (sibling) keyword, which holds a certain weight and potential. Result is used
// to better evaluate which siblings have the best odds when trying to arrive to a target word.
//
// The weight element represents how close the word is to the target (num matched characters * (100 / len(word)))
//
// The potential element represents the number of siblings the word has, for context on how many routes can this
// word take
type Result struct {
	word      string
	siblings  []string
	weight    int
	potential int
}

// String method will return a string representation of a result
func (r *Result) String() string {
	var sb = &strings.Builder{}

	sb.WriteString("{[")
	sb.WriteString(r.word)
	sb.WriteString("] weight: ")
	sb.WriteString(strconv.Itoa(r.weight))
	sb.WriteString(" potential: ")
	sb.WriteString(strconv.Itoa(r.potential))
	sb.WriteString("}")

	return sb.String()
}

// setWeight method will specify the weight of this word in comparison to the target
//
// A weight unit is 100 divided by the length of the word (gold has 4 weightUnits, each of which weighing 25)
//
// The more characters match the target, the bigger the weight
func (r *Result) setWeight(target string) {
	var weight int = 0
	var weightUnit = 100 / len(target)

	for i := 0; i < len(target); i++ {
		if len(r.word) >= len(target) && r.word[i] == target[i] {
			weight = weight + weightUnit
		}
	}

	r.weight = weight
}

// setPotential function will define the results' potential by the number of generated siblings it has
//
// The more siblings a word has, the bigger the potential in finding a quick route to the target.
func (r *Result) setPotential(target string, matches []string) {
	r.potential = len(matches)
	r.siblings = matches
}

// newResult function will take a target and a generated string, as well as the siblings to the generated string,
// and return a built Result profile of the word
func newResult(target, gen string, matches []string) *Result {
	result := &Result{
		word: gen,
	}

	result.setWeight(target)
	result.setPotential(target, matches)

	return result
}
