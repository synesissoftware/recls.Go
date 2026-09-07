// Copyright 2019-2026, Matthew Wilson and Synesis Information Systems. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build unix

/*
 * Created: 4th September 2026
 * Updated: 4th September 2026
 */

package internal

import (
	"os"
	"syscall"
)

// Returns the hard-link count from info when Sys() is *syscall.Stat_t.
// On Unix, directories normally report at least 2 (`.` and `..`).
//
// Parameters:
//   - path — unused on Unix (count comes from info);
//   - info — file metadata from Stat / Lstat / ReadDir;
//
// Returns:
//   - the link count, and true when available;
func LinkCount(
	path string,
	info os.FileInfo,
) (uint64, bool) {

	_ = path

	if info == nil {
		return 0, false
	}

	st, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, false
	} else {
		return uint64(st.Nlink), true
	}
}
