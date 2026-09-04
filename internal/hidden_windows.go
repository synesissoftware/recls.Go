// Copyright 2019-2026, Matthew Wilson and Synesis Information Systems. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

/*
 * Created: 3rd September 2026
 * Updated: 4th September 2026
 */

package internal

import (
	"os"
	"syscall"
)

// Reports whether a Windows entry name looks hidden by convention
// (leading '.'). Full attribute checks use file attributes from
// os.FileInfo when available.
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

// Reports whether info carries FILE_ATTRIBUTE_HIDDEN. Prefers attributes
// already present on FileInfo.Sys(); falls back to GetFileAttributes only
// when Sys() is unavailable.
func IsHiddenFile(
	path string,
	info os.FileInfo,
) bool {

	if info != nil {
		switch sys := info.Sys().(type) {
		case *syscall.Win32FileAttributeData:
			return sys.FileAttributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0
		case *syscall.Win32finddata:
			return sys.FileAttributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0
		}
	}

	ptr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return false
	}
	attrs, err := syscall.GetFileAttributes(ptr)
	if err != nil {
		return false
	}
	return attrs&syscall.FILE_ATTRIBUTE_HIDDEN != 0
}

// Combines name and attribute checks on Windows.
func ProbeHidden(
	path string,
	name string,
	info os.FileInfo,
) bool {

	if IsHiddenFile(path, info) {
		return true
	}

	if IsHiddenName(name) {
		return true
	}

	return false
}
