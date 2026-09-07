package recls_test

import (
	recls "github.com/synesissoftware/recls.Go"

	"github.com/stretchr/testify/require"

	"testing"
)

const (
	Expected_VersionMajor uint16 = 0
	Expected_VersionMinor uint16 = 1
	Expected_VersionPatch uint16 = 0
	Expected_VersionAB    uint16 = 0x4001 // ver2go.Alpha1
)

func Test_Version_Elements(t *testing.T) {
	require.Equal(t, Expected_VersionMajor, recls.VersionMajor)
	require.Equal(t, Expected_VersionMinor, recls.VersionMinor)
	require.Equal(t, Expected_VersionPatch, recls.VersionPatch)
	require.Equal(t, Expected_VersionAB, recls.VersionAB)
}

func Test_Version(t *testing.T) {
	require.Equal(t, uint64(0x0000_0001_0000_4001), recls.Version())
}

func Test_Version_String(t *testing.T) {
	require.Equal(t, "0.1.0-alpha1", recls.VersionString())
}
