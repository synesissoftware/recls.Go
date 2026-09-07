// Copyright 2019-2026, Matthew Wilson and Synesis Information Systems. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 * Created: 3rd September 2026
 * Updated: 3rd September 2026
 */

package internal

import (
	"os"
	"path/filepath"
	"strings"
)

// Appends a trailing path-name separator if missing.
func EnsureTrailingSep(
	path string,
	sep string,
) string {
	if path == "" {
		return path
	}
	if strings.HasSuffix(path, sep) {
		return path
	}
	// Also accept the alternate separator already present.
	if sep == "/" && strings.HasSuffix(path, "\\") {
		return path
	}
	if sep == "\\" && strings.HasSuffix(path, "/") {
		return path
	}
	return path + sep
}

// Removes trailing path-name separators, preserving a root path such as "/"
// or "C:\".
func StripTrailingSeps(path string) string {
	if path == "" {
		return path
	}
	for len(path) > 1 {
		last := path[len(path)-1]
		if last != '/' && last != '\\' {
			break
		}
		// Keep Windows drive root "C:\"
		if len(path) == 3 && path[1] == ':' {
			break
		}
		path = path[:len(path)-1]
	}
	return path
}

// Expands a leading "~" or "~/" using the home directory.
func ExpandTilde(path string) (string, error) {
	if path == "" || path[0] != '~' {
		return path, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return path, err
	}
	if path == "~" {
		return home, nil
	}
	if path[1] == '/' || path[1] == '\\' {
		return filepath.Join(home, path[2:]), nil
	}
	return path, nil
}

// Resolves an empty/relative/tilde search root to an absolute directory
// path that ends with a trailing separator.
func ResolveSearchRoot(
	root string,
	useTildeOnEmpty bool,
	sep string,
) (string, error) {
	if root == "" {
		if useTildeOnEmpty {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			return EnsureTrailingSep(home, sep), nil
		}
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return EnsureTrailingSep(wd, sep), nil
	}

	expanded, err := ExpandTilde(root)
	if err != nil {
		return "", err
	}

	abs, err := filepath.Abs(expanded)
	if err != nil {
		return "", err
	}
	return EnsureTrailingSep(filepath.Clean(abs), sep), nil
}

// Expands tilde and makes path absolute relative to the given reference
// directory (usually the search root or cwd).
func ResolveEntryPath(
	path string,
	referenceDir string,
) (string, error) {
	expanded, err := ExpandTilde(path)
	if err != nil {
		return "", err
	}
	if filepath.IsAbs(expanded) {
		return filepath.Clean(expanded), nil
	}
	ref := referenceDir
	if ref == "" {
		ref, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}
	ref = StripTrailingSeps(ref)
	return filepath.Clean(filepath.Join(ref, expanded)), nil
}

// Reports whether name is "." or "..".
func IsDots(name string) bool {
	return name == "." || name == ".."
}

// Splits a directory string into parts that each end with the given
// separator (recls convention), preserving a leading root separator as its
// own part when present.
func SplitDirectoryParts(
	directory string,
	sep string,
) []string {
	if directory == "" {
		return nil
	}
	// Prefer libpath-produced parts when available; this helper is a
	// fallback for search-relative directories.
	raw := strings.SplitAfter(directory, sep)
	parts := make([]string, 0, len(raw))
	for _, p := range raw {
		if p == "" {
			continue
		}
		parts = append(parts, p)
	}
	return parts
}

// Returns path relative to origin, following recls.Ruby
// Ximpl::Util.derive_relative_path semantics (not filepath.Rel alone).
func DeriveRelativePath(
	origin string,
	path string,
	sep string,
) string {
	if path == "" {
		return ""
	}
	if origin == "" {
		return path
	}

	trailing := ""
	if strings.HasSuffix(path, "/") || strings.HasSuffix(path, "\\") {
		trailing = sep
	}

	origin = filepath.Clean(origin)
	path = filepath.Clean(path)

	// Bare "." origin → path as-is with optional trailing semantics handled
	// by caller; Ruby returns path when origin matches /^\.[\\\/]*$/.
	if origin == "." {
		if trailing != "" && !strings.HasSuffix(path, sep) {
			return path + trailing
		}
		return path
	}

	origin = StripTrailingSeps(origin)
	path = StripTrailingSeps(path)

	// Windows: different drives → absolute path.
	if len(path) >= 2 && len(origin) >= 2 && path[1] == ':' && origin[1] == ':' {
		if strings.EqualFold(path[:1], origin[:1]) == false {
			return path + trailing
		}
	}

	pathParts := splitPathElements(path)
	originParts := splitPathElements(origin)

	for len(pathParts) > 0 && len(originParts) > 0 {
		if !pathElementsEqual(pathParts[0], originParts[0]) {
			break
		}
		pathParts = pathParts[1:]
		originParts = originParts[1:]
	}

	if len(pathParts) == 0 && len(originParts) == 0 {
		return "." + trailing
	}

	var b strings.Builder
	for range originParts {
		b.WriteString("..")
		b.WriteString(sep)
	}
	for i, p := range pathParts {
		b.WriteString(p)
		if i+1 < len(pathParts) {
			b.WriteString(sep)
		}
	}
	result := b.String()
	if trailing != "" && !strings.HasSuffix(result, sep) {
		result += trailing
	}
	return result
}

func splitPathElements(path string) []string {
	path = strings.ReplaceAll(path, "\\", "/")
	parts := strings.Split(path, "/")
	out := make([]string, 0, len(parts))
	for i, p := range parts {
		if p == "" {
			// Keep a single empty leading element to represent root "/".
			if i == 0 {
				out = append(out, "")
			}
			continue
		}
		out = append(out, p)
	}
	return out
}

func pathElementsEqual(
	a string,
	b string,
) bool {
	// Exact match for now. Windows should use case-insensitive comparison
	// (e.g. strings.EqualFold); prefer upstreaming compare-path / case-fold
	// helpers into libpath.Go (Ruby already has make_compare_path) rather
	// than inventing a third variant here — see TODO.md.
	return a == b
}
