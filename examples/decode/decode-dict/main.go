package main

import (
	"fmt"

	"github.com/sircoon4/bencodex-go"
)

func main() {
	var b []byte

	b = []byte("du3:onei1eu5:threei3eu3:twoi2ee")
	printDecodeResult(b)

	b = []byte("du4:apam7:carrotsu5:spam14:eggsu5:spam2u3:mooe")
	printDecodeResult(b)

	b = []byte("d5:spam36:potatou4:apam7:carrotsu5:spam14:eggsu5:spam2u3:mooe")
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
