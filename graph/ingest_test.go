package graph

import (
	"os"
	"testing"
)

func TestFromWordList(t *testing.T) {
	module := "Graph"
	funcname := "FromWordList()"

	_ = module
	_ = funcname

	type test struct {
		name string
		path string
		ok   bool
	}

	var tests = []test{
		{
			name: "from path -- os.Getenv(\"WORD_LIST\")",
			path: os.Getenv("WORD_LIST"),
			ok:   true,
		},
		{
			name: "from empty path",
			path: "",
			ok:   true,
		},
	}

	var verify = func(idx int, test test) {
		list, err := FromWordList(test.path)

		if err != nil && test.ok {
			t.Errorf(
				"#%v -- FAILED -- [%s] [%s] unexpected error: %v -- action: %s",
				idx,
				module,
				funcname,
				err,
				test.name,
			)
			return
		}

		if len(list) == 0 {
			t.Errorf(
				"#%v -- FAILED -- [%s] [%s] retrieved list is zero-length -- action: %s",
				idx,
				module,
				funcname,
				test.name,
			)
			return
		}

	}

	for idx, test := range tests {
		verify(idx, test)
	}

}
