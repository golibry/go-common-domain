package main

import (
	"fmt"

	f "github.com/golibry/go-common-domain/domain/finance"
)

func main() {
	rate, _ := f.NewPercentageRateFromString("19.5")
	money, _ := f.NewMoneyFromString("100.00", "EUR")

	amount, _ := rate.ApplyTo(money)
	total, _ := rate.AddTo(money)
	discounted, _ := rate.SubtractFrom(money)
	basisPoints := rate.BasisPoints()

	fromDB := f.ReconstitutePercentageRate(basisPoints)

	fmt.Println(rate.BasisPoints())
	fmt.Println(rate.Percent().String())
	fmt.Println(rate.Fraction().String())
	fmt.Println(amount.String())
	fmt.Println(total.String())
	fmt.Println(discounted.String())
	fmt.Println(fromDB.String())
}
