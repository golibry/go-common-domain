package main

import (
	"encoding/json"
	"fmt"

	g "github.com/golibry/go-common-domain/domain/geography"
)

func main() {
	point, _ := g.NewGeoPoint(44.4268, 26.1025)
	fromDB := g.ReconstituteGeoPoint(point.Latitude(), point.Longitude())

	jsonData, _ := json.Marshal(point)

	fmt.Println(point.String())
	fmt.Println(point.Equals(fromDB))
	fmt.Println(string(jsonData))
}
