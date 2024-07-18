package encode

import (
	"bytes"
	"reflect"
	"sort"
	"strconv"
)

type UnsupportedTypeError struct {
	Type reflect.Type
	Kind reflect.Kind
}

func (e *UnsupportedTypeError) Error() string {
	return "unsupported type: " + e.Type.String() + " " + e.Kind.String()
}

// check if the reflect.Value wraps nil value
func isNil(val reflect.Value) bool {
	return !val.IsValid()
}

// Unicode strings do not appear earlier than byte strings.
// Byte strings are sorted as raw strings, not alphanumerics.
// Unicode strings are sorted as their UTF-8 byte representations, not any collation order or chart order listed by Unicode.
// For example, b (62) should be followed by á (C3 A1), because the byte 62 is less than the byte C3.
func sortEncodedPropertySlice(val [][]byte) [][]byte {
	sort.Slice(val, func(i, j int) bool {
		if val[i][0] > val[j][0] {
			return false
		} else if val[i][0] < val[j][0] {
			return true
		} else {
			iValue := getKeyFromEncodedMapPairData(val[i])
			jValue := getKeyFromEncodedMapPairData(val[j])
			cr := bytes.Compare(iValue, jValue)
			if cr == 0 {
				return len(iValue) < len(jValue)
			}
			return cr < 0
		}
	})
	return val
}

func getKeyFromEncodedMapPairData(data []byte) []byte {
	indicator := bytes.Split(data, []byte(":"))[0]
	preKey := bytes.Split(data, []byte(":"))[1]
	var dlen int
	var err error
	switch indicator[0] {
	case 'u':
		dlen, err = strconv.Atoi(string(indicator[1:]))
		if err != nil {
			panic(err)
		}
	default:
		dlen, err = strconv.Atoi(string(indicator))
		if err != nil {
			panic(err)
		}
	}
	return preKey[:dlen]
}
