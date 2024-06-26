package decode

import (
	"bytes"
	"fmt"
	"reflect"
)

// get the first byte of the byte array and remove it
func popByte(b *[]byte) (byte, error) {
	if len(*b) <= 0 {
		return 0, fmt.Errorf("data is not compatible with bencodex")
	}

	var v byte = (*b)[0]
	*b = (*b)[1:]

	return v, nil
}

// get the first c bytes of the byte array and remove them
func popBytes(b *[]byte, c int) ([]byte, error) {
	if len(*b) < c {
		return nil, fmt.Errorf("data is not compatible with bencodex")
	}

	var v []byte = (*b)[:c]
	*b = (*b)[c:]

	return v, nil
}

// get the bytes until the first c byte of the byte array and remove them
func popBytesUntil(b *[]byte, c byte) ([]byte, error) {
	until := bytes.IndexByte(*b, c)

	if until == -1 {
		return nil, fmt.Errorf("data is not compatible with bencodex")
	}

	var v []byte = (*b)[:until]
	*b = (*b)[until+1:]

	return v, nil
}

// get nil value if v is invalid type of reflect.Value or get the interface of v
func safeInterface(v reflect.Value) any {
	if !v.IsValid() {
		return nil
	}
	return v.Interface()
}
