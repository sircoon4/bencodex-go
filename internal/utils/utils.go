package utils

import "reflect"

func FindValue(m map[reflect.Value]any, k any) any {
	var kr reflect.Value
	kr, ok := k.(reflect.Value)
	if ok {
		k = kr.Interface()
	} else {
		kr = reflect.ValueOf(k)
	}
	switch kr.Kind() {
	case reflect.String, reflect.Slice:
		for key, value := range m {
			if reflect.DeepEqual(key.Interface(), k) {
				return value
			}
		}
	case reflect.Array:
		if kr.Type().Elem().Kind() == reflect.Uint8 {
			ks := reflect.MakeSlice(reflect.SliceOf(kr.Type().Elem()), kr.Len(), kr.Len())
			for i := 0; i < kr.Len(); i++ {
				ks.Index(i).Set(kr.Index(i))
			}
			return FindValue(m, ks.Interface())
		}
	default:
		panic("key must be a string, a byte array, or a slice of byte")
	}

	return nil
}
