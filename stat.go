// Copyright 2019-2026, Matthew Wilson and Synesis Information Systems. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 * Created: 3rd September 2026
 * Updated: 3rd September 2026
 */

package recls

import (
	"github.com/synesissoftware/recls.Go/internal"

	"os"
)

// Returns an Entry for path. When no type flags are set, Files is assumed
// (via NormaliseFlags) but does not filter Stat results — Stat always
// describes the named path.
//
// Parameters:
//   - path — the path to examine;
//   - opts — options that moderate Stat (e.g. DetailsLater, MarkDirs);
func Stat(
	path string,
	opts StatOptions,
) (Entry, error) {

	flags := NormaliseFlags(opts)

	ref, err := os.Getwd()
	if err != nil {
		return Entry{}, err
	}
	ref = internal.EnsureTrailingSep(ref, PathNameSeparator)

	abs, err := internal.ResolveEntryPath(path, ref)
	if err != nil {
		return Entry{}, err
	}

	info, err := statInfo(abs, flags)
	if err != nil {
		return Entry{}, err
	}
	// DetailsLater: info may be nil when path does not exist.

	// Leave SearchRoot empty — Stat is free of search context (matches
	// recls.Ruby Stat without search_dir).
	e, err := buildEntry(abs, ref, "", info, flags|DirectoryParts)
	if err != nil {
		return Entry{}, err
	}
	return e, nil
}
