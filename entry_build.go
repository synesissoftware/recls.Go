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

	libpath_parse "github.com/synesissoftware/libpath.Go/parse"
	libpath_common "github.com/synesissoftware/libpath.Go/parse/common"

	"os"
)

// Constructs an Entry from an absolute (or fullest) path, optional
// FileInfo, and optional search root.
func buildEntry(
	entryPath string,
	referenceDir string,
	searchRoot string,
	info os.FileInfo,
	flags SearchFlags,
) (Entry, error) {

	pd, err := libpath_parse.ParsePathString(entryPath, referenceDir)
	if err != nil {
		return Entry{}, &InvalidPathError{Path: entryPath, Err: err}
	}

	e := entryFromDescriptor(pd, flags)

	if searchRoot != "" {
		e.SearchRoot = searchRoot
		relPath := internal.DeriveRelativePath(searchRoot, e.FullPath, PathNameSeparator)
		relDir := internal.DeriveRelativePath(searchRoot, e.Location, PathNameSeparator)
		e.SearchRelativePath = relPath
		e.SearchRelativeDirectory = relDir
		// Always populate search-relative parts for discoverability in
		// v0.1; DirectoryParts still gates absolute DirectoryParts.
		e.SearchRelativeDirectoryParts = internal.SplitDirectoryParts(
			ensureDirPartsForm(relDir),
			PathNameSeparator,
		)
	}

	if info != nil {
		applyFileInfo(&e, info)
		e.hidden = internal.ProbeHidden(e.FullPath, e.EntryName, info)
	}

	if e.IsDir() && 0 != (flags&MarkDirs) {
		e.FullPath = internal.EnsureTrailingSep(e.FullPath, PathNameSeparator)
		if e.SearchRelativePath != "" {
			e.SearchRelativePath = internal.EnsureTrailingSep(e.SearchRelativePath, PathNameSeparator)
		}
	}

	return e, nil
}

func entryFromDescriptor(
	pd libpath_common.PathDescriptor,
	flags SearchFlags,
) Entry {

	e := Entry{
		PathDescriptor: pd,
	}
	if 0 == (flags & DirectoryParts) {
		e.DirectoryParts = nil
	}
	return e
}

func ensureDirPartsForm(dir string) string {
	if dir == "" || dir == "." {
		return dir
	} else {
		return internal.EnsureTrailingSep(dir, PathNameSeparator)
	}
}

func applyFileInfo(e *Entry, info os.FileInfo) {
	e.fileInfo = info
	e.readonly = info.Mode().Perm()&0200 == 0
}

// Obtains FileInfo according to flags. Returns (nil, nil) when the path
// does not exist and DetailsLater is set; otherwise returns the error.
func statInfo(path string, flags SearchFlags) (os.FileInfo, error) {
	var (
		info os.FileInfo
		err  error
	)
	if 0 != (flags & NoFollowLinks) {
		info, err = os.Lstat(path)
	} else {
		info, err = os.Stat(path)
	}
	if err == nil {
		return info, nil
	}
	if os.IsNotExist(err) && 0 != (flags&DetailsLater) {
		return nil, nil
	}
	return nil, err
}
