package recls_test

import (
	recls "github.com/synesissoftware/recls.Go"

	"github.com/stretchr/testify/require"

	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func Test_Entry_LinkCount_FILE(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "a.txt")
	require.NoError(t, os.WriteFile(path, []byte("x"), 0o644))

	e, err := recls.Stat(path, 0)
	require.NoError(t, err)
	require.Equal(t, uint64(1), e.LinkCount())

	linkPath := filepath.Join(root, "a-link.txt")
	if err := os.Link(path, linkPath); err != nil {
		if runtime.GOOS == "windows" {
			t.Skipf("os.Link unavailable: %v", err)
		}
		require.NoError(t, err)
	}

	e, err = recls.Stat(path, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, e.LinkCount(), uint64(2))

	e2, err := recls.Stat(linkPath, 0)
	require.NoError(t, err)
	require.Equal(t, e.LinkCount(), e2.LinkCount())
}

func Test_Entry_LinkCount_DIRECTORY(t *testing.T) {
	root := t.TempDir()
	e, err := recls.Stat(root, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, e.LinkCount(), uint64(1))
	if runtime.GOOS != "windows" {
		// Unix directories report at least 2 (`.` and `..`).
		require.GreaterOrEqual(t, e.LinkCount(), uint64(2))
	}
}
