package recls

import "github.com/synesissoftware/ver2go"

const (
	VersionMajor uint16 = 0
	VersionMinor uint16 = 0
	VersionPatch uint16 = 0
	VersionAB    uint16 = 0x0005
	Version      uint64 = (uint64(VersionMajor) << 48) + (uint64(VersionMinor) << 32) + (uint64(VersionPatch) << 16) + (uint64(VersionAB) << 0)
)

var (
	versionString string = ver2go.CalcVersionString(VersionMajor, VersionMinor, VersionPatch, VersionAB)
)

func VersionString() string {
	return versionString
}
