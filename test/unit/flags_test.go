package recls_test

import (
	recls "github.com/synesissoftware/recls.Go"

	"github.com/stretchr/testify/require"

	"testing"
)

func Test_NormaliseFlags_DEFAULTS_TO_Files(t *testing.T) {
	require.Equal(t, recls.Files, recls.NormaliseFlags(0))
	require.Equal(t, recls.Files|recls.Recursive, recls.NormaliseFlags(recls.Recursive))
}

func Test_NormaliseFlags_PRESERVES_TYPES(t *testing.T) {
	flags := recls.Directories | recls.Recursive
	require.Equal(t, flags, recls.NormaliseFlags(flags))

	both := recls.Files | recls.Directories
	require.Equal(t, both, recls.NormaliseFlags(both))
}

func Test_PathNameSeparator_NON_EMPTY(t *testing.T) {
	require.NotEmpty(t, recls.PathNameSeparator)
}

func Test_PathSeparator_NON_EMPTY(t *testing.T) {
	require.NotEmpty(t, recls.PathSeparator)
}
