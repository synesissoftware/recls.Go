// Copyright 2019-2026, Matthew Wilson and Synesis Information Systems. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build unix

/*
 * Created: 3rd September 2026
 * Updated: 4th September 2026
 */

package internal

import (
	"os"
)

// Reports whether a Unix entry name is hidden (leading '.'), excluding
// "." and "..".
func IsHiddenName(name string) bool {
	switch name {
	case "":
		return false
	case ".", "..":
		return false
	default:
		return name[0] == '.'
	}
}

// Reports whether a Unix entry is hidden by name.
func ProbeHidden(
	path string,
	name string,
	info os.FileInfo,
) bool {

	_ = path
	_ = info

	return IsHiddenName(name)
}
