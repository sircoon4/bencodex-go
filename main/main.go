package main

import (
	"reflect"

	"github.com/planetarium/bencodex-go"
)

func main() {
	var b []byte
	var err error

	b, err = bencodex.Encode(reflect.ValueOf(nil))
	if err == nil {
		bencodex.Decode(b)
	}
	b, err = bencodex.Encode(reflect.Value{})
	if err == nil {
		bencodex.Decode(b)
	}
	b, err = bencodex.Encode(reflect.ValueOf(true))
	if err == nil {
		bencodex.Decode(b)
	}
	b, err = bencodex.Encode(reflect.ValueOf(false))
	if err == nil {
		bencodex.Decode(b)
	}
	b, err = bencodex.Encode(reflect.ValueOf(340))
	if err == nil {
		bencodex.Decode(b)
	}
	b, err = bencodex.Encode(reflect.ValueOf(-1725))
	if err == nil {
		bencodex.Decode(b)
	}
	b, err = bencodex.Encode(reflect.ValueOf(uint(340)))
	if err == nil {
		bencodex.Decode(b)
	}
	b, err = bencodex.Encode(reflect.ValueOf("hello"))
	if err == nil {
		bencodex.Decode(b)
	}
	b, err = bencodex.Encode(reflect.ValueOf("단팥빵"))
	if err == nil {
		bencodex.Decode(b)
	}
	var emptyList []int
	b, err = bencodex.Encode(reflect.ValueOf(emptyList))
	if err == nil {
		bencodex.Decode(b)
	}
	b, err = bencodex.Encode(reflect.ValueOf([]byte("hello")))
	if err == nil {
		bencodex.Decode(b)
	}
	b, err = bencodex.Encode(reflect.ValueOf([]int{1, 2, 3}))
	if err == nil {
		bencodex.Decode(b)
	}
	b, err = bencodex.Encode(reflect.ValueOf([]any{-1, "we", []byte("byteee")}))
	if err == nil {
		bencodex.Decode(b)
	}
	var emptyMap map[string]int
	b, err = bencodex.Encode(reflect.ValueOf(emptyMap))
	if err == nil {
		bencodex.Decode(b)
	}
	b, err = bencodex.Encode(reflect.ValueOf(map[string]int{"one": 1, "two": 2, "three": 3}))
	if err == nil {
		bencodex.Decode(b)
	}
	byteArrayKey := [3]byte{101, 101, 101}
	b, err = bencodex.Encode(reflect.ValueOf(map[any]any{"apam": []byte("eggs"), byteArrayKey: "moo", "spam1": []byte("eggs")}))
	if err == nil {
		bencodex.Decode(b)
	}

}
