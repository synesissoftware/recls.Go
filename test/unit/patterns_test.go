package recls_test

import (
	"github.com/synesissoftware/recls.Go/internal"

	"github.com/stretchr/testify/require"

	"runtime"
	"testing"
)

func Test_SplitPatterns_EMPTY_IS_WILDCARDS_ALL(t *testing.T) {
	require.Equal(t, []string{"*"}, internal.SplitPatterns("", ":"))
}

func Test_SplitPatterns_NORMATIVE(t *testing.T) {
	require.Equal(t, []string{"Gemfile"}, internal.SplitPatterns("Gemfile", "|"))
	require.Equal(t, []string{"Gemfile", "Gemfile.lock"}, internal.SplitPatterns("Gemfile|Gemfile.lock", "|"))
	require.Equal(t, []string{"*.go", "*.md"}, internal.SplitPatterns("*.go|*.md", ":"))
	require.Equal(t, []string{"*.go", "*.md"}, internal.SplitPatterns("*.go|*.md", ":"))
}

func Test_SplitPatterns_EDGE_CASES(t *testing.T) {
	require.Equal(t, []string{"Gemfile"}, internal.SplitPatterns("|Gemfile", "|"))
	require.Equal(t, []string{"Gemfile"}, internal.SplitPatterns("Gemfile|", "|"))
	require.Equal(t, []string{"Gemfile"}, internal.SplitPatterns("Gemfile||||", "|"))
	require.Equal(t, []string{"Gemfile", "Gemfile.lock"}, internal.SplitPatterns("Gemfile|||||Gemfile.lock", "|"))
}

func Test_SplitPatterns_PIPE_AND_COLON(t *testing.T) {
	require.Equal(t, []string{"*.go", "*.md"}, internal.SplitPatterns("*.go|*.md", ":"))
	require.Equal(t, []string{"*.go", "*.md"}, internal.SplitPatterns("*.go:*.md", ":"))
}

func Test_SplitPatterns_SEMICOLON_ONLY_WHEN_OnlyWhenPathListSep(t *testing.T) {
	// ';' is the Windows path-list separator; on Unix it is not special.
	require.Equal(t, []string{"*.go;*.md"}, internal.SplitPatterns("*.go;*.md", ":"))
	require.Equal(t, []string{"*.go", "*.md"}, internal.SplitPatterns("*.go;*.md", ";"))
	if runtime.GOOS == "windows" {
		require.Equal(t, ";", string(';'))
	} else {
		require.Equal(t, ";", string(';'))

	}
}

func Test_SplitPatterns_Dedup(t *testing.T) {
	require.Equal(t, []string{"*.go"}, internal.SplitPatterns("*.go|*.go", ":"))
}

func Test_NormalisePatterns_STRING(t *testing.T) {
	require.Equal(t, []string{"*.go", "*.md"}, internal.NormalisePatterns("*.go|*.md", ":"))
	require.Equal(t, []string{"*"}, internal.NormalisePatterns("", ":"))
}

func Test_NormalisePatterns_SLICE(t *testing.T) {
	require.Equal(t, []string{"*"}, internal.NormalisePatterns([]string{}, ":"))
	require.Equal(t, []string{"*"}, internal.NormalisePatterns([]string{"", "  "}, ":"))
	require.Equal(t, []string{"*.go", "*.md"}, internal.NormalisePatterns([]string{"*.go", "*.md"}, ":"))
	require.Equal(t, []string{"*.go"}, internal.NormalisePatterns([]string{"*.go", "*.go"}, ":"))
	// Slice elements are discrete: '|' / ':' inside an element are not split.
	require.Equal(t, []string{"*.go|*.md"}, internal.NormalisePatterns([]string{"*.go|*.md"}, ":"))
	require.Equal(t, []string{"a:b"}, internal.NormalisePatterns([]string{"a:b"}, ":"))
}

func Test_CompileAndMatch_EntryName(t *testing.T) {
	cp, err := internal.CompilePatterns([]string{"*.txt", "readme"})
	require.NoError(t, err)

	ok, err := cp.MatchesAny("notes.txt")
	require.NoError(t, err)
	require.True(t, ok)

	ok, err = cp.MatchesAny("readme")
	require.NoError(t, err)
	require.True(t, ok)

	ok, err = cp.MatchesAny("notes.go")
	require.NoError(t, err)
	require.False(t, ok)
}

func Test_ValidatePatternsForFlags_DotRecursive(t *testing.T) {
	require.Error(t, internal.ValidatePatternsForFlags([]string{"."}, true))
	require.Error(t, internal.ValidatePatternsForFlags([]string{".."}, true))
	require.NoError(t, internal.ValidatePatternsForFlags([]string{"."}, false))
	require.NoError(t, internal.ValidatePatternsForFlags([]string{"*.go"}, true))
}
