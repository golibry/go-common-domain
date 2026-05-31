package main

import (
	"encoding/json"
	"fmt"

	pid "github.com/golibry/go-common-domain/domain/product/identifier"
)

func main() {
	ean, _ := pid.NewEAN("4006 3813 3393 1")
	gtin, _ := pid.NewGTIN("10012345678902")
	cnk, _ := pid.NewCNK("123-4566")

	fromDB := pid.ReconstituteGTIN(gtin.Value())
	jsonData, _ := json.Marshal(ean)

	fmt.Println(ean.Value())
	fmt.Println(gtin.Value())
	fmt.Println(cnk.Value())
	fmt.Println(gtin.Equals(fromDB))
	fmt.Println(string(jsonData))
}
