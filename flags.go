// Copyright 2019-2026, Matthew Wilson and Synesis Information Systems. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 * Created: 3rd September 2026
 * Updated: 3rd September 2026
 */

package recls

// Control search and Stat behaviour. Values align with the C recls
// RECLS_F_* flags where practical.
type SearchFlags uint32

const (
	// Includes regular files. Default when no type flag is set.
	Files SearchFlags = 0x0000_0001
	// Includes directories.
	Directories SearchFlags = 0x0000_0002
	// Includes symbolic links (Unix).
	Links SearchFlags = 0x0000_0004
	// Reserved; not implemented in v0.1.
	Devices SearchFlags = 0x0000_0008
	// Reserved; not implemented in v0.1.
	Sockets SearchFlags = 0x0000_0010
	// Selects the entry-type filter bits.
	TypeMask SearchFlags = 0x0000_0FFF

	// Requests per-directory progress callbacks (via
	// SearchOptions.OnDirectory when provided).
	DirProgress SearchFlags = 0x0000_1000
	// Aborts when a directory cannot be read or an entry cannot be stated;
	// default is skip-and-continue.
	StopOnAccessFailure SearchFlags = 0x0000_2000
	// Reserved; not implemented in v0.1.
	LinkCount SearchFlags = 0x0000_4000
	// Reserved; not implemented in v0.1.
	NodeIndex SearchFlags = 0x0000_8000

	// Searches subdirectories depth-first.
	Recursive SearchFlags = 0x0001_0000
	// Uses Lstat for metadata and does not descend into directory
	// symlinks.
	NoFollowLinks SearchFlags = 0x0002_0000
	// Populates DirectoryParts and search-relative parts.
	DirectoryParts SearchFlags = 0x0004_0000
	// Obtains a path-only Entry without requiring existence.
	DetailsLater SearchFlags = 0x0008_0000

	// Appends a trailing path-name separator to directory paths.
	MarkDirs SearchFlags = 0x0020_0000
	// Treats an empty search root as the home directory rather than the
	// current working directory.
	UseTildeOnNoSearchRoot SearchFlags = 0x0400_0000
	// Skips hidden entries (leading '.' on Unix; FILE_ATTRIBUTE_HIDDEN on
	// Windows).
	IgnoreHiddenEntries SearchFlags = 0x0800_0000
	// Disables the Unix/macOS device+inode loop guard; by default the
	// guard is active.
	NoBreakInfiniteLoops SearchFlags = 0x1000_0000
)

// Alias for SearchFlags used with Stat.
type StatOptions = SearchFlags

// Holds flags and optional search hooks.
type SearchOptions struct {
	Flags SearchFlags
	// Invoked for each traversed directory when DirProgress is set. A
	// non-nil return aborts the search.
	OnDirectory func(dir string) error
}

// Applies recls defaults: when no type bit is set, Files is assumed.
func NormaliseFlags(flags SearchFlags) SearchFlags {
	if 0 == (flags & TypeMask) {
		flags |= Files
	}
	return flags
}
