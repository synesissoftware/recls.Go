// Copyright 2019-2026, Matthew Wilson and Synesis Information Systems. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 * Created: 3rd September 2026
 * Updated: 4th September 2026
 */

package recls

import (
	"github.com/synesissoftware/recls.Go/internal"

	"errors"
	"io/fs"
	"iter"
	"os"
	"path/filepath"
)

// Accepted by Search and SearchFunc: either a multi-pattern string (split
// on '|' and the platform path-list separator) or a slice of discrete
// pattern strings (not re-split). Empty string or empty/blank-only slice
// matches all names ("*"). Matching is against the entry basename via
// shwild, not filepath.Match.
type PatternSource interface {
	string | []string
}

// Returns a depth-first sequence of matching entries under root.
//
// Parameters:
//   - root — the root directory to search;
//   - patterns — a PatternSource (string or []string; see PatternSource);
//   - opts — options that moderate the search;
func Search[P PatternSource](
	root string,
	patterns P,
	opts SearchOptions,
) iter.Seq2[Entry, error] {

	return func(yield func(Entry, error) bool) {
		err := SearchFunc(root, patterns, opts, func(e Entry) error {
			if !yield(e, nil) {
				return errStopIteration
			}
			return nil
		})

		if err != nil && !errors.Is(err, errStopIteration) {
			yield(Entry{}, err)
		}
	}
}

var errStopIteration = errors.New("recls: stop iteration")

// Invokes fn for each matching entry under root, depth-first. A non-nil
// error from fn aborts the search and is returned.
//
// Parameters:
//   - root — the root directory to search;
//   - patterns — a PatternSource (string or []string; see PatternSource);
//   - opts — options that moderate the search;
//   - fn — callback invoked for each matching entry;
func SearchFunc[P PatternSource](
	root string,
	patterns P,
	opts SearchOptions,
	fn func(Entry) error,
) error {

	flags := NormaliseFlags(opts.Flags)

	searchRoot, err := internal.ResolveSearchRoot(
		root,
		0 != (flags&UseTildeOnNoSearchRoot),
		PathNameSeparator,
	)
	if err != nil {
		return err
	}

	rawPatterns := internal.NormalisePatterns(patterns, PathSeparator)
	if err := internal.ValidatePatternsForFlags(rawPatterns, 0 != (flags&Recursive)); err != nil {
		return ErrDotRecursiveSearch
	}
	compiled, err := internal.CompilePatterns(rawPatterns)
	if err != nil {
		return err
	}

	var guard *internal.LoopGuard
	if 0 == (flags & NoBreakInfiniteLoops) {
		guard = internal.NewLoopGuard()
		if !guard.Enter(searchRoot) {
			return nil
		}
	}

	return searchDirectory(searchRoot, searchRoot, compiled, flags, opts.OnDirectory, guard, fn)
}

func searchDirectory(
	searchRoot string,
	dir string,
	patterns internal.CompiledPatterns,
	flags SearchFlags,
	onDirectory func(string) error,
	guard *internal.LoopGuard,
	fn func(Entry) error,
) error {

	dir = internal.EnsureTrailingSep(dir, PathNameSeparator)

	if 0 != (flags&DirProgress) && onDirectory != nil {
		if err := onDirectory(dir); err != nil {
			return err
		}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		if 0 != (flags & StopOnAccessFailure) {
			return &AccessDeniedError{Path: dir, Err: err}
		}
		return nil
	}

	type pendingDir struct {
		path string
		info os.FileInfo
	}
	var subdirs []pendingDir

	for _, de := range entries {
		name := de.Name()
		if internal.IsDots(name) {
			continue
		}

		entryPath := filepath.Join(internal.StripTrailingSeps(dir), name)

		info, err := dirEntryInfo(de, entryPath, flags)
		if err != nil {
			if 0 != (flags & StopOnAccessFailure) {
				return &AccessDeniedError{Path: entryPath, Err: err}
			}
			continue
		}

		hidden := internal.ProbeHidden(entryPath, name, info)
		if hidden && 0 != (flags&IgnoreHiddenEntries) {
			continue
		}

		match, err := patterns.MatchesAny(name)
		if err != nil {
			return err
		}

		if match && typeMatches(info, flags) {
			e, err := buildEntry(entryPath, searchRoot, searchRoot, info, flags)
			if err != nil {
				if 0 != (flags & StopOnAccessFailure) {
					return err
				}
				continue
			}
			if err := fn(e); err != nil {
				return err
			}
		}

		if info.IsDir() && 0 != (flags&Recursive) {
			// Do not descend into directory symlinks when NoFollowLinks.
			if isSymlink(info) && 0 != (flags&NoFollowLinks) {
				continue
			}
			if hidden && 0 != (flags&IgnoreHiddenEntries) {
				continue
			}
			subdirs = append(subdirs, pendingDir{path: entryPath, info: info})
		}
	}

	for _, sd := range subdirs {
		if guard != nil && !guard.Enter(sd.path) {
			continue
		}
		if err := searchDirectory(searchRoot, sd.path, patterns, flags, onDirectory, guard, fn); err != nil {
			return err
		}
	}

	return nil
}

func dirEntryInfo(
	de os.DirEntry,
	path string,
	flags SearchFlags,
) (os.FileInfo, error) {

	if 0 != (flags & NoFollowLinks) {
		return os.Lstat(path)
	}
	// Prefer Info() but fall back to Stat for follow-symlink semantics.
	info, err := de.Info()
	if err != nil {
		return os.Stat(path)
	}
	if isSymlink(info) {
		return os.Stat(path)
	}
	return info, nil
}

func isSymlink(info os.FileInfo) bool {
	return info.Mode()&fs.ModeSymlink != 0
}

func typeMatches(info os.FileInfo, flags SearchFlags) bool {
	types := flags & TypeMask
	if types == 0 {
		types = Files
	}

	isLink := isSymlink(info)
	isDir := info.IsDir()
	isFile := info.Mode().IsRegular()

	// When following links, Stat may report the target type; Lstat keeps
	// the link itself. Count links toward Links when NoFollowLinks is set.
	if isLink && 0 != (flags&Links) && 0 != (flags&NoFollowLinks) {
		return true
	}

	matched := false
	if 0 != (types&Files) && isFile {
		matched = true
	}
	if 0 != (types&Directories) && isDir {
		matched = true
	}
	if 0 != (types&Links) && isLink {
		matched = true
	}
	return matched
}
