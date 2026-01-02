package main

import (
	"fmt"

	f "github.com/golibry/go-common-domain/domain/finance"
	"github.com/shopspring/decimal"
)

func main() {
	m1, _ := f.NewMoneyFromString("10.50", "USD")
	m2, _ := f.NewMoneyFromString("2.25", "USD")

	sum, _ := m1.Add(m2)
	diff, _ := m1.Subtract(m2)
	prod, _ := m2.Multiply(decimal.NewFromInt(3))

	fmt.Println(sum.String())
	fmt.Println(diff.String())
	fmt.Println(prod.String())
}
