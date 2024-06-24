package main

import (
	"reflect"

	"github.com/planetarium/bencodex-go"
)

func main() {
	bencodex.Encode(reflect.ValueOf(true))
	bencodex.Encode(reflect.ValueOf(false))
	bencodex.Encode(reflect.ValueOf(340))
	bencodex.Encode(reflect.ValueOf(-1725))
	bencodex.Encode(reflect.ValueOf(uint(340)))
	bencodex.Encode(reflect.ValueOf("hello"))
	bencodex.Encode(reflect.ValueOf("단팥빵"))
	bencodex.Encode(reflect.ValueOf([]byte("hello")))
	bencodex.Encode(reflect.ValueOf([]int{1, 2, 3}))
	bencodex.Encode(reflect.ValueOf([]any{-1, "we", []byte("byteee")}))
	bencodex.Encode(reflect.ValueOf(map[string]int{"one": 1, "two": 2, "three": 3}))
	byteArrayKey := [3]byte{101, 101, 101}
	bencodex.Encode(reflect.ValueOf(map[any]any{"apam": []byte("eggs"), byteArrayKey: "moo", "spam1": []byte("eggs")}))

}
