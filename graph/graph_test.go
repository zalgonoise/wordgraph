package graph

import (
	"reflect"
	"regexp"
	"testing"
)

func FuzzAddAndFind(f *testing.F) {
	module := "Graph"
	funcname := "Add() >> Find()"
	action := "fuzz testing graph's Add(...string), then Find(string)"

	f.Add("cat")
	f.Add("cot")
	f.Add("dot")
	f.Add("dog")

	f.Fuzz(func(t *testing.T, a string) {
		n := New()
		n.Add(a)

		if !n.Find(a) {
			t.Errorf(
				"FAILED -- [%s] [%s] unable to find the added word %s in the graph -- action: %s",
				module,
				funcname,
				a,
				action,
			)
			return

		}
	})
}

func TestAdd(t *testing.T) {
	module := "Graph"
	funcname := "Add() >> Find()"

	_ = module
	_ = funcname

	type test struct {
		name  string
		input []string
	}

	var tests = []test{
		{
			name:  "single word input",
			input: []string{"cat"},
		},
		{
			name:  "multi word input",
			input: []string{"cat", "cot", "dot", "dog"},
		},
		{
			name:  "zero-length input",
			input: []string{},
		},
		{
			name:  "nil input",
			input: nil,
		},
	}

	var verify = func(idx int, test test) {
		n := New()
		n.Add(test.input...)

		if len(test.input) > 0 {
			for _, word := range test.input {
				if !n.Find(word) {
					t.Errorf(
						"#%v -- FAILED -- [%s] [%s] unable to find the added word %s in the graph -- action: %s",
						idx,
						module,
						funcname,
						word,
						test.name,
					)
					return
				}
			}
		}

	}

	for idx, test := range tests {
		verify(idx, test)
	}
}

func verifyByte(word string, node *Node) bool {
	if len(word) == 0 {
		return true
	}

	n, ok := node.charMap[word[0]]

	if !ok {
		return false
	}

	if n.Byte() != word[0] {
		return false
	}

	return verifyByte(word[1:], n)
}

func FuzzByte(f *testing.F) {
	module := "Graph"
	funcname := "Byte()"
	action := "fuzz testing graph's Byte() call in its Nodes"

	f.Add("cat")
	f.Add("cot")
	f.Add("dot")
	f.Add("dog")

	f.Fuzz(func(t *testing.T, a string) {
		n := New()
		n.Add(a)

		if ok := verifyByte(a, n); !ok {
			t.Errorf(
				"FAILED -- [%s] [%s] invalid output when calling Byte() on a node, from input %s -- action: %s",
				module,
				funcname,
				a,
				action,
			)
			return
		}
	})
}

func verifyString(word string, node *Node) bool {
	if len(word) == 0 {
		return true
	}

	n, ok := node.charMap[word[0]]

	if !ok {
		return false
	}

	if n.String() != string(word[0]) {
		return false
	}

	return verifyByte(word[1:], n)
}

func FuzzString(f *testing.F) {
	module := "Graph"
	funcname := "String()"
	action := "fuzz testing graph's String() call in its Nodes"

	f.Add("cat")
	f.Add("cot")
	f.Add("dot")
	f.Add("dog")

	f.Fuzz(func(t *testing.T, a string) {
		n := New()
		n.Add(a)

		if ok := verifyString(a, n); !ok {
			t.Errorf(
				"FAILED -- [%s] [%s] invalid output when calling String() on a node, from input %s -- action: %s",
				module,
				funcname,
				a,
				action,
			)
			return
		}
	})
}

func verifyRune(word string, node *Node) bool {
	if len(word) == 0 {
		return true
	}

	n, ok := node.charMap[word[0]]

	if !ok {
		return false
	}

	if n.Rune() != rune(word[0]) {
		return false
	}

	return verifyByte(word[1:], n)
}

func FuzzRune(f *testing.F) {
	module := "Graph"
	funcname := "Rune()"
	action := "fuzz testing graph's Rune() call in its Nodes"

	f.Add("cat")
	f.Add("cot")
	f.Add("dot")
	f.Add("dog")

	f.Fuzz(func(t *testing.T, a string) {
		n := New()
		n.Add(a)

		if ok := verifyRune(a, n); !ok {
			t.Errorf(
				"FAILED -- [%s] [%s] invalid output when calling Rune() on a node, from input %s -- action: %s",
				module,
				funcname,
				a,
				action,
			)
			return
		}
	})
}

func TestPrint(t *testing.T) {
	module := "Graph"
	funcname := "Add() >> Print()"

	_ = module
	_ = funcname

	type test struct {
		name  string
		input []string
		wants string
	}

	var tests = []test{
		{
			name:  "single word input",
			input: []string{"cat"},
			wants: `\[c\]\n\s*\+-\[a\]\n\s*\+-\[t\]\s*;\s*`,
		},
		{
			name:  "multi word input",
			input: []string{"cat", "car", "cab"},
			wants: `\[c\]\n\s*\+-\[a\](\n\s*\+-\[(t|r|b)\]\s*;\s*){3}`,
		},
	}

	var verify = func(idx int, test test) {
		n := New()
		n.Add(test.input...)

		result := n.Print("", "")

		rgx := regexp.MustCompile(test.wants)

		if !rgx.MatchString(result) {
			t.Errorf(
				"#%v -- FAILED -- [%s] [%s] unexpected output from provided keywords -- input %v ; \n-- wanted: \n%s\n-- got \n%s\n -- action: %s",
				idx,
				module,
				funcname,
				test.input,
				test.wants,
				result,
				test.name,
			)
			return
		}

	}

	for idx, test := range tests {
		verify(idx, test)
	}
}

func FuzzGetRoot(f *testing.F) {
	module := "Graph"
	funcname := "getRoot()"
	action := "fuzz testing graph's Add(...string), then getRoot() on each node"

	f.Add("cat")
	f.Add("cot")
	f.Add("dot")
	f.Add("dog")

	f.Fuzz(func(t *testing.T, a string) {
		n := New()
		n.Add(a)

		nodes := n.GetNodes(a)

		for _, node := range nodes {
			root := node.getRoot()

			if !reflect.DeepEqual(root, n) {
				t.Errorf(
					"FAILED -- [%s] [%s] output mismatch error when getting to the root from different nodes; input: %s -- action: %s",
					module,
					funcname,
					a,
					action,
				)
				return
			}
		}
	})
}
