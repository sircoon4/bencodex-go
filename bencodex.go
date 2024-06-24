package bencodex

import (
	"fmt"
	"reflect"

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
