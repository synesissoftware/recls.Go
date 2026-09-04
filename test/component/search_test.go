package recls_test

import (
	recls "github.com/synesissoftware/recls.Go"

	"github.com/stretchr/testify/require"

	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func makeTree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "sub", "deep"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "a.txt"), []byte("a"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "b.go"), []byte("b"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "sub", "c.txt"), []byte("c"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "sub", "deep", "d.txt"), []byte("d"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, ".hidden"), []byte("h"), 0o644))
	require.NoError(t, os.Mkdir(filepath.Join(root, ".hiddendir"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, ".hiddendir", "x.txt"), []byte("x"), 0o644))
	return root
}

func collect(
	t *testing.T,
	root string,
	patterns string,
	opts recls.SearchOptions,
) []recls.Entry {

	t.Helper()
	var out []recls.Entry
	err := recls.SearchFunc(root, patterns, opts, func(e recls.Entry) error {
		out = append(out, e)
		return nil
	})
	require.NoError(t, err)
	return out
}

func names(entries []recls.Entry) []string {
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		out = append(out, e.EntryName)
	}
	return out
}

func Test_Search_NON_RECURSIVE_PATTERN(t *testing.T) {
	root := makeTree(t)
	got := collect(t, root, "*.txt", recls.SearchOptions{})
	require.Equal(t, []string{"a.txt"}, names(got))
}

func Test_Search_Recursive_PATTERN(t *testing.T) {
	root := makeTree(t)
	got := collect(t, root, "*.txt", recls.SearchOptions{Flags: recls.Recursive})
	// Without IgnoreHiddenEntries, .hiddendir/x.txt is included.
	require.ElementsMatch(t, []string{"a.txt", "c.txt", "d.txt", "x.txt"}, names(got))
}

func Test_Search_Directories_AND_MarkDirs(t *testing.T) {
	root := makeTree(t)
	got := collect(t, root, "*", recls.SearchOptions{
		Flags: recls.Directories | recls.Recursive | recls.MarkDirs,
	})
	require.NotEmpty(t, got)
	for _, e := range got {
		require.True(t, e.IsDir())
		require.Equal(t, recls.PathNameSeparator, string(e.Path()[len(e.Path())-1]))
		require.NotEmpty(t, e.SearchRelativePath)
	}
}

func Test_Search_IgnoreHidden(t *testing.T) {
	root := makeTree(t)
	got := collect(t, root, "*", recls.SearchOptions{
		Flags: recls.Files | recls.Recursive | recls.IgnoreHiddenEntries,
	})
	for _, e := range got {
		require.False(t, e.IsHidden(), "unexpected hidden entry %q", e.Path())
		require.NotEqual(t, ".hidden", e.EntryName)
		require.NotEqual(t, "x.txt", e.EntryName)
	}
	require.ElementsMatch(t, []string{"a.txt", "b.go", "c.txt", "d.txt"}, names(got))
}

func Test_Search_MULTI_PATTERN(t *testing.T) {
	root := makeTree(t)
	got := collect(t, root, "*.txt|*.go", recls.SearchOptions{Flags: recls.Recursive})
	require.ElementsMatch(t, []string{"a.txt", "b.go", "c.txt", "d.txt", "x.txt"}, names(got))
}

func Test_Search_iter(t *testing.T) {
	root := makeTree(t)
	var count int
	for e, err := range recls.Search(root, "*.txt", recls.SearchOptions{Flags: recls.Recursive}) {
		require.NoError(t, err)
		require.Equal(t, ".txt", e.Extension)
		count++
	}
	require.Equal(t, 4, count) // a,c,d,x (hidden dir still traversed without IgnoreHidden)
}

func Test_Search_StopOnAccessFailure(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("chmod-based access denial is Unix-specific")
	}
	root := t.TempDir()
	denied := filepath.Join(root, "denied")
	require.NoError(t, os.Mkdir(denied, 0o000))
	t.Cleanup(func() { _ = os.Chmod(denied, 0o755) })

	err := recls.SearchFunc(root, "*", recls.SearchOptions{
		Flags: recls.Recursive | recls.StopOnAccessFailure | recls.Files | recls.Directories,
	}, func(e recls.Entry) error {
		return nil
	})
	var ade *recls.AccessDeniedError
	require.True(t, errors.As(err, &ade))
	require.NotEmpty(t, ade.Path)
	require.Error(t, ade.Err)
}

func Test_Search_Symlink_NoFollow(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink policy coverage focused on Unix")
	}
	root := t.TempDir()
	outside := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(outside, "inside.txt"), []byte("i"), 0o644))
	link := filepath.Join(root, "linkdir")
	require.NoError(t, os.Symlink(outside, link))
	require.NoError(t, os.WriteFile(filepath.Join(root, "local.txt"), []byte("l"), 0o644))

	got := collect(t, root, "*", recls.SearchOptions{
		Flags: recls.Files | recls.Directories | recls.Links | recls.Recursive | recls.NoFollowLinks,
	})
	var namesFound []string
	for _, e := range got {
		namesFound = append(namesFound, e.EntryName)
	}
	require.Contains(t, namesFound, "local.txt")
	require.NotContains(t, namesFound, "inside.txt")
}
