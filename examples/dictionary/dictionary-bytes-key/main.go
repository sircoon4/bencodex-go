package main

import (
	"fmt"

	"github.com/sircoon4/bencodex-go/bencodextype"
)

func main() {
	dict := bencodextype.NewDictionary()
	dict.Set("spam1", []byte("eggs"))
	dict.Set([]byte("spam2"), "moo")
	dict.Set("spam3", 34)

	for _, key := range dict.Keys() {
		value := dict.Get(key)
		fmt.Println(key, value)
	}
	fmt.Println()

	if !dict.CanConvertToMap() {
		fmt.Println("Cannot convert to map")
	}
}
