package main

import (
	recls "github.com/synesissoftware/recls.Go"
	"github.com/synesissoftware/ver2go"

	"fmt"
)

func main() {
	fmt.Printf("recls v%s\n", ver2go.CalcVersionString(recls.VersionMajor, recls.VersionMinor, recls.VersionPatch, recls.VersionAB))
	fmt.Printf("ver2go v%s\n", ver2go.CalcVersionString(ver2go.VersionMajor, ver2go.VersionMinor, ver2go.VersionPatch, ver2go.VersionAB))
}
