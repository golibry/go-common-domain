package main

import (
	"fmt"

	w "github.com/golibry/go-common-domain/domain/web"
)

func main() {
	ip4, _ := w.NewIPAddress("  192.168.001.010  ")
	fmt.Println(ip4.Value())
	fmt.Println(ip4.IsIPv4())
	fmt.Println(ip4.IsIPv6())
}
