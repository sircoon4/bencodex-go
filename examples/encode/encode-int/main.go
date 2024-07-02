package main

import (
	"fmt"
	"math/big"

	"github.com/sircoon4/bencodex-go"
)

func main() {
	var val any

	val = 340
	printEncodeResult(val)

	val = -1725
	printEncodeResult(val)

	val, ok := new(big.Int).SetString("1234567890123456789012345678901234567890", 10)
	if ok {
		printEncodeResult(val)
	}

	val = uint(340)
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
