package main

import (
	"fmt"

	f "github.com/golibry/go-common-domain/domain/finance"
)

func main() {
	cur, _ := f.NewCurrency(" usd ")
	var scanned f.Currency
	_ = scanned.Scan([]byte(" eur "))

	fmt.Println(cur.Value())
	fmt.Println(scanned.Value())
}
