package main

import (
	"fmt"

	_ "bad/pkgxyz/v1" // want `expected pkgxyzV1, but import was named _`
)

func main() {
	fmt.Println("violates the namedimport rule")
}
