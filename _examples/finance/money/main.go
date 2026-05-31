package main

import (
	"encoding/json"
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
	rounded, _ := m1.DivideWithRounding(decimal.NewFromInt(3), f.RoundHalfUp)
	eur, _ := f.NewCurrency("EUR")
	minor, _ := f.NewMoneyFromMinorUnits(1099, eur, 2)
	jpy, _ := f.NewMoneyFromString("1000", "JPY")
	kwd, _ := f.NewMoneyFromString("10.999", "KWD")

	vat, _ := f.NewPercentageRateFromString("19")
	vatAmount, _ := vat.ApplyTo(m1)
	gross, _ := vat.AddTo(m1)

	jsonValue, _ := json.Marshal(m1)
	fromDB := f.ReconstituteMoneyFromMinorUnits(1050, f.ReconstituteCurrency("USD"), 2)

	fmt.Println(sum.String())
	fmt.Println(diff.String())
	fmt.Println(prod.String())
	fmt.Println(rounded.String())
	fmt.Println(minor.String())
	fmt.Println(jpy.String(), jpy.Scale())
	fmt.Println(kwd.String(), kwd.Scale())
	fmt.Println(vatAmount.String())
	fmt.Println(gross.String())
	fmt.Println(string(jsonValue))
	fmt.Println(fromDB.String())
}
