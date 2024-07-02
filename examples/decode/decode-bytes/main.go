package main

import (
	"fmt"

	"github.com/sircoon4/bencodex-go"
)

func main() {
	b := []byte("5:hello")
	printDecodeResult(b)
}

func printDecodeResult(b []byte) {
	rv, err := bencodex.Decode(b)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Decode Result of %s: \n", string(b))
		fmt.Printf("value: %v\n type: %T\n", rv, rv)
		fmt.Println()
	}
}
