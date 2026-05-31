package main

import (
	"encoding/json"
	"fmt"

	g "github.com/golibry/go-common-domain/domain/geography"
)

func main() {
	country, _ := g.NewCountryCode("ro")
	address, _ := g.NewAddress(" Main   Street 1 ", "Apt 2", "Bucharest", "Bucuresti", "010101", country)
	fromDB := g.ReconstituteAddress(address.Line1(), address.Line2(), address.City(), address.Region(), address.PostalCode(), address.Country())

	jsonData, _ := json.Marshal(address)

	fmt.Println(address.String())
	fmt.Println(address.Equals(fromDB))
	fmt.Println(string(jsonData))
}
