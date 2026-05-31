package main

import (
	"fmt"

	f "github.com/golibry/go-common-domain/domain/finance"
)

func main() {
	cur, _ := f.NewCurrency(" usd ")
	jpy, _ := f.NewCurrency("jpy")
	fromDB := f.ReconstituteCurrency("EUR")

	scale, _ := jpy.MinorUnitScale()

	fmt.Println(cur.Value())
	fmt.Println(fromDB.Value())
	fmt.Println(scale)
}
