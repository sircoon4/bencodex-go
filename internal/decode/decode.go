package decode

import (
	"fmt"
	"reflect"
)

func DecodeValue(b *[]byte) (reflect.Value, error) {
	if *b == nil {
		return reflect.Value{}, fmt.Errorf("nil byte array")
	}

	switch (*b)[0] {
	case 'n':
		return decodeNil(b)
	case 't', 'f':
		return decodeBool(b)
	case 'i':
		return decodeInteger(b)
	case 'u':
		return decodeString(b)
	case 'l':
		return decodeList(b)
	case 'd':
		return decodeDictionary(b)
	default:
		return decodeBytes(b)
	}
}
