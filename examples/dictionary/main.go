package main

import (
	"fmt"

	"github.com/planetarium/bencodex-go/bencodextype"
)

func main() {
	dict := bencodextype.NewDictionary()
	dict.Set("spam1", []byte("eggs"))
	dict.Set("spam2", "moo")
	dict.Set("spam3", 34)

	for _, key := range dict.Keys() {
		value := dict.Get(key)
		fmt.Println(key, value)
	}
	fmt.Println()

	if dict.CanConvertToMap() {
		mapValue := dict.ConvertToMap()
		fmt.Printf("mapValue type: %T\n", mapValue)
		for key, value := range mapValue {
			fmt.Println(key, value)
		}
	}
	fmt.Println()

	dict = bencodextype.NewDictionary()
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
