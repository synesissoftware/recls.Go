// Copyright 2019-2026, Matthew Wilson and Synesis Information Systems. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 * Created: 19th February 2025
 * Updated: 3rd September 2026
 */

package recls

import (
	"errors"
	"fmt"
)

// Indicates that a search has no further entries, or that an empty pattern
// string was supplied (matching C RECLS_RC_NO_MORE_DATA).
var ErrNoMoreData = errors.New("recls: no more data")

// Indicates that "." or ".." was used as a pattern in a recursive search.
var ErrDotRecursiveSearch = errors.New("recls: dot/dot-dot pattern not allowed in recursive search")

// Reports that a directory or entry could not be accessed when
// StopOnAccessFailure was set.
type AccessDeniedError struct {
	Path string
	Err  error
}

func (e *AccessDeniedError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("recls: access denied for %q: %v", e.Path, e.Err)
	} else {
		return fmt.Sprintf("recls: access denied for %q", e.Path)
	}
}

func (e *AccessDeniedError) Unwrap() error {
	return e.Err
}

// Reports that a path string could not be interpreted.
type InvalidPathError struct {
	Path string
	Err  error
}

func (e *InvalidPathError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("recls: invalid path %q: %v", e.Path, e.Err)
	} else {
		return fmt.Sprintf("recls: invalid path %q", e.Path)
	}
}

func (e *InvalidPathError) Unwrap() error {
	return e.Err
}
