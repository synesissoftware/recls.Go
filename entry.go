// Copyright 2019-2026, Matthew Wilson and Synesis Information Systems. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 * Created: 3rd September 2026
 * Updated: 3rd September 2026
 */

package recls

import (
	libpath_common "github.com/synesissoftware/libpath.Go/parse/common"

	"io/fs"
	"os"
	"time"
)

// A file-system entry discovered by Stat or Search. Embeds libpath's
// PathDescriptor for path decomposition; search-relative fields and
// file-system nature are layered on top. Implements fs.FileInfo.
type Entry struct {
	libpath_common.PathDescriptor

	// Absolute search root when produced by Search; empty for Stat.
	SearchRoot string
	// Path relative to SearchRoot.
	SearchRelativePath string
	// Location relative to SearchRoot.
	SearchRelativeDirectory string
	// Split SearchRelativeDirectory (parts end with PathNameSeparator).
	SearchRelativeDirectoryParts []string

	fileInfo os.FileInfo
	hidden   bool
	readonly bool
}

// Fullest establishable path (recls name for PathDescriptor.FullPath).
func (e Entry) Path() string {
	return e.FullPath
}

// True when file-system metadata was obtained.
func (e Entry) Exists() bool {
	return e.fileInfo != nil
}

// True when the entry is a regular file.
func (e Entry) IsFile() bool {
	return e.fileInfo != nil && e.fileInfo.Mode().IsRegular()
}

// True when the entry is a symbolic link.
func (e Entry) IsLink() bool {
	return e.fileInfo != nil && e.fileInfo.Mode()&fs.ModeSymlink != 0
}

// True when the entry is considered hidden.
func (e Entry) IsHidden() bool {
	return e.hidden
}

// True when the entry is not writable by the owner (Unix: mode lacks the
// owner-write bit; Windows modelling is incomplete in v0.1).
func (e Entry) IsReadonly() bool {
	return e.readonly
}

// fs.FileInfo
func (e Entry) Name() string {
	return e.EntryName
}

func (e Entry) Size() int64 {
	if e.fileInfo == nil {
		return 0
	} else {
		return e.fileInfo.Size()
	}
}

func (e Entry) Mode() fs.FileMode {
	if e.fileInfo == nil {
		return 0
	} else {
		return e.fileInfo.Mode()
	}
}

func (e Entry) ModTime() time.Time {
	if e.fileInfo == nil {
		return time.Time{}
	} else {
		return e.fileInfo.ModTime()
	}
}

func (e Entry) IsDir() bool {
	return e.fileInfo != nil && e.fileInfo.IsDir()
}

func (e Entry) Sys() any {
	if e.fileInfo == nil {
		return nil
	} else {
		return e.fileInfo.Sys()
	}
}

// Underlying os.FileInfo when Exists(); otherwise nil.
func (e Entry) FileInfo() os.FileInfo {
	return e.fileInfo
}
