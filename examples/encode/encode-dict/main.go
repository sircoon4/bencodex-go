package main

import (
	"fmt"

	"github.com/sircoon4/bencodex-go"
	"github.com/sircoon4/bencodex-go/bencodextype"
)

func main() {
	var val any

	val = map[string]int{}
	printEncodeResult(val)

	val = map[string]int{"one": 1, "two": 2, "three": 3}
	printEncodeResult(val)

	val = map[string]any{"apam": []byte("carrots"), "spam2": "moo", "spam1": []byte("eggs")}
	printEncodeResult(val)

	dict := bencodextype.NewDictionaryFromMap(val.(map[string]any))
	dict.Set([]byte("spam3"), []byte("potato"))
	printEncodeResult(dict)
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
