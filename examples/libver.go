package main

import (
	recls "github.com/synesissoftware/recls.Go"
	"github.com/synesissoftware/ver2go"

	"fmt"
)

func main() {
	fmt.Printf("recls v%s\n", recls.VersionString())
	fmt.Printf("ver2go v%s\n", ver2go.VersionString())
}
