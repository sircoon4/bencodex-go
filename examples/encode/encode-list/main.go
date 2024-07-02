package main

import (
	"fmt"

	"github.com/sircoon4/bencodex-go"
)

func main() {
	var val any

	val = []int{}
	printEncodeResult(val)

	val = []int{1, 2, 3}
	printEncodeResult(val)

	val = []string{"test", "for", "string"}
	printEncodeResult(val)

	val = []any{1, "we", []byte("byteee")}
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
