package main

import (
	"fmt"

	"github.com/sircoon4/bencodex-go/bencodextype"
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

	subDict := bencodextype.NewDictionary()
	subDict.Set("spam1", []byte("eggs"))
	subDict.Set("spam2", "moo")
	subDict.Set("spam3", 34)
	dict.Set("spam4", subDict)
	if dict.CanConvertToMap() {
		mapValue := dict.ConvertToMap()
		fmt.Printf("mapValue type: %T\n", mapValue)
		for key, value := range mapValue {
			fmt.Println(key, value)
		}
	}
	fmt.Println()
}
