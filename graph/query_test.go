package graph

import (
	"errors"
	"reflect"
	"testing"
)

func TestFind(t *testing.T) {
	module := "Graph"
	funcname := "Add() >> Find()"

	_ = module
	_ = funcname

	type test struct {
		name  string
		query string
		ok    bool
	}

	root := New()
	root.Add("cat", "dog", "cot", "cog")

	var tests = []test{
		{
			name:  "valid query",
			query: "cat",
			ok:    true,
		},
		{
			name:  "zero-length query",
			query: "",
		},
		{
			name:  "query for a non-existent item",
			query: "log",
		},
	}

	var verify = func(idx int, test test) {
		if !root.Find(test.query) && test.ok {
			t.Errorf(
				"#%v -- FAILED -- [%s] [%s] unable to find the added word %s in the graph -- action: %s",
				idx,
				module,
				funcname,
				test.query,
				test.name,
			)
			return
		}
	}

	for idx, test := range tests {
		verify(idx, test)
	}
}

func TestGetNodes(t *testing.T) {
	module := "Graph"
	funcname := "GetNodes()"

	_ = module
	_ = funcname

	type test struct {
		name  string
		query string
		len   int
	}

	root := New()
	root.Add("cat", "dog", "cot", "cog")

	var tests = []test{
		{
			name:  "valid query",
			query: "cat",
			len:   3,
		},
		{
			name:  "zero-length query",
			query: "",
			len:   0,
		},
		{
			name:  "query for a non-existent item",
			query: "log",
			len:   0,
		},
	}

	var verify = func(idx int, test test) {

		if nodes := root.GetNodes(test.query); len(nodes) != test.len {
			t.Errorf(
				"#%v -- FAILED -- [%s] [%s] fetched nodes doesn't match expected length: wanted %v ; got %v -- action: %s",
				idx,
				module,
				funcname,
				test.len,
				len(nodes),
				test.name,
			)
			return
		}
	}

	for idx, test := range tests {
		verify(idx, test)
	}
}

func TestSiblings(t *testing.T) {
	module := "Graph"
	funcname := "Siblings()"

	_ = module
	_ = funcname

	type test struct {
		name  string
		query string
		wants []string
		err   error
	}

	root := New()
	root.Add("cat", "dog", "pat", "cot", "fur")

	var tests = []test{
		{
			name:  "valid query",
			query: "cat",
			wants: []string{"pat", "cot"},
		},
		{
			name:  "invalid query -- word does not exist",
			query: "fog",
			err:   ErrNonExistent,
		},
	}

	var verify = func(idx int, test test) {
		siblings, err := root.Siblings(test.query)

		if err != nil {
			if !errors.Is(test.err, err) {
				t.Errorf(
					"#%v -- FAILED -- [%s] [%s] unexpected error occurred: %v -- action: %s",
					idx,
					module,
					funcname,
					err,
					test.name,
				)
				return
			} else {
				// intended error, can exit safely
				return
			}

		}

		if !reflect.DeepEqual(siblings, test.wants) {
			t.Errorf(
				"#%v -- FAILED -- [%s] [%s] output mismatch error: wanted %v ; got %v -- action: %s",
				idx,
				module,
				funcname,
				test.wants,
				siblings,
				test.name,
			)
			return
		}
	}

	for idx, test := range tests {
		verify(idx, test)
	}
}

func TestTargetSiblings(t *testing.T) {
	module := "Graph"
	funcname := "TargetSiblings()"

	_ = module
	_ = funcname

	type test struct {
		name   string
		query  string
		target string
		wants  []*Result
		err    error
	}

	root := New()
	root.Add("cat", "dog", "pat", "cot", "cog", "fur", "fir", "zap", "catt")

	var tests = []test{
		{
			name:   "valid query",
			query:  "cat",
			target: "dog",
			wants: []*Result{
				{
					word:      "cot",
					siblings:  []string{"cat", "cog"},
					weight:    33,
					potential: 2,
				},
				{
					word:      "catt",
					siblings:  []string{"cat"},
					weight:    0,
					potential: 1,
				},
				{
					word:      "pat",
					siblings:  []string{"cat"},
					weight:    0,
					potential: 1,
				},
			},
		},
		{
			name:   "invalid query -- word does not exist",
			query:  "fog",
			target: "far",
			err:    ErrNonExistent,
		},
		{
			name:   "invalid query -- zero match error",
			query:  "zap",
			target: "far",
			err:    ErrNoRoute,
		},
		{
			name:   "valid query",
			query:  "fur",
			target: "for",
			wants: []*Result{
				{
					word:      "fir",
					siblings:  []string{"fur"},
					weight:    66,
					potential: 1,
				},
			},
		},
	}

	var verify = func(idx int, test test) {
		results, err := root.TargetSiblings(test.query, test.target)

		if err != nil {
			if !errors.Is(test.err, err) {
				t.Errorf(
					"#%v -- FAILED -- [%s] [%s] unexpected error occurred: %v -- action: %s",
					idx,
					module,
					funcname,
					err,
					test.name,
				)
				return
			} else {
				// intended error, can exit safely
				return
			}

		}

		if !reflect.DeepEqual(results, test.wants) {
			t.Errorf(
				"#%v -- FAILED -- [%s] [%s] output mismatch error: wanted %v ; got %v -- action: %s",
				idx,
				module,
				funcname,
				test.wants,
				results,
				test.name,
			)
			return
		}

	}

	for idx, test := range tests {
		verify(idx, test)
	}
}
