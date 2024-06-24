package bencodex

import (
	"fmt"
	"reflect"

	"github.com/planetarium/bencodex-go/internal/decode"
	"github.com/planetarium/bencodex-go/internal/encode"
)

func Encode(val reflect.Value) ([]byte, error) {
	fmt.Printf("Encode Result of %v: \n", val)
	encodedValue, err := encode.EncodeValue(val)
	if err != nil {
		fmt.Println(err)
		return encodedValue, err
	}
	fmt.Println(encodedValue, string(encodedValue))
	return encodedValue, nil
}

func Decode(b []byte) (any, error) {
	fmt.Printf("Decode Result of %s: \n", string(b))
	decodedValue, err := decode.DecodeValue(&b)
	if err != nil {
		fmt.Println(err)
		return reflect.Value{}, err
	}
	fmt.Println(decodedValue)
	return decodedValue, nil
}
