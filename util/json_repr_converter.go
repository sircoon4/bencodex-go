package util

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/sircoon4/bencodex-go/bencodextype"
)

func MarshalJsonRepr(data any) ([]byte, error) {
	jsonReprData, err := convertBencodexDataToJsonReprData(data)
	if err != nil {
		return nil, err
	}

	out, err := json.MarshalIndent(jsonReprData, "", "  ")
	if err != nil {
		return nil, err
	}

	return out, nil
}

const base64Threshold = 256

func convertBencodexDataToJsonReprData(data any) (any, error) {
	if data == nil {
		return nil, nil
	}

	val := reflect.ValueOf(data)
	switch val.Kind() {
	case reflect.Bool:
		return data.(bool), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(val.Uint(), 10), nil
	case reflect.String:
		return fmt.Sprintf("\ufeff%s", data.(string)), nil
	case reflect.Slice, reflect.Array:
		// If the slice is a byte slice, encode it as a byte slice
		if val.Type().Elem().Kind() == reflect.Uint8 {
			if len(data.([]byte)) < base64Threshold {
				return fmt.Sprintf("0x%x", data.([]byte)), nil
			} else {
				return fmt.Sprintf("b64:%s", base64.StdEncoding.EncodeToString(data.([]byte))), nil
			}
		}

		list := make([]any, 0)
		for _, preItem := range data.([]any) {
			item, err := convertBencodexDataToJsonReprData(preItem)
			if err != nil {
				return nil, err
			}
			list = append(list, item)
		}
		return list, nil
	case reflect.Map:
		mapData := data.(map[string]any)
		for key, value := range mapData {
			keyData, err := convertBencodexDataToJsonReprData(key)
			if err != nil {
				return nil, err
			}
			valData, err := convertBencodexDataToJsonReprData(value)
			if err != nil {
				return nil, err
			}
			mapData[keyData.(string)] = valData
		}
		return mapData, nil
	case reflect.Pointer:
		_, ok := val.Interface().(*bencodextype.Dictionary)
		if ok {
			dict := val.Interface().(*bencodextype.Dictionary)
			mapData := make(map[string]any)
			for _, key := range dict.Keys() {
				keyData, err := convertBencodexDataToJsonReprData(key)
				if err != nil {
					return nil, err
				}
				valData, err := convertBencodexDataToJsonReprData(dict.Get(key))
				if err != nil {
					return nil, err
				}
				mapData[keyData.(string)] = valData
			}
			return mapData, nil
		}
		_, ok = val.Interface().(*big.Int)
		if ok {
			return val.Interface().(*big.Int).String(), nil
		}
		return nil, fmt.Errorf("ConvertToBencodexJsonReprData: unsupported type")
	default:
		return nil, fmt.Errorf("ConvertToBencodexJsonReprData: unsupported type")
	}
}

func UnmarshalJsonRepr(data []byte) (any, error) {
	var jsonData any
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		return nil, err
	}

	return convertJsonReprDataToBencodexData(jsonData)
}

func convertJsonReprDataToBencodexData(jsonData any) (any, error) {
	switch jsonData.(type) {
	case nil:
		return nil, nil
	}

	switch val := reflect.ValueOf(jsonData); val.Kind() {
	case reflect.Bool:
		return val.Bool(), nil
	case reflect.String:
		return decodeJsonReprDataString(val.String())
	case reflect.Slice:
		list := make([]any, 0)
		for i := 0; i < val.Len(); i++ {
			data, err := convertJsonReprDataToBencodexData(val.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			list = append(list, data)
		}
		return list, nil
	case reflect.Map:
		dict := bencodextype.NewDictionary()
		for _, key := range val.MapKeys() {
			keyData, err := convertJsonReprDataToBencodexData(key.Interface())
			if err != nil {
				return nil, err
			}
			valData, err := convertJsonReprDataToBencodexData(val.MapIndex(key).Interface())
			if err != nil {
				return nil, err
			}
			dict.Set(keyData, valData)
		}
		return dict, nil
	default:
		return nil, fmt.Errorf("unsupported type")
	}
}

func decodeJsonReprDataString(s string) (any, error) {
	if len(s) == 0 {
		return nil, fmt.Errorf("empty string is not allowed")
	} else if len(s) > 0 && (s[0] == '-' ||
		(s[0] >= '1' && s[0] <= '9') ||
		(s[0] == '0' && (len(s) < 2 || s[1] != 'x'))) {
		val, err := strconv.Atoi(s)
		if err != nil {
			// If the value is too large to fit in an int, it is stored as a big.Int
			bigInt := new(big.Int)
			bigInt, ok := bigInt.SetString(s, 10)
			if !ok {
				return nil, err
			} else {
				return bigInt, nil
			}
		}
		return val, nil
	} else if strings.HasPrefix(s, "\ufeff") {
		runeVal := []rune(s)
		return string(runeVal[1:]), nil
	} else if strings.HasPrefix(s, "0x") {
		val, err := hex.DecodeString(s[2:])
		if err != nil {
			return nil, err
		}
		return val, nil
	} else if strings.HasPrefix(s, "b64:") {
		val, err := base64.StdEncoding.DecodeString(s[4:])
		if err != nil {
			return nil, err
		}
		return val, nil
	}

	return nil, fmt.Errorf("invalid string format")
}
