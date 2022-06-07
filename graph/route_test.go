package graph

import (
	"errors"
	"reflect"
	"testing"
)

func TestFindRoute(t *testing.T) {
	module := "Graph"
	funcname := "FindRoute()"

	_ = module
	_ = funcname

	type test struct {
		name   string
		origin string
		target string
		wants  []string
		err    error
	}

	root := New()
	root.Add("ruby", "rudy", "tubb", "tuby", "rubb", "rudd", "mudd", "muda", "bare", "rare", "rabi")

	var tests = []test{
		{
			name:   "valid 5-step route",
			origin: "ruby",
			target: "muda",
			wants: []string{
				"ruby", "rudy", "rudd", "mudd", "muda",
			},
		},
		{
			name:   "valid 4-step route",
			origin: "ruby",
			target: "mudd",
			wants: []string{
				"ruby", "rudy", "rudd", "mudd",
			},
		},
		{
			name:   "valid 3-step route",
			origin: "ruby",
			target: "rudd",
			wants: []string{
				"ruby", "rudy", "rudd",
			},
		},
		{
			name:   "valid 2-step route",
			origin: "ruby",
			target: "rudy",
			wants: []string{
				"ruby", "rudy",
			},
		},
		{
			name:   "invalid call -- origin is same as target",
			origin: "ruby",
			target: "ruby",
			err:    ErrSameWord,
		},
		{
			name:   "invalid call -- target not in graph",
			origin: "ruby",
			target: "raly",
			err:    ErrNonExistent,
		},
		{
			name:   "invalid call -- origin not in graph",
			origin: "raly",
			target: "ruby",
			err:    ErrNonExistent,
		},
		{
			name:   "zero routes found",
			origin: "rabi",
			target: "rudy",
			err:    ErrNoRoute,
		},
	}

	var verify = func(idx int, test test) {
		route, err := root.FindRoute(test.origin, test.target)

		if err != nil {
			if !errors.Is(err, test.err) {
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

		if !reflect.DeepEqual(route, test.wants) {
			t.Errorf(
				"#%v -- FAILED -- [%s] [%s] output mismatch error: wanted %v ; got %v -- action: %s",
				idx,
				module,
				funcname,
				test.wants,
				route,
				test.name,
			)
			return
		}
	}

	for idx, test := range tests {
		verify(idx, test)
	}
}
