package main

import (
	"fmt"
	"reflect"

	"github.com/planetarium/bencodex-go"
)

func main() {
	var b []byte
	var v any
	var err error

	b, err = bencodex.Encode(nil)
	if err == nil {
		bencodex.Decode(b)
	}
	b, err = bencodex.Encode(true)
	if err == nil {
		bencodex.Decode(b)
	}
	b, err = bencodex.Encode(false)
	if err == nil {
		bencodex.Decode(b)
	}
	b, err = bencodex.Encode(340)
	if err == nil {
		bencodex.Decode(b)
	}
	b, err = bencodex.Encode(-1725)
	if err == nil {
		bencodex.Decode(b)
	}
	b, err = bencodex.Encode(uint(340))
	if err == nil {
		bencodex.Decode(b)
	}
	b, err = bencodex.Encode("hello")
	if err == nil {
		bencodex.Decode(b)
	}
	b, err = bencodex.Encode("단팥빵")
	if err == nil {
		bencodex.Decode(b)
	}
	var emptyList []int
	b, err = bencodex.Encode(emptyList)
	if err == nil {
		bencodex.Decode(b)
	}
	b, err = bencodex.Encode([]byte("hello"))
	if err == nil {
		bencodex.Decode(b)
	}
	b, err = bencodex.Encode([]int{1, 2, 3})
	if err == nil {
		v, err = bencodex.Decode(b)
		if err == nil {
			bencodex.Encode(v)
		}
	}
	b, err = bencodex.Encode([]string{"test", "for", "string"})
	if err == nil {
		bencodex.Decode(b)
	}
	b, err = bencodex.Encode([]any{-1, "we", []byte("byteee")})
	if err == nil {
		bencodex.Decode(b)
	}
	var emptyMap map[string]int
	b, err = bencodex.Encode(emptyMap)
	if err == nil {
		bencodex.Decode(b)
	}
	b, err = bencodex.Encode(map[string]int{"one": 1, "two": 2, "three": 3})
	if err == nil {
		v, err = bencodex.Decode(b)
		if err == nil {
			bencodex.Encode(v)
		}
	}
	byteArrayKey := [3]byte{101, 101, 101}
	b, err = bencodex.Encode(map[any]any{"apam": []byte("eggs"), byteArrayKey: "moo", "spam1": []byte("eggsk")})
	if err == nil {
		v, err = bencodex.Decode(b)
		vm := v.(map[reflect.Value]any)
		fmt.Println()
		fmt.Println(bencodex.FindValue(vm, []byte{101, 101, 101}))
		fmt.Println(bencodex.FindValue(vm, byteArrayKey))
		fmt.Println(bencodex.FindValue(vm, "apam"))
		fmt.Println(bencodex.FindValue(vm, reflect.ValueOf("spam1")))
		if err == nil {
			bencodex.Encode(v)
		}
	}
}
