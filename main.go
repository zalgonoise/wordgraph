package main

import (
	"fmt"
	"os"

	"github.com/zalgonoise/wordgraph/graph"
)

func main() {
	keyword := "gopher"
	target := "logic"

	w, err := graph.FromWordList(os.Getenv("WORD_LIST"))
	if err != nil {
		fmt.Printf("failed to grab word list: %v -- please check if your ${WORD_LIST} env variable is set", err)
		os.Exit(1)
	}

	n := graph.New()
	n.Add(w...)

	result, err := n.FindRoute(keyword, target)
	if err != nil {
		fmt.Printf("failed to find route with an error: %v", err)
		os.Exit(1)
	}

	fmt.Println(result)

}
