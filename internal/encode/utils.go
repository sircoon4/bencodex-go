package encode

import (
	"bytes"
	"reflect"
	"sort"
)

type UnsupportedTypeError struct {
	Type reflect.Type
	Kind reflect.Kind
}

func (e *UnsupportedTypeError) Error() string {
	return "unsupported type: " + e.Type.String() + " " + e.Kind.String()
}

func isNil(val reflect.Value) bool {
	if !val.IsValid() {
		return true
	}
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return val.IsNil()
	}
	return false
}

func sortEncodedPropertySlice(val [][]byte) [][]byte {
	sort.Slice(val, func(i, j int) bool {
		if val[i][0] > val[j][0] {
			return false
		} else if val[i][0] < val[j][0] {
			return true
		} else {
			iValue := bytes.Split(val[i], []byte(":"))[1]
			jValue := bytes.Split(val[j], []byte(":"))[1]
			return bytes.Compare(iValue, jValue) < 0
		}
	})
	return val
}
