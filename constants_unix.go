// Copyright 2019-2026, Matthew Wilson and Synesis Information Systems. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build unix

/*
 * Created: 3rd September 2026
 * Updated: 3rd September 2026
 */

package recls

import (
	libpath_unix "github.com/synesissoftware/libpath.Go/util/unix"
)

const (
	// Platform path-name (directory) separator.
	PathNameSeparator = string(libpath_unix.PathElementSeparator)
	// Platform path-list separator (e.g. for multi-pattern search strings).
	PathSeparator = string(libpath_unix.PathSeparator)
)
