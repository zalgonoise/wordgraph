package graph

import (
	"testing"
)

func TestResultPrint(t *testing.T) {
	module := "Result"
	funcname := "Print()"

	_ = module
	_ = funcname

	type test struct {
		name   string
		query  string
		target string
		print  map[string]string
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
			print: map[string]string{
				"cot":  "{[cot] weight: 33 potential: 2}",
				"catt": "{[catt] weight: 0 potential: 1}",
				"pat":  "{[pat] weight: 0 potential: 1}",
			},
		},
	}

	var verify = func(idx int, test test) {
		results, _ := root.TargetSiblings(test.query, test.target)

		for _, result := range results {
			for k, v := range test.print {
				if result.word == k {
					output := result.String()

					if output != v {
						t.Errorf(
							"#%v -- FAILED -- [%s] [%s] output mismatch error: wanted %v ; got %v -- action: %s",
							idx,
							module,
							funcname,
							v,
							output,
							test.name,
						)
						return
					}
				}
			}

		}
	}

	for idx, test := range tests {
		verify(idx, test)
	}
}
