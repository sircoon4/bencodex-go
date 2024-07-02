package main

import (
	"fmt"

	"github.com/sircoon4/bencodex-go"
)

func main() {
	var b []byte

	b = []byte("le")
	printDecodeResult(b)

	b = []byte("li1ei2ei3ee")
	printDecodeResult(b)

	b = []byte("lu4:testu3:foru6:stringe")
	printDecodeResult(b)

	b = []byte("li-1eu2:we6:byteeee")
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
