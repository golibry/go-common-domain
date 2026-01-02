package main

import (
	"errors"
	"fmt"

	a "github.com/golibry/go-common-domain/domain/auth"
)

func main() {
	pwd, _ := a.NewPassword("MySecure123!@")

	// Never print raw hashes; String() returns a safe placeholder
	fmt.Println(pwd.String())

	// Verify the correct password
	fmt.Println(pwd.Verify("MySecure123!@") == nil)

	// Verify mismatch returns the expected error
	err := pwd.Verify("wrong")
	fmt.Println(errors.Is(err, a.ErrPasswordVerifyFailed))
}
