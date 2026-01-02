package main

import (
	"errors"
	"fmt"

	d "github.com/golibry/go-common-domain/domain"
)

func main() {
	base := errors.New("base err")
	wrapped := d.NewErrorWithWrap(base, "context err")

	fmt.Println(errors.Is(wrapped, base))
	fmt.Println(wrapped.Error())
}
