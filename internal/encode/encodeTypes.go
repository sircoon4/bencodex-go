package encode

import (
	"reflect"
	"strconv"

	"github.com/planetarium/bencodex-go/bencodextype"
)

func encodeNil() []byte {
	return []byte("n")
}

func encodeBool(val reflect.Value) []byte {
	if val.Bool() {
		return []byte("t")
	}
	return []byte("f")
}

func encodeInteger(val reflect.Value) []byte {
	return []byte("i" + strconv.FormatInt(val.Int(), 10) + "e")
}

func encodeUnsignedInteger(val reflect.Value) []byte {
	return []byte("i" + strconv.FormatUint(val.Uint(), 10) + "e")
}

func encodeString(val reflect.Value) []byte {
	return []byte("u" + strconv.Itoa(val.Len()) + ":" + val.String())
}

func encodeBytes(val reflect.Value) []byte {
	if !val.CanAddr() {
		tmpVal := reflect.New(val.Type())
		tmpVal.Elem().Set(val)
		val = tmpVal.Elem()
	}
	return []byte(strconv.Itoa(val.Len()) + ":" + string(val.Bytes()))
}

func encodeList(val reflect.Value) ([]byte, error) {
	buf := []byte("l")
	for i := 0; i < val.Len(); i++ {
		encoded, err := EncodeValue(val.Index(i))
		if err != nil {
			return nil, err
		}
		buf = append(buf, encoded...)
	}
	buf = append(buf, 'e')
	return buf, nil
}

func encodeMap(val reflect.Value) ([]byte, error) {
	encodedPropertySlice := make([][]byte, 0)
	keys := val.MapKeys()
	for _, key := range keys {
		encodedKey, err := EncodeValue(key)
		if err != nil {
			return nil, err
		}
		encodedVal, err := EncodeValue(val.MapIndex(key))
		if err != nil {
			return nil, err
		}

		encodedProperty := make([]byte, 0)
		encodedProperty = append(encodedProperty, encodedKey...)
		encodedProperty = append(encodedProperty, encodedVal...)

		encodedPropertySlice = append(encodedPropertySlice, encodedProperty)
	}

	buf := []byte("d")
	for _, encodedProperty := range sortEncodedPropertySlice(encodedPropertySlice) {
		buf = append(buf, encodedProperty...)
	}
	buf = append(buf, 'e')
	return buf, nil
}

func encodeDictionary(val reflect.Value) ([]byte, error) {
	dict := val.Interface().(*bencodextype.Dictionary)

	encodedPropertySlice := make([][]byte, 0)
	keys := dict.Keys()
	for _, key := range keys {
		encodedKey, err := EncodeValue(reflect.ValueOf(key))
		if err != nil {
			return nil, err
		}
		encodedVal, err := EncodeValue(reflect.ValueOf(dict.Get(key)))
		if err != nil {
			return nil, err
		}

		encodedProperty := make([]byte, 0)
		encodedProperty = append(encodedProperty, encodedKey...)
		encodedProperty = append(encodedProperty, encodedVal...)

		encodedPropertySlice = append(encodedPropertySlice, encodedProperty)
	}

	buf := []byte("d")
	for _, encodedProperty := range sortEncodedPropertySlice(encodedPropertySlice) {
		buf = append(buf, encodedProperty...)
	}
	buf = append(buf, 'e')
	return buf, nil
}
