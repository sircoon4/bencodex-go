package decode

import (
	"bytes"
	"fmt"
)

func popByte(b *[]byte) (byte, error) {
	if len(*b) <= 0 {
		return 0, fmt.Errorf("data is not compatible with bencodex")
	}

	var v byte = (*b)[0]
	*b = (*b)[1:]

	return v, nil
}

func popBytes(b *[]byte, c int) ([]byte, error) {
	if len(*b) < c {
		return nil, fmt.Errorf("data is not compatible with bencodex")
	}

	var v []byte = (*b)[:c]
	*b = (*b)[c:]

	return v, nil
}

func popBytesUntil(b *[]byte, c byte) ([]byte, error) {
	until := bytes.IndexByte(*b, c)

	if until == -1 {
		return nil, fmt.Errorf("data is not compatible with bencodex")
	}

	var v []byte = (*b)[:until]
	*b = (*b)[until+1:]

	return v, nil
}
