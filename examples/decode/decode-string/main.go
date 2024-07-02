package main

import (
	"fmt"

	"github.com/sircoon4/bencodex-go"
)

func main() {
	var b []byte

	b = []byte("u5:hello")
	printDecodeResult(b)

	b = []byte("u9:단팥빵")
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
