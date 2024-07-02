package bencodex

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/sircoon4/bencodex-go/bencodextype"
)

// remove slice2 from slice1
func difference(slice1 []string, slice2 []string) []string {
	var diff []string

	for _, s1 := range slice1 {
		found := false
		for _, s2 := range slice2 {
			if s1 == s2 {
				found = true
				break
			}
		}
		// if string is not foundec, add it to return slice
		if !found {
			diff = append(diff, s1)
		}
	}

	return diff
}

// customizedAssertEqual is a function that compares the real values of the result and decoded data.
func customizedAssertEqual(t *testing.T, result any, decoded any) {
	// If the decoded data is a dictionary type, compare the values of the result and decoded data.
	dDict, ok := decoded.(*bencodextype.Dictionary)
	if ok {
		rDict, ok := result.(*bencodextype.Dictionary)
		if !ok {
			t.Fatalf("result and decoded are not equal")
		}
		if dDict.Length() != rDict.Length() {
			t.Fatalf("result and decoded are not equal")
		}
		for _, key := range dDict.Keys() {
			if rDict.Get(key) == nil {
				t.Fatalf("result and decoded are not equal")
			} else {
				customizedAssertEqual(t, rDict.Get(key), dDict.Get(key))
			}
		}
	} else {
		// If the decoded data is a big.Int type, compare the values of the result and decoded data.
		dBigInt, ok := decoded.(*big.Int)
		if ok {
			rBigInt, ok := result.(*big.Int)
			if !ok {
				t.Fatalf("result and decoded are not equal")
			}
			if dBigInt.Cmp(rBigInt) != 0 {
				t.Fatalf("result and decoded are not equal")
			}
		} else {
			// If the decoded data is not a dictionary or big.Int type, compare the values of the result and decoded data.
			rvr := reflect.ValueOf(result)
			rvd := reflect.ValueOf(decoded)
			if rvr.Kind() == rvd.Kind() {
				switch rvd.Kind() {
				case reflect.Slice:
					if rvr.Len() == rvd.Len() {
						for i := 0; i < rvd.Len(); i++ {
							customizedAssertEqual(t, rvr.Index(i).Interface(), rvd.Index(i).Interface())
						}
					} else {
						t.Fatalf("result and decoded are not equal")
					}
				default:
					if !reflect.DeepEqual(result, decoded) {
						t.Fatalf("result and decoded are not equal")
					}
				}
			} else {
				t.Fatalf("result and decoded are not equal")
			}
		}
	}
}
