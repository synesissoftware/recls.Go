package recls

import (
	. "github.com/synesissoftware/recls.Go"

	"github.com/stretchr/testify/require"

	"testing"
)

const (
	Expected_VersionMajor uint16 = 0
	Expected_VersionMinor uint16 = 0
	Expected_VersionPatch uint16 = 0
	Expected_VersionAB    uint16 = 5
)

func Test_Version_Elements(t *testing.T) {
	require.Equal(t, Expected_VersionMajor, VersionMajor)
	require.Equal(t, Expected_VersionMinor, VersionMinor)
	require.Equal(t, Expected_VersionPatch, VersionPatch)
	require.Equal(t, Expected_VersionAB, VersionAB)
}

func Test_Version(t *testing.T) {
	require.Equal(t, uint64(5), Version)
}

func Test_Version_String(t *testing.T) {
	require.Equal(t, "0.0.0", VersionString())
}
