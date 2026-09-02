package recls

import "github.com/synesissoftware/ver2go"

const (
	VersionMajor uint16 = 0
	VersionMinor uint16 = 0
	VersionPatch uint16 = 2
	VersionAB    uint16 = ver2go.Release
)

var (
	version              = ver2go.CombineVersion(VersionMajor, VersionMinor, VersionPatch, VersionAB)
	versionString string = ver2go.CalcVersionString(VersionMajor, VersionMinor, VersionPatch, VersionAB)
)

// Version returns this library's version as a packed 64-bit integer, formed
// by ver2go.CombineVersion from VersionMajor, VersionMinor, VersionPatch,
// and VersionAB. The result is suitable for numeric comparison: a later
// release has a strictly greater value than an earlier one that uses the
// same packing.
func Version() uint64 {
	return version
}

// VersionString returns this library's version as a human-readable string,
// formed by ver2go.CalcVersionString from VersionMajor, VersionMinor,
// VersionPatch, and VersionAB. For a final (non-prerelease) version the
// result is of the form "MAJOR.MINOR.PATCH", e.g. "0.0.1".
func VersionString() string {
	return versionString
}
