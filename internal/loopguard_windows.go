// Copyright 2019-2026, Matthew Wilson and Synesis Information Systems. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

/*
 * Created: 3rd September 2026
 * Updated: 3rd September 2026
 */

package internal

// No-op on Windows (reparse-point loops are out of scope for v0.1).
type LoopGuard struct{}

// Constructs a no-op guard.
func NewLoopGuard() *LoopGuard {
	return &LoopGuard{}
}

// Always returns true on Windows.
func (g *LoopGuard) Enter(path string) bool {
	_ = path

	return true
}
