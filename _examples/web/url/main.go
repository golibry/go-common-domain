package main

import (
	"fmt"

	w "github.com/golibry/go-common-domain/domain/web"
)

func main() {
	u, _ := w.NewURL("https://example.com/path?q=1")
	fmt.Println(u.Value())
	fmt.Println(u.Parsed().Scheme)
	fmt.Println(u.Parsed().Host)
	fmt.Println(u.Parsed().Path)
	fmt.Println(u.Parsed().RawQuery)
}
