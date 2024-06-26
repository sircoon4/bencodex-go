package main

import (
	"fmt"
	"os"

	"github.com/planetarium/bencodex-go"
	"github.com/planetarium/bencodex-go/bencodextype"
)

func main() {
	var rv any
	var err error

	testMap := map[string]any{"apam": []byte("carrots"), "spam2": "moo", "spam1": []byte("eggs")}
	f, _ := os.Create("output1.txt")
	err = bencodex.EncodeTo(f, testMap)
	f.Close()
	if err != nil {
		fmt.Println(err)
	} else {
		f, _ = os.Open("output1.txt")
		rv, err = bencodex.DecodeFrom(f)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("value: %v\n type: %T\n", rv, rv)
			fmt.Println()
		}
		f.Close()
	}

	fmt.Println()
	dict := bencodextype.NewDictionaryFromMap(testMap)
	dict.Set([]byte("spam3"), []byte("potato"))
	f, _ = os.Create("output2.txt")
	err = bencodex.EncodeTo(f, dict)
	f.Close()
	if err != nil {
		fmt.Println(err)
	} else {
		f, _ = os.Open("output2.txt")
		rv, err = bencodex.DecodeFrom(f)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("value: %v\n type: %T\n", rv, rv)
			fmt.Println()
		}
		f.Close()
	}
}
