package main

import (
	recls "github.com/synesissoftware/recls.Go"

	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const timeFormat = "02/01/2006 03:04:05 PM"

func main() {
	var (
		displayNameOnly    bool
		displaySearchRoots bool
		displayTotals      bool
		displayWithCommas  bool
		markDirectories    bool
		succinct           bool
		includeDirectories bool
		includeHidden      bool
		includeFiles       bool
		nonRecursive       bool
		trimToCwd          bool
		trimToSearchRoot   bool
		showTotalSize      bool
		minSize            int64
		maxSize            int64
	)

	flag.BoolVar(&displayNameOnly, "f", false, "display entry name only")
	flag.BoolVar(&displayNameOnly, "display-name-only", false, "display entry name only")
	flag.BoolVar(&displaySearchRoots, "d", false, "display search roots (and skip listing)")
	flag.BoolVar(&displaySearchRoots, "display-search-roots", false, "display search roots (and skip listing)")
	flag.BoolVar(&displayTotals, "c", false, "display totals")
	flag.BoolVar(&displayTotals, "n", false, "display totals")
	flag.BoolVar(&displayTotals, "display-totals", false, "display totals")
	flag.BoolVar(&displayWithCommas, ",", false, "insert commas into sizes and totals")
	flag.BoolVar(&displayWithCommas, "display-with-commas", false, "insert commas into sizes and totals")
	flag.BoolVar(&markDirectories, "m", false, "mark directories with a trailing separator")
	flag.BoolVar(&markDirectories, "mark-directories", false, "mark directories with a trailing separator")
	flag.BoolVar(&succinct, "s", false, "succinct output — path only")
	flag.BoolVar(&succinct, "succinct", false, "succinct output — path only")
	flag.BoolVar(&includeDirectories, "dirs", false, "include directories")
	flag.BoolVar(&includeDirectories, "directories", false, "include directories")
	flag.BoolVar(&includeDirectories, "include-directories", false, "include directories")
	flag.BoolVar(&includeHidden, "hidden", false, "include hidden entries")
	flag.BoolVar(&includeHidden, "include-hidden", false, "include hidden entries")
	flag.BoolVar(&includeFiles, "files", false, "include files")
	flag.BoolVar(&includeFiles, "include-files", false, "include files")
	flag.BoolVar(&nonRecursive, "R", false, "non-recursive search")
	flag.BoolVar(&nonRecursive, "non-recursive-search", false, "non-recursive search")
	flag.BoolVar(&trimToCwd, "t", false, "trim path relative to the current directory")
	flag.BoolVar(&trimToCwd, "trim-to-cwd", false, "trim path relative to the current directory")
	flag.BoolVar(&trimToSearchRoot, "T", false, "trim path relative to the search root")
	flag.BoolVar(&trimToSearchRoot, "trim-to-search-root", false, "trim path relative to the search root")
	flag.BoolVar(&showTotalSize, "Z", false, "display total size of listed entries")
	flag.BoolVar(&showTotalSize, "show-total-size", false, "display total size of listed entries")
	flag.Int64Var(&minSize, "min-size", -1, "minimum file size (bytes); ignored for directories")
	flag.Int64Var(&maxSize, "max-size", -1, "maximum file size (bytes); ignored for directories")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] [root [pattern ...]]\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\nInspired by Synesis rls (Ruby); recursive listing via recls.Go.\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if trimToCwd && trimToSearchRoot {
		fmt.Fprintln(os.Stderr, "rls: cannot specify both -t/--trim-to-cwd and -T/--trim-to-search-root")
		os.Exit(2)
	}
	if minSize >= 0 && maxSize >= 0 && minSize > maxSize {
		fmt.Fprintln(os.Stderr, "rls: max-size cannot be smaller than min-size")
		os.Exit(2)
	}

	root := "."
	patterns := []string{"*"}
	args := flag.Args()
	if len(args) >= 1 {
		root = args[0]
	}
	if len(args) >= 2 {
		patterns = args[1:]
	}

	var flags recls.SearchFlags
	if includeFiles {
		flags |= recls.Files
	}
	if includeDirectories {
		flags |= recls.Directories
	}
	if 0 == (flags & (recls.Files | recls.Directories | recls.Links)) {
		flags |= recls.Files
	}
	if !nonRecursive {
		flags |= recls.Recursive
	}
	if !includeHidden {
		flags |= recls.IgnoreHiddenEntries
	}
	if markDirectories {
		flags |= recls.MarkDirs
	}

	if displaySearchRoots {
		fmt.Printf("searching %q with %q\n", root, strings.Join(patterns, "|"))
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "rls: %v\n", err)
		os.Exit(1)
	}

	tty := isTTY()
	var numFound int64
	var totalSize int64

	opts := recls.SearchOptions{Flags: flags}
	err = recls.SearchFunc(root, patterns, opts, func(e recls.Entry) error {
		if e.IsFile() {
			if minSize >= 0 && e.Size() < minSize {
				return nil
			}
			if maxSize >= 0 && e.Size() > maxSize {
				return nil
			}
		}

		numFound++

		var path string
		switch {
		case trimToCwd:
			rel, relErr := filepath.Rel(cwd, e.Path())
			if relErr != nil {
				path = e.Path()
			} else {
				path = rel
			}
		case trimToSearchRoot:
			path = e.SearchRelativePath
		case displayNameOnly:
			path = e.EntryName
		default:
			path = e.Path()
		}

		if succinct {
			fmt.Println(path)
		} else {
			date := e.ModTime().Format(timeFormat)
			attr := makeAttr(e)
			sizeStr := ""
			if !e.IsDir() {
				sizeStr = formatSize(e.Size(), displayWithCommas)
				totalSize += e.Size()
			}
			if tty {
				sizeStr = padLeft(sizeStr, 20)
			}
			fmt.Printf("%s\t%s\t%s\t%s\n", date, attr, sizeStr, path)
		}

		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "rls: %v\n", err)
		os.Exit(1)
	}

	if showTotalSize {
		sizeStr := formatSize(totalSize, true)
		if tty {
			sep := strings.Repeat("-", len(sizeStr))
			fmt.Printf("%s%s\n", strings.Repeat(" ", 40), padLeft(sep, 20))
			fmt.Printf("%s%s\n", strings.Repeat(" ", 40), padLeft(sizeStr, 20))
		} else {
			fmt.Printf("\t\t\t%s\n", sizeStr)
		}
	}

	if displayTotals {
		n := strconv.FormatInt(numFound, 10)
		if displayWithCommas {
			n = insertCommas(n)
		}
		suffix := "s"
		if numFound == 1 {
			suffix = ""
		}
		fmt.Printf("\t%s file%s found\n", n, suffix)
	}
}

func makeAttr(e recls.Entry) string {
	r := []byte("--------")
	if e.IsReadonly() {
		r[0] = 'R'
	}
	if e.IsHidden() {
		r[1] = 'H'
	}
	if e.IsDir() {
		r[4] = 'D'
	}
	return string(r)
}

func formatSize(n int64, withCommas bool) string {
	s := strconv.FormatInt(n, 10)
	if withCommas {
		return insertCommas(s)
	} else {
		return s
	}
}

func insertCommas(s string) string {
	neg := false
	if strings.HasPrefix(s, "-") {
		neg = true
		s = s[1:]
	}
	if len(s) <= 3 {
		if neg {
			return "-" + s
		} else {
			return s
		}
	}

	var b strings.Builder
	lead := len(s) % 3
	if lead == 0 {
		lead = 3
	}
	b.WriteString(s[:lead])
	for i := lead; i < len(s); i += 3 {
		b.WriteByte(',')
		b.WriteString(s[i : i+3])
	}
	if neg {
		return "-" + b.String()
	} else {
		return b.String()
	}
}

func padLeft(s string, width int) string {
	if len(s) >= width {
		return s
	} else {
		return strings.Repeat(" ", width-len(s)) + s
	}
}

func isTTY() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	} else {
		return (fi.Mode() & os.ModeCharDevice) != 0
	}
}
