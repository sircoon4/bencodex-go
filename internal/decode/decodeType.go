package decode

import (
	"fmt"
	"math/big"
	"reflect"
	"strconv"

	"github.com/sircoon4/bencodex-go/bencodextype"
)

func decodeNil(b *[]byte) (reflect.Value, error) {
	_, err := popByte(b)
	if err != nil {
		return reflect.Value{}, err
	}
	return reflect.Value{}, nil
}

func decodeBool(b *[]byte) (reflect.Value, error) {
	ch, err := popByte(b)
	if err != nil {
		return reflect.Value{}, err
	}

	if ch == 't' {
		return reflect.ValueOf(true), nil
	} else if ch == 'f' {
		return reflect.ValueOf(false), nil
	} else {
		return reflect.Value{}, fmt.Errorf("data is not compatible with bencodex")
	}
}

func decodeInteger(b *[]byte) (reflect.Value, error) {
	_, err := popByte(b)
	if err != nil {
		return reflect.Value{}, err
	}

	v, err := popBytesUntil(b, 'e')
	if err != nil {
		return reflect.Value{}, err
	}

	vInt, err := strconv.Atoi(string(v))
	if err != nil {
		// If the value is too large to fit in an int, it is stored as a big.Int
		vBigInt := new(big.Int)
		vBigInt, ok := vBigInt.SetString(string(v), 10)
		if !ok {
			return reflect.Value{}, err
		} else {
			return reflect.ValueOf(vBigInt), nil
		}
	}

	return reflect.ValueOf(vInt), nil
}

func decodeString(b *[]byte) (reflect.Value, error) {
	_, err := popByte(b)
	if err != nil {
		return reflect.Value{}, err
	}

	v, err := popBytesUntil(b, ':')
	if err != nil {
		return reflect.Value{}, err
	}

	length, err := strconv.ParseInt(string(v), 10, 64)
	if err != nil {
		return reflect.Value{}, err
	}

	v, err = popBytes(b, int(length))
	if err != nil {
		return reflect.Value{}, err
	}

	return reflect.ValueOf(string(v)), nil
}

func decodeBytes(b *[]byte) (reflect.Value, error) {
	v, err := popBytesUntil(b, ':')
	if err != nil {
		return reflect.Value{}, err
	}
	length, err := strconv.ParseInt(string(v), 10, 64)
	if err != nil {
		return reflect.Value{}, err
	}
	v, err = popBytes(b, int(length))
	if err != nil {
		return reflect.Value{}, err
	}
	return reflect.ValueOf(v), nil
}

func decodeList(b *[]byte) (reflect.Value, error) {
	_, err := popByte(b)
	if err != nil {
		return reflect.Value{}, err
	}

	var list []any
	for (*b)[0] != 'e' {
		val, err := DecodeValue(b)
		if err != nil {
			return reflect.Value{}, err
		}
		list = append(list, safeInterface(val))
	}

	_, err = popByte(b)
	if err != nil {
		return reflect.Value{}, err
	}

	return reflect.ValueOf(list), nil
}

func decodeDictionary(b *[]byte) (reflect.Value, error) {
	_, err := popByte(b)
	if err != nil {
		return reflect.Value{}, err
	}

	dict := bencodextype.NewDictionary()
	for (*b)[0] != 'e' {
		key, err := DecodeValue(b)
		if err != nil {
			return reflect.Value{}, err
		}

		val, err := DecodeValue(b)
		if err != nil {
			return reflect.Value{}, err
		}
		dict.Set(safeInterface(key), safeInterface(val))
	}

	_, err = popByte(b)
	if err != nil {
		return reflect.Value{}, err
	}

	return reflect.ValueOf(dict), nil
}
