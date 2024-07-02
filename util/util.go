package util

import (
	"math/big"
	"reflect"

	"github.com/sircoon4/bencodex-go/bencodextype"
)

func BencodexValueEqual(decoded1 any, decoded2 any) bool {
	isEqual := true

	// If the decoded data is a dictionary type, compare the values of the result and decoded data.
	d2Dict, ok := decoded2.(*bencodextype.Dictionary)
	if ok {
		d1Dict, ok := decoded1.(*bencodextype.Dictionary)
		if !ok {
			isEqual = false
		}
		if d2Dict.Length() != d1Dict.Length() {
			isEqual = false
		}
		for _, key := range d2Dict.Keys() {
			if d1Dict.Get(key) == nil {
				isEqual = false
			} else {
				isEqual = BencodexValueEqual(d1Dict.Get(key), d2Dict.Get(key))
			}
		}
	} else {
		// If the decoded data is a big.Int type, compare the values of the result and decoded data.
		d2BigInt, ok := decoded2.(*big.Int)
		if ok {
			d1BigInt, ok := decoded1.(*big.Int)
			if !ok {
				isEqual = false
			}
			if d2BigInt.Cmp(d1BigInt) != 0 {
				isEqual = false
			}
		} else {
			// If the decoded data is not a dictionary or big.Int type, compare the values of the result and decoded data.
			rvd1 := reflect.ValueOf(decoded1)
			rvd2 := reflect.ValueOf(decoded2)
			if rvd1.Kind() == rvd2.Kind() {
				switch rvd2.Kind() {
				case reflect.Slice:
					if rvd1.Len() == rvd2.Len() {
						for i := 0; i < rvd2.Len(); i++ {
							isEqual = BencodexValueEqual(rvd1.Index(i).Interface(), rvd2.Index(i).Interface())
						}
					} else {
						isEqual = false
					}
				default:
					if !reflect.DeepEqual(decoded1, decoded2) {
						isEqual = false
					}
				}
			} else {
				isEqual = false
			}
		}
	}

	return isEqual
}
