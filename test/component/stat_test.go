package recls_test

import (
	recls "github.com/synesissoftware/recls.Go"

	"github.com/stretchr/testify/require"

	"os"
	"path/filepath"
	"testing"
)

func Test_Stat_FILE(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "hello.txt")
	require.NoError(t, os.WriteFile(path, []byte("hi"), 0o644))

	e, err := recls.Stat(path, 0)
	require.NoError(t, err)
	require.True(t, e.Exists())
	require.True(t, e.IsFile())
	require.False(t, e.IsDir())
	require.Equal(t, "hello.txt", e.EntryName)
	require.Equal(t, "hello", e.Stem)
	require.Equal(t, ".txt", e.Extension)
	require.Equal(t, int64(2), e.Size())
	require.NotEmpty(t, e.Path())
	require.NotEmpty(t, e.Location)
}

func Test_Stat_Directory_MarkDirs(t *testing.T) {
	dir := t.TempDir()

	e, err := recls.Stat(dir, recls.Directories|recls.MarkDirs|recls.DirectoryParts)
	require.NoError(t, err)
	require.True(t, e.Exists())
	require.True(t, e.IsDir())
	require.True(t, len(e.Path()) > 0)
	require.Equal(t, recls.PathNameSeparator, string(e.Path()[len(e.Path())-1]))
	require.NotNil(t, e.DirectoryParts)
}

func Test_Stat_DetailsLater_MISSING(t *testing.T) {
	dir := t.TempDir()
	missing := filepath.Join(dir, "no-such-file")

	e, err := recls.Stat(missing, recls.DetailsLater|recls.DirectoryParts)
	require.NoError(t, err)
	require.False(t, e.Exists())
	require.Equal(t, "no-such-file", e.EntryName)
}

func Test_Stat_MISSING_WITHOUT_DetailsLater(t *testing.T) {
	dir := t.TempDir()
	missing := filepath.Join(dir, "no-such-file")

	_, err := recls.Stat(missing, 0)
	require.Error(t, err)
}
