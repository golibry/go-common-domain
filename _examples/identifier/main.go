package main

import (
	"fmt"

	id "github.com/golibry/go-common-domain/domain/identifier"
)

func main() {
	i, _ := id.NewIntIdentifier(42)
	fmt.Println(i.Value())
	fmt.Println(i.String())

	// Recreate from string
	j, _ := id.NewIntIdentifierFromString("42")
	fmt.Println(i.Equals(j))

	var scanned id.IntIdentifier
	_ = scanned.Scan(int64(42))
	fmt.Println(scanned.Value())
}
