package main

import (
	libpath "github.com/synesissoftware/libpath.Go"
	recls "github.com/synesissoftware/recls.Go"
	shwild "github.com/synesissoftware/shwild.Go"
	ver2go "github.com/synesissoftware/ver2go"

	"fmt"
)

func main() {
	fmt.Printf("recls v%s\n", recls.VersionString())
	fmt.Printf("libpath v%s\n", libpath.VersionString())
	fmt.Printf("shwild v%s\n", shwild.VersionString())
	fmt.Printf("ver2go v%s\n", ver2go.VersionString())
}
