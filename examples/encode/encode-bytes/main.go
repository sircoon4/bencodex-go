package main

import (
	"fmt"

	"github.com/sircoon4/bencodex-go"
)

func main() {
	val := []byte("hello")
	printEncodeResult(val)
}

func printEncodeResult(val any) {
	b, err := bencodex.Encode(val)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Encode Result of %v: \n", val)
		fmt.Println(string(b))
		fmt.Println()
	}
}
