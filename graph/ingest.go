package graph

import (
	"io/ioutil"
	"net/http"
	"os"
)

const delimiter byte = 10 // bytes for newline (\n)
const dwylEnglishWordsRepo = `https://raw.githubusercontent.com/dwyl/english-words/master/words_alpha.txt`

// FromWordList function will get a list of strings from the file retrieved from the path
// provided. It separates words by newlines.
func FromWordList(path string) ([]string, error) {
	if path == "" {
		return FromOnlineSource(dwylEnglishWordsRepo)
	}

	b, err := os.ReadFile(path)

	if err != nil {
		return FromOnlineSource(dwylEnglishWordsRepo)
	}

	return fromBytes(b), nil
}

// FromOnlineSource function will get a list of strings (separated by newlines) from an
// endpoint on the internet. If no page URL is provided, the default repo (dwyl/english-words)'s
// word list is fetched.
func FromOnlineSource(pageURL string) ([]string, error) {
	if pageURL == "" {
		return FromOnlineSource(dwylEnglishWordsRepo)
	}

	response, err := http.Get(pageURL)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	b, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	return fromBytes(b), nil
}

// fromBytes function is a converter to materialize raw data (a slice of bytes)
// into a list of strings (separated by newlines, or byte 10)
func fromBytes(b []byte) []string {
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

	return out
}
