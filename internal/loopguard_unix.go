// Copyright 2019-2026, Matthew Wilson and Synesis Information Systems. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build unix

/*
 * Created: 3rd September 2026
 * Updated: 3rd September 2026
 */

package internal

import (
	"os"
	"syscall"
)

// A device+inode pair used to break symlink / hard-link directory cycles.
type DirIdentity struct {
	Dev uint64
	Ino uint64
}

// Returns the device+inode identity for path, or ok=false.
func IdentityOf(path string) (DirIdentity, bool) {
	fi, err := os.Lstat(path)
	if err != nil {
		return DirIdentity{}, false
	}
	st, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return DirIdentity{}, false
	}
	return DirIdentity{Dev: uint64(st.Dev), Ino: uint64(st.Ino)}, true
}

// Remembers visited directory identities so recursive search does not
// follow a cycle (symlink or bind-mount) forever. Matches C recls
// ReclsFileSearchDirectoryControlPreventInfiniteLoops. Disabled when
// NoBreakInfiniteLoops is set. No-op on Windows (see loopguard_windows.go).
type LoopGuard struct {
	seen map[DirIdentity]struct{}
}

// Constructs an empty guard.
func NewLoopGuard() *LoopGuard {
	return &LoopGuard{seen: make(map[DirIdentity]struct{})}
}

// Records path's identity. Returns false if already seen (cycle).
func (g *LoopGuard) Enter(path string) bool {
	id, ok := IdentityOf(path)
	if !ok {
		return true
	}
	if _, exists := g.seen[id]; exists {
		return false
	}
	g.seen[id] = struct{}{}
	return true
}
