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

	opts := recls.SearchOptions{
		Flags: recls.Files | recls.Recursive | recls.MarkDirs,
		// Flags: recls.Files | recls.Directories | recls.Recursive | recls.MarkDirs,
	}

	err := recls.SearchFunc(root, "*", opts, func(e recls.Entry) error {
		n := e.LinkCount()
		if n <= 1 {
			return nil
		}

		kind := "file"
		if e.IsDir() {
			kind = "dir"
		} else if e.IsLink() {
			kind = "link"
		}

		fmt.Printf("%6d  %-4s  %s\n", n, kind, e.SearchRelativePath)
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "hard_links: %v\n", err)
		os.Exit(1)
	}
}
