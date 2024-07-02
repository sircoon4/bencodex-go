package main

import (
	"fmt"

	"github.com/sircoon4/bencodex-go"
)

func main() {
	var b []byte

	b = []byte("i340e")
	printDecodeResult(b)

	b = []byte("i-1725e")
	printDecodeResult(b)

	b = []byte("i123456789012345678901234567890e")
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
