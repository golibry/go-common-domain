package main

import (
	"fmt"

	w "github.com/golibry/go-common-domain/domain/web"
)

func main() {
	e, _ := w.NewEmail("User.Name+tag@Example.COM")
	fmt.Println(e.LocalPart())
	fmt.Println(e.DomainPart())
	fmt.Println(e.String())
	fmt.Println(e.Value())
}
