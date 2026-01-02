package main

import (
	"fmt"

	f "github.com/golibry/go-common-domain/domain/finance"
)

func main() {
	cur, _ := f.NewCurrency(" usd ")
	fmt.Println(cur.Value())
}
