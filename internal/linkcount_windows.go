// Copyright 2019-2026, Matthew Wilson and Synesis Information Systems. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

/*
 * Created: 4th September 2026
 * Updated: 4th September 2026
 */

package internal

import (
	"os"
	"syscall"
)

// Returns the hard-link count via GetFileInformationByHandle.
// Win32FileAttributeData / Win32finddata do not carry nNumberOfLinks.
//
// Parameters:
//   - path — absolute or openable path of the entry;
//   - info — unused on Windows (count requires a handle);
//
// Returns:
//   - the link count, and true when available;
func LinkCount(
	path string,
	info os.FileInfo,
) (uint64, bool) {

	_ = info

	if path == "" {
		return 0, false
	}

	p, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return 0, false
	}

	h, err := syscall.CreateFile(
		p,
		0,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_DELETE,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_FLAG_BACKUP_SEMANTICS,
		0,
	)
	if err != nil {
		return 0, false
	}
	defer syscall.CloseHandle(h)

	var fi syscall.ByHandleFileInformation
	if err := syscall.GetFileInformationByHandle(h, &fi); err != nil {
		return 0, false
	} else {
		return uint64(fi.NumberOfLinks), true
	}
}
