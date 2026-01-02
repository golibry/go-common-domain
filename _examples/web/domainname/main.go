package main

import (
	"fmt"

	w "github.com/golibry/go-common-domain/domain/web"
)

func main() {
	d, _ := w.NewDomainName("Sub.Example.COM")
	fmt.Println(d.String())
	fmt.Println(d.Value())
}
