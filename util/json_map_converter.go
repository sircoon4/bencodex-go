package util

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strconv"

	"github.com/planetarium/bencodex-go/bencodextype"
)

func MarshalJsonMap(data any) ([]byte, error) {
	jsonMapData, err := convertBencodexDataToJsonMapData(data)
	if err != nil {
		return nil, err
	}

	out, err := json.MarshalIndent(jsonMapData, "", "  ")
	if err != nil {
		return nil, err
	}

	return out, nil
}

func convertBencodexDataToJsonMapData(data any) (map[string]any, error) {
	if data == nil {
		return map[string]any{"type": "null"}, nil
	}

	val := reflect.ValueOf(data)
	switch val.Kind() {
	case reflect.Bool:
		return map[string]any{"type": "boolean", "value": data.(bool)}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return map[string]any{"type": "integer", "decimal": strconv.FormatInt(val.Int(), 10)}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return map[string]any{"type": "integer", "decimal": strconv.FormatUint(val.Uint(), 10)}, nil
	case reflect.String:
		return map[string]any{"type": "text", "value": data.(string)}, nil
	case reflect.Slice, reflect.Array:
		// If the slice is a byte slice, encode it as a byte slice
		if val.Type().Elem().Kind() == reflect.Uint8 {
			return map[string]any{"type": "binary", "base64": base64.StdEncoding.EncodeToString(data.([]byte))}, nil
		}

		list := make([]map[string]any, 0)
		for _, preItem := range data.([]any) {
			item, err := convertBencodexDataToJsonMapData(preItem)
			if err != nil {
				return nil, err
			}
			list = append(list, item)
		}
		return map[string]any{"type": "list", "values": list}, nil
	case reflect.Map:
		dict := data.(map[string]any)
		pairs := make([]map[string]any, 0)
		for key, value := range dict {
			keyData, err := convertBencodexDataToJsonMapData(key)
			if err != nil {
				return nil, err
			}
			valData, err := convertBencodexDataToJsonMapData(value)
			if err != nil {
				return nil, err
			}
			pairs = append(pairs, map[string]any{"key": keyData, "value": valData})
		}
		return map[string]any{"type": "dictionary", "pairs": pairs}, nil
	case reflect.Pointer:
		_, ok := val.Interface().(*bencodextype.Dictionary)
		if ok {
			dict := val.Interface().(*bencodextype.Dictionary)
			pairs := make([]map[string]any, 0)
			for _, key := range dict.Keys() {
				keyData, err := convertBencodexDataToJsonMapData(key)
				if err != nil {
					return nil, err
				}
				valData, err := convertBencodexDataToJsonMapData(dict.Get(key))
				if err != nil {
					return nil, err
				}
				pairs = append(pairs, map[string]any{"key": keyData, "value": valData})
			}
			return map[string]any{"type": "dictionary", "pairs": pairs}, nil
		}
		_, ok = val.Interface().(*big.Int)
		if ok {
			return map[string]any{"type": "integer", "decimal": val.Interface().(*big.Int).String()}, nil
		}
		return nil, fmt.Errorf("ConvertToBencodexJsonMapData: unsupported type")
	default:
		return nil, fmt.Errorf("ConvertToBencodexJsonMapData: unsupported type")
	}
}

func UnmarshalJsonMap(jsonData []byte) (any, error) {
	var jsonMapData map[string]any

	err := json.Unmarshal(jsonData, &jsonMapData)
	if err != nil {
		return nil, err
	}

	return convertJsonMapDataToBencodexData(jsonMapData)
}

func convertJsonMapDataToBencodexData(jsonMapData map[string]any) (any, error) {
	switch jsonMapData["type"] {
	case "null":
		return nil, nil
	case "boolean":
		if jsonMapData["value"] == nil {
			return nil, fmt.Errorf("invalid map data")
		}
		return jsonMapData["value"].(bool), nil
	case "integer":
		if jsonMapData["decimal"] == nil {
			return nil, fmt.Errorf("invalid map data")
		}
		data, error := strconv.Atoi(jsonMapData["decimal"].(string))
		if error != nil {
			// If the value is too large to fit in an int, it is stored as a big.Int
			bigInt := new(big.Int)
			bigInt, ok := bigInt.SetString(jsonMapData["decimal"].(string), 10)
			if !ok {
				return nil, error
			} else {
				return bigInt, nil
			}
		}
		return data, nil
	case "binary":
		if jsonMapData["base64"] == nil {
			return nil, fmt.Errorf("invalid map data")
		}
		data, err := base64.StdEncoding.DecodeString(jsonMapData["base64"].(string))
		if err != nil {
			return nil, err
		}
		return data, nil
	case "text":
		if jsonMapData["value"] == nil {
			return nil, fmt.Errorf("invalid map data")
		}
		return jsonMapData["value"].(string), nil
	case "list":
		list := make([]any, 0)
		if jsonMapData["values"] == nil {
			return nil, fmt.Errorf("invalid map data")
		}
		for _, preItem := range jsonMapData["values"].([]any) {
			item := preItem.(map[string]any)
			val, err := convertJsonMapDataToBencodexData(item)
			if err != nil {
				return nil, err
			}
			list = append(list, val)
		}
		return list, nil
	case "dictionary":
		// if jsonData is dictionary type, return bencodex dictionary type
		dict := bencodextype.NewDictionary()
		if jsonMapData["pairs"] == nil {
			return nil, fmt.Errorf("invalid map data")
		}
		for _, prePair := range jsonMapData["pairs"].([]any) {
			pair := prePair.(map[string]any)
			if pair["key"] == nil || pair["value"] == nil {
				return nil, fmt.Errorf("invalid map data")
			}
			keyData, err := convertJsonMapDataToBencodexData(pair["key"].(map[string]any))
			if err != nil {
				return nil, err
			}
			key := keyData

			valData, err := convertJsonMapDataToBencodexData(pair["value"].(map[string]any))
			if err != nil {
				return nil, err
			}
			dict.Set(key, valData)
		}
		return dict, nil
	default:
		return nil, fmt.Errorf("invalid map data")
	}
}
