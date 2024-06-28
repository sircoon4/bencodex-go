package util

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strconv"

	"github.com/planetarium/bencodex-go/bencodextype"
	"gopkg.in/yaml.v2"
)

func ParseBencodexJasonMapData(jsonMapData map[string]any) (any, error) {
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
			val, err := ParseBencodexJasonMapData(item)
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
			keyData, err := ParseBencodexJasonMapData(pair["key"].(map[string]any))
			if err != nil {
				return nil, err
			}
			key := keyData

			valData, err := ParseBencodexJasonMapData(pair["value"].(map[string]any))
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

func ConvertToBencodexJsonMapData(data any) (map[string]any, error) {
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
			item, err := ConvertToBencodexJsonMapData(preItem)
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
			keyData, err := ConvertToBencodexJsonMapData(key)
			if err != nil {
				return nil, err
			}
			valData, err := ConvertToBencodexJsonMapData(value)
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
				keyData, err := ConvertToBencodexJsonMapData(key)
				if err != nil {
					return nil, err
				}
				valData, err := ConvertToBencodexJsonMapData(dict.Get(key))
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

func MarshalJson(data any) ([]byte, error) {
	jsonMapData, err := ConvertToBencodexJsonMapData(data)
	if err != nil {
		return nil, err
	}

	out, err := json.MarshalIndent(jsonMapData, "", "  ")
	if err != nil {
		return nil, err
	}

	return out, nil
}

func ConvertToBencodexYamlData(data any) (any, error) {
	if data == nil {
		return nil, nil
	}

	val := reflect.ValueOf(data)
	switch val.Kind() {
	case reflect.Bool:
		return data.(bool), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val.Uint(), nil
	case reflect.String:
		return data.(string), nil
	case reflect.Slice, reflect.Array:
		// If the slice is a byte slice, encode it as a byte slice
		if val.Type().Elem().Kind() == reflect.Uint8 {
			return fmt.Sprintf("!!binary \"%s\"", base64.StdEncoding.EncodeToString(data.([]byte))), nil
		}

		list := make([]any, 0)
		for _, preItem := range data.([]any) {
			item, err := ConvertToBencodexYamlData(preItem)
			if err != nil {
				return nil, err
			}
			list = append(list, item)
		}
		return list, nil
	case reflect.Map:
		mapData := make([]map[string]any, 0)
		for key, value := range data.(map[string]any) {
			keyData, err := ConvertToBencodexYamlData(key)
			if err != nil {
				return nil, err
			}
			valData, err := ConvertToBencodexYamlData(value)
			if err != nil {
				return nil, err
			}
			mapData = append(mapData, map[string]any{keyData.(string): valData})
		}
		return mapData, nil
	case reflect.Pointer:
		_, ok := val.Interface().(*bencodextype.Dictionary)
		if ok {
			dict := val.Interface().(*bencodextype.Dictionary)
			mapData := make([]map[string]any, 0)
			for _, key := range dict.Keys() {
				keyData, err := ConvertToBencodexYamlData(key)
				if err != nil {
					return nil, err
				}
				valData, err := ConvertToBencodexYamlData(dict.Get(key))
				if err != nil {
					return nil, err
				}
				mapData = append(mapData, map[string]any{keyData.(string): valData})
			}
			return mapData, nil
		}
		_, ok = val.Interface().(*big.Int)
		if ok {
			return fmt.Sprintf("\"BencodexBigIntInspector\"%s", val.Interface().(*big.Int).String()), nil
		}
		return nil, fmt.Errorf("ConverToBencodexYamlData: unsupported type")
	default:
		return nil, fmt.Errorf("ConverToBencodexYamlData: unsupported type")
	}
}

func MarshalYaml(data any) ([]byte, error) {
	yamlData, err := ConvertToBencodexYamlData(data)
	if err != nil {
		return nil, err
	}

	out, err := yaml.Marshal(yamlData)
	if err != nil {
		return nil, err
	}

	out = bytes.ReplaceAll(out, []byte("'"), []byte(""))
	out = bytes.ReplaceAll(out, []byte("\"BencodexBigIntInspector\""), []byte(""))

	return out, nil
}

func BencodexValueEqual(result any, decoded any) error {
	// If the decoded data is a dictionary type, compare the values of the result and decoded data.
	dDict, ok := decoded.(*bencodextype.Dictionary)
	if ok {
		rDict, ok := result.(*bencodextype.Dictionary)
		if !ok {
			return fmt.Errorf("result and decoded are not equal")
		}
		if dDict.Length() != rDict.Length() {
			return fmt.Errorf("result and decoded are not equal")
		}
		for _, key := range dDict.Keys() {
			if rDict.Get(key) == nil {
				return fmt.Errorf("result and decoded are not equal")
			} else {
				BencodexValueEqual(rDict.Get(key), dDict.Get(key))
			}
		}
	} else {
		// If the decoded data is a big.Int type, compare the values of the result and decoded data.
		dBigInt, ok := decoded.(*big.Int)
		if ok {
			rBigInt, ok := result.(*big.Int)
			if !ok {
				return fmt.Errorf("result and decoded are not equal")
			}
			if dBigInt.Cmp(rBigInt) != 0 {
				return fmt.Errorf("result and decoded are not equal")
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
							BencodexValueEqual(rvr.Index(i).Interface(), rvd.Index(i).Interface())
						}
					} else {
						return fmt.Errorf("result and decoded are not equal")
					}
				default:
					if !reflect.DeepEqual(result, decoded) {
						return fmt.Errorf("result and decoded are not equal")
					}
				}
			} else {
				return fmt.Errorf("result and decoded are not equal")
			}
		}
	}

	return nil
}
