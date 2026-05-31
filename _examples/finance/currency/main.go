package main

import (
	"fmt"

	f "github.com/golibry/go-common-domain/domain/finance"
)

func main() {
	cur, _ := f.NewCurrency(" usd ")
	jpy, _ := f.NewCurrency("jpy")
	var scanned f.Currency
	_ = scanned.Scan([]byte(" eur "))

	scale, _ := jpy.MinorUnitScale()

	fmt.Println(cur.Value())
	fmt.Println(scanned.Value())
	fmt.Println(scale)
}
