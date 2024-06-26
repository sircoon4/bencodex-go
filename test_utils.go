package bencodex

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func parseJsonData(jsonData map[string]any) (any, error) {
	switch jsonData["type"] {
	case "null":
		return nil, nil
	case "boolean":
		if jsonData["value"] == nil {
			return nil, fmt.Errorf("invalid json data")
		}
		return jsonData["value"].(bool), nil
	case "integer":
		if jsonData["decimal"] == nil {
			return nil, fmt.Errorf("invalid json data")
		}
		data, error := strconv.ParseInt(jsonData["decimal"].(string), 10, 64)
		if error != nil {
			return nil, error
		}
		return int64(data), nil
	case "binary":
		if jsonData["base64"] == nil {
			return nil, fmt.Errorf("invalid json data")
		}
		data, err := base64.StdEncoding.DecodeString(jsonData["base64"].(string))
		if err != nil {
			return nil, err
		}
		return data, nil
	case "text":
		if jsonData["value"] == nil {
			return nil, fmt.Errorf("invalid json data")
		}
		return jsonData["value"].(string), nil
	case "list":
		list := make([]any, 0)
		if jsonData["values"] == nil {
			return nil, fmt.Errorf("invalid json data")
		}
		for _, preItem := range jsonData["values"].([]any) {
			item := preItem.(map[string]any)
			val, err := parseJsonData(item)
			if err != nil {
				return nil, err
			}
			list = append(list, val)
		}
		return list, nil
	case "dictionary":
		dict := make(map[reflect.Value]any)
		if jsonData["pairs"] == nil {
			return nil, fmt.Errorf("invalid json data")
		}
		for _, prePair := range jsonData["pairs"].([]any) {
			pair := prePair.(map[string]any)
			if pair["key"] == nil || pair["value"] == nil {
				return nil, fmt.Errorf("invalid json data")
			}
			keyData, err := parseJsonData(pair["key"].(map[string]any))
			if err != nil {
				return nil, err
			}
			key := reflect.ValueOf(keyData)

			valData, err := parseJsonData(pair["value"].(map[string]any))
			if err != nil {
				return nil, err
			}
			dict[key] = valData
		}
		return dict, nil
	default:
		return nil, fmt.Errorf("invalid json data")
	}
}

func customizedAssertEqual(t *testing.T, result any, decoded any) {
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
		case reflect.Map:
			if rvr.Len() == rvd.Len() {
				d, ok := decoded.(map[reflect.Value]any)
				if !ok {
					t.Fatalf("decoded is not a valid map")
				}
				r, ok := result.(map[reflect.Value]any)
				if !ok {
					t.Fatalf("result is not a valid map")
				}
				for key, value := range d {
					if FindValue(r, key) == nil {
						t.Fatalf("result and decoded are not equal")
					} else {
						customizedAssertEqual(t, FindValue(r, key), value)
					}
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
