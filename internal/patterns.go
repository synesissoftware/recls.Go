// Copyright 2019-2026, Matthew Wilson and Synesis Information Systems. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 * Created: 3rd September 2026
 * Updated: 4th September 2026
 */

package internal

import (
	"github.com/synesissoftware/shwild.Go"

	"strings"
)

// Accepted by NormalisePatterns: a multi-pattern string, or a slice of
// discrete pattern strings.
type PatternSource interface {
	string | []string
}

// Normalises patterns into a non-empty slice of discrete pattern strings.
// A string is split on '|' and the path-list separator (see SplitPatterns).
// A slice is treated as already-discrete patterns: each element is trimmed
// and empty elements are dropped; elements are not re-split on '|' or the
// path-list separator. An empty string or empty/blank-only slice yields
// ["*"].
//
// Parameters:
//   - patterns — a string or []string pattern source;
//   - separator — the platform path-list separator (used only for strings);
//
// Returns:
//   - a slice of pattern strings;
func NormalisePatterns[P PatternSource](
	patterns P,
	separator string,
) []string {

	switch v := any(patterns).(type) {
	case string:
		return SplitPatterns(v, separator)
	case []string:
		return normalisePatternSlice(v)
	default:
		panic("recls: unexpected pattern source type")
	}
}

// Splits a multi-pattern string on '|' and the platform path-list
// separator. An empty/blank patterns string yields ["*"]. Does not invent
// additional separators beyond those two.
//
// Parameters:
//   - patterns — a multi-pattern string;
//   - separator — the platform path-list separator;
//
// Returns:
//   - a slice of pattern strings;
func SplitPatterns(
	patterns string,
	separator string,
) []string {

	if patterns == "" {
		return []string{"*"}
	}

	s := patterns
	if separator != "" && separator != "|" {
		s = strings.ReplaceAll(s, separator, "|")
	}

	raw := strings.Split(s, "|")
	out := make([]string, 0, len(raw))
	seen := make(map[string]struct{}, len(raw))
	for _, p := range raw {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		out = append(out, p)
	}
	if len(out) == 0 {
		return []string{"*"}
	}
	return out
}

func normalisePatternSlice(patterns []string) []string {
	if len(patterns) == 0 {
		return []string{"*"}
	}

	out := make([]string, 0, len(patterns))
	seen := make(map[string]struct{}, len(patterns))
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		out = append(out, p)
	}
	if len(out) == 0 {
		return []string{"*"}
	}
	return out
}

// Holds precompiled shwild patterns for a search.
type CompiledPatterns struct {
	Patterns []shwild.CompiledPattern
	Raw      []string
}

// Compiles each pattern string with shwild.
func CompilePatterns(patterns []string) (CompiledPatterns, error) {
	cps := make([]shwild.CompiledPattern, 0, len(patterns))
	for _, p := range patterns {
		cp, err := shwild.Compile(p)
		if err != nil {
			return CompiledPatterns{}, err
		}
		cps = append(cps, cp)
	}
	return CompiledPatterns{Patterns: cps, Raw: patterns}, nil
}

// Rejects "." / ".." patterns under Recursive.
func ValidatePatternsForFlags(
	patterns []string,
	recursive bool,
) error {

	if !recursive {
		return nil
	}
	for _, p := range patterns {
		if p == "." || p == ".." {
			return errDotRecursive
		}
	}
	return nil
}

var errDotRecursive = errDotRecursiveType{}

type errDotRecursiveType struct{}

func (errDotRecursiveType) Error() string {
	return "recls: dot/dot-dot pattern not allowed in recursive search"
}

// Reports whether name matches any compiled pattern. Matching is against
// the entry name (basename), not the full path.
func (cp CompiledPatterns) MatchesAny(name string) (bool, error) {
	for _, p := range cp.Patterns {
		ok, err := p.Match(name)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}
