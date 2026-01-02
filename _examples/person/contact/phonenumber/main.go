package main

import (
	"fmt"

	c "github.com/golibry/go-common-domain/domain/person/contact"
)

func main() {
	pn, _ := c.NewPhoneNumber(" +1 (234) 567-890 ")
	fmt.Println(pn.Value())
	fmt.Println(pn.String())
}
