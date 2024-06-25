package bencodex

import (
	"fmt"
	"reflect"

	"github.com/planetarium/bencodex-go/internal/decode"
	"github.com/planetarium/bencodex-go/internal/encode"
)

func Encode(val any) ([]byte, error) {
	fmt.Printf("Encode Result of %v: \n", val)
	encodedValue, err := encode.EncodeValue(reflect.ValueOf(val))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(encodedValue, string(encodedValue))
	return encodedValue, nil
}

func Decode(b []byte) (any, error) {
	fmt.Printf("Decode Result of %s: \n", string(b))
	val, err := decode.DecodeValue(&b)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var decodedValue any
	if !val.IsValid() {
		decodedValue = nil
	} else {
		decodedValue = val.Interface()
	}
	fmt.Println(decodedValue)
	return decodedValue, nil
}
