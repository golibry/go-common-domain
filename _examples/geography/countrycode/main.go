package main

import (
	"fmt"

	g "github.com/golibry/go-common-domain/domain/geography"
)

func main() {
	cc, _ := g.NewCountryCode(" us ")
	fmt.Println(cc.Value())
}
