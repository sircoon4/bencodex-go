package main

import (
	"fmt"
	"math/big"

	"github.com/sircoon4/bencodex-go"
	"github.com/sircoon4/bencodex-go/bencodextype"
)

func main() {
	var val any
	var b []byte
	var err error

	val = nil
	b, err = bencodex.Encode(val)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Encode Result of %v: \n", val)
		fmt.Println(b, string(b))
		fmt.Println()
	}

	val = true
	b, err = bencodex.Encode(val)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Encode Result of %v: \n", val)
		fmt.Println(b, string(b))
		fmt.Println()
	}

	val = false
	b, err = bencodex.Encode(val)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Encode Result of %v: \n", val)
		fmt.Println(b, string(b))
		fmt.Println()
	}

	val = 340
	b, err = bencodex.Encode(val)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Encode Result of %v: \n", val)
		fmt.Println(b, string(b))
		fmt.Println()
	}

	val = -1725
	b, err = bencodex.Encode(val)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Encode Result of %v: \n", val)
		fmt.Println(b, string(b))
		fmt.Println()
	}

	val = new(big.Int).SetInt64(-1725)
	b, err = bencodex.Encode(val)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Encode Result of %v: \n", val)
		fmt.Println(b, string(b))
		fmt.Println()
	}

	val = uint(340)
	b, err = bencodex.Encode(val)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Encode Result of %v: \n", val)
		fmt.Println(b, string(b))
		fmt.Println()
	}

	val = "hello"
	b, err = bencodex.Encode(val)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Encode Result of %v: \n", val)
		fmt.Println(b, string(b))
		fmt.Println()
	}

	val = "단팥빵"
	b, err = bencodex.Encode(val)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Encode Result of %v: \n", val)
		fmt.Println(b, string(b))
		fmt.Println()
	}

	val = []int{}
	b, err = bencodex.Encode(val)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Encode Result of %v: \n", val)
		fmt.Println(b, string(b))
		fmt.Println()
	}

	val = []byte("hello")
	b, err = bencodex.Encode(val)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Encode Result of %v: \n", val)
		fmt.Println(b, string(b))
		fmt.Println()
	}

	val = []int{1, 2, 3}
	b, err = bencodex.Encode(val)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Encode Result of %v: \n", val)
		fmt.Println(b, string(b))
		fmt.Println()
	}

	val = []string{"test", "for", "string"}
	b, err = bencodex.Encode(val)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Encode Result of %v: \n", val)
		fmt.Println(b, string(b))
		fmt.Println()
	}

	val = []any{1, "we", []byte("byteee")}
	b, err = bencodex.Encode(val)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Encode Result of %v: \n", val)
		fmt.Println(b, string(b))
		fmt.Println()
	}

	val = map[string]int{}
	b, err = bencodex.Encode(val)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Encode Result of %v: \n", val)
		fmt.Println(b, string(b))
		fmt.Println()
	}

	val = map[string]int{"one": 1, "two": 2, "three": 3}
	b, err = bencodex.Encode(val)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Encode Result of %v: \n", val)
		fmt.Println(b, string(b))
		fmt.Println()
	}

	val = map[string]any{"apam": []byte("carrots"), "spam2": "moo", "spam1": []byte("eggs")}
	b, err = bencodex.Encode(val)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Encode Result of %v: \n", val)
		fmt.Println(b, string(b))
		fmt.Println()
	}

	dict := bencodextype.NewDictionaryFromMap(val.(map[string]any))
	dict.Set([]byte("spam3"), []byte("potato"))
	b, err = bencodex.Encode(dict)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Encode Result of %v: \n", dict)
		fmt.Println(b, string(b))
		fmt.Println()
	}
}
