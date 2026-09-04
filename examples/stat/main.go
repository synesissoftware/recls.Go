package main

import (
	recls "github.com/synesissoftware/recls.Go"

	"fmt"
	"os"
)

func main() {
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	e, err := recls.Stat(path, recls.DirectoryParts|recls.MarkDirs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "stat: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("path:      %s\n", e.Path())
	fmt.Printf("location:  %s\n", e.Location)
	fmt.Printf("entry:     %s\n", e.EntryName)
	fmt.Printf("stem:      %s\n", e.Stem)
	fmt.Printf("extension: %s\n", e.Extension)
	fmt.Printf("exists:    %v\n", e.Exists())
	fmt.Printf("is_dir:    %v\n", e.IsDir())
	fmt.Printf("is_file:   %v\n", e.IsFile())
	fmt.Printf("size:      %d\n", e.Size())
}
