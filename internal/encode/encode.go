package encode

import (
	"math/big"
	"reflect"

	"github.com/sircoon4/bencodex-go/bencodextype"
)

func EncodeValue(val reflect.Value) ([]byte, error) {
	if isNil(val) {
		return encodeNil(), nil
	}

	switch val.Kind() {
	case reflect.Bool:
		return encodeBool(val), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return encodeInteger(val), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return encodeUnsignedInteger(val), nil
	case reflect.String:
		return encodeString(val), nil
	case reflect.Slice, reflect.Array:
		// If the slice is a byte slice, encode it as a byte slice
		if val.Type().Elem().Kind() == reflect.Uint8 {
			return encodeBytes(val), nil
		}

		return encodeList(val)
	case reflect.Map:
		return encodeMap(val)
	case reflect.Interface:
		//handle elements of []any type
		return EncodeValue(val.Elem())
	case reflect.Pointer:
		_, ok := val.Interface().(*bencodextype.Dictionary)
		if ok {
			return encodeDictionary(val)
		}
		_, ok = val.Interface().(*big.Int)
		if ok {
			return encodeBigInteger(val), nil
		}
		return nil, &UnsupportedTypeError{val.Type(), val.Kind()}
	default:
		return nil, &UnsupportedTypeError{val.Type(), val.Kind()}
	}
}
