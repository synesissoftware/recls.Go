package recls_test

import (
	"github.com/synesissoftware/recls.Go/internal"

	"github.com/stretchr/testify/require"

	"testing"
)

func Test_DeriveRelativePath_SAME_DIR(t *testing.T) {
	sep := "/"
	origin := "/Users/matthewwilson/dev/freelibs/recls/100/recls.Ruby/trunk"
	path := "/Users/matthewwilson/dev/freelibs/recls/100/recls.Ruby/trunk"
	require.Equal(t, ".", internal.DeriveRelativePath(origin, path, sep))
	require.Equal(t, "./", internal.DeriveRelativePath(origin, path+"/", sep))
}

func Test_DeriveRelativePath_UNDER_HOME(t *testing.T) {
	// libpath.Ruby test_recls_stat_case_2 shape: search root = home,
	// entry = pwd under home.
	sep := "/"
	home := "/Users/matthewwilson"
	pwd := "/Users/matthewwilson/dev/freelibs/recls/100/recls.Ruby/trunk"
	got := internal.DeriveRelativePath(home, pwd+"/", sep)
	require.Equal(t, "dev/freelibs/recls/100/recls.Ruby/trunk/", got)
}

func Test_DeriveRelativePath_UP_TO_HOME(t *testing.T) {
	// libpath.Ruby test_recls_stat_case_3 shape: search root = pwd,
	// entry = home (ancestor).
	sep := "/"
	home := "/Users/matthewwilson"
	pwd := "/Users/matthewwilson/dev/freelibs/recls/100/recls.Ruby/trunk"
	got := internal.DeriveRelativePath(pwd, home+"/", sep)
	require.Equal(t, "../../../../../../", got)
}

func Test_DeriveRelativePath_EMPTY_ORIGIN(t *testing.T) {
	require.Equal(t, "/tmp/x", internal.DeriveRelativePath("", "/tmp/x", "/"))
}

func Test_EnsureTrailingSep_TRAILING_SEP(t *testing.T) {
	require.Equal(t, "/tmp/", internal.EnsureTrailingSep("/tmp", "/"))
	require.Equal(t, "/tmp/", internal.EnsureTrailingSep("/tmp/", "/"))
}

func Test_IsHiddenName_DOT_DOT(t *testing.T) {
	require.False(t, internal.IsHiddenName("."))
	require.False(t, internal.IsHiddenName(".."))
	require.False(t, internal.IsHiddenName("visible"))
	require.True(t, internal.IsHiddenName(".hidden"))
}

func Test_SplitDirectoryParts_DIRECTORY_PARTS(t *testing.T) {
	parts := internal.SplitDirectoryParts("dev/freelibs/trunk/", "/")
	require.Equal(t, []string{"dev/", "freelibs/", "trunk/"}, parts)
}
