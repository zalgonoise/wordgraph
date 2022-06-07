package graph

import "os"

const delimiter byte = 10 // bytes for newline (\n)

// FromWordList function will get a list of strings from the file retrieved from the path
// provided. It separates words by newlines.
func FromWordList(path string) ([]string, error) {
	b, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var out []string
	var word []byte

	for _, u := range b {
		if u != delimiter {
			word = append(word, u)
		} else {
			out = append(out, string(word))
			word = []byte{}
		}
	}

	if len(word) > 0 {
		out = append(out, string(word))
	}

	return out, nil
}
