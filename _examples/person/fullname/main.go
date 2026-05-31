package main

import (
	"encoding/json"
	"fmt"

	p "github.com/golibry/go-common-domain/domain/person"
)

func main() {
	fn, err := p.NewFullName(" John ", "F.", " Doe ")
	fmt.Println(err)
	fmt.Println(fn.FirstName())
	fmt.Println(fn.MiddleName())
	fmt.Println(fn.LastName())
	fmt.Println(fn.String())

	jsonValue, _ := json.Marshal(fn)
	fromDB := p.ReconstituteFullName("John", "F.", "Doe")

	fmt.Println(string(jsonValue))
	fmt.Println(fromDB.String())
}
