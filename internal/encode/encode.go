package encode

import "reflect"

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
		if val.Type().Elem().Kind() == reflect.Uint8 {
			return encodeBytes(val), nil
		}

		return encodeList(val)
	case reflect.Map:
		return encodeDictionary(val)
	case reflect.Interface:
		return EncodeValue(val.Elem())
	default:
		return nil, &UnsupportedTypeError{val.Type()}
	}
}
