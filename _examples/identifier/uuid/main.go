package main

import (
	"encoding/json"
	"fmt"

	id "github.com/golibry/go-common-domain/domain/identifier"
)

func main() {
	identifier, _ := id.NewUUIDIdentifier(" 550E8400-E29B-41D4-A716-446655440000 ")
	fromDB := id.ReconstituteUUIDIdentifier(identifier.Value())

	jsonData, _ := json.Marshal(identifier)

	fmt.Println(identifier.Value())
	fmt.Println(identifier.Equals(fromDB))
	fmt.Println(string(jsonData))
}
