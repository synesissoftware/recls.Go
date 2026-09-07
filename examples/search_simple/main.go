package main

import (
	recls "github.com/synesissoftware/recls.Go"

	"fmt"
	"os"
)

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	pattern := "*"
	if len(os.Args) > 2 {
		pattern = os.Args[2]
	}

	opts := recls.SearchOptions{
		Flags: recls.Files | recls.Recursive | recls.DirectoryParts,
	}

	err := recls.SearchFunc(root, pattern, opts, func(e recls.Entry) error {
		fmt.Println(e.SearchRelativePath)
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "search: %v\n", err)
		os.Exit(1)
	}
}
