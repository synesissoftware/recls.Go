// Copyright 2019-2026, Matthew Wilson and Synesis Information Systems. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 * Created: 3rd September 2026
 * Updated: 3rd September 2026
 */

package internal

import (
	"github.com/synesissoftware/shwild.Go"

	"strings"
)

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
