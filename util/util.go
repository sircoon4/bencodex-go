package util

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strconv"

	"github.com/planetarium/bencodex-go/bencodextype"
	"gopkg.in/yaml.v3"
)

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

func UnmarshalYaml(yamlData []byte) (any, error) {
	var yamlNode yaml.Node

	err := yaml.Unmarshal(yamlData, &yamlNode)
	if err != nil {
		return nil, err
	}

	if yamlNode.Kind == yaml.DocumentNode {
		if len(yamlNode.Content) > 0 {
			return convertYamlNodeToBencodexData(yamlNode.Content[0])
		} else {
			return nil, nil
		}
	}

	return convertYamlNodeToBencodexData(&yamlNode)
}

func convertYamlNodeToBencodexData(yamlNode *yaml.Node) (any, error) {
	switch yamlNode.Kind {
	case yaml.ScalarNode:
		switch yamlNode.Tag {
		case "!!null":
			return nil, nil
		case "!!bool":
			return yamlNode.Value == "true", nil
		case "!!int":
			data, err := strconv.Atoi(yamlNode.Value)
			if err != nil {
				// If the value is too large to fit in an int, it is stored as a big.Int
				bigInt := new(big.Int)
				bigInt, ok := bigInt.SetString(yamlNode.Value, 10)
				if !ok {
					return nil, err
				} else {
					return bigInt, nil
				}
			}
			return data, nil
		case "!!str":
			return yamlNode.Value, nil
		case "!!binary":
			data, err := base64.StdEncoding.DecodeString(yamlNode.Value)
			if err != nil {
				return nil, err
			}
			return data, nil
		default:
			return nil, fmt.Errorf("invalid yaml data")
		}
	case yaml.SequenceNode:
		list := make([]any, 0)
		for _, preItem := range yamlNode.Content {
			item, err := convertYamlNodeToBencodexData(preItem)
			if err != nil {
				return nil, err
			}
			list = append(list, item)
		}
		return list, nil
	case yaml.MappingNode:
		dict := bencodextype.NewDictionary()
		for i := 0; i < len(yamlNode.Content); i += 2 {
			key, err := convertYamlNodeToBencodexData(yamlNode.Content[i])
			if err != nil {
				return nil, err
			}
			val, err := convertYamlNodeToBencodexData(yamlNode.Content[i+1])
			if err != nil {
				return nil, err
			}
			dict.Set(key, val)
		}
		return dict, nil
	default:
		return nil, fmt.Errorf("invalid yaml data")
	}
}

func MarshalYaml(data any) ([]byte, error) {
	yamlNode, err := convertBencodexDataToYamlNode(data)
	if err != nil {
		return nil, err
	}

	rootNode := &yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{yamlNode}}

	out, err := yaml.Marshal(rootNode)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func convertBencodexDataToYamlNode(data any) (*yaml.Node, error) {
	if data == nil {
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null", Value: "null"}, nil
	}

	val := reflect.ValueOf(data)
	switch val.Kind() {
	case reflect.Bool:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: fmt.Sprintf("%t", data.(bool))}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: fmt.Sprintf("%d", val.Int())}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: fmt.Sprintf("%d", val.Uint())}, nil
	case reflect.String:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: data.(string)}, nil
	case reflect.Slice, reflect.Array:
		// If the slice is a byte slice, encode it as a byte slice
		if val.Type().Elem().Kind() == reflect.Uint8 {
			return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!binary", Value: base64.StdEncoding.EncodeToString(data.([]byte))}, nil
		}

		list := make([]*yaml.Node, 0)
		for _, preItem := range data.([]any) {
			item, err := convertBencodexDataToYamlNode(preItem)
			if err != nil {
				return nil, err
			}
			list = append(list, item)
		}
		return &yaml.Node{Kind: yaml.SequenceNode, Content: list}, nil
	case reflect.Map:
		dict := data.(map[string]any)
		content := make([]*yaml.Node, 0)
		for key, value := range dict {
			keyNode, err := convertBencodexDataToYamlNode(key)
			if err != nil {
				return nil, err
			}
			valNode, err := convertBencodexDataToYamlNode(value)
			if err != nil {
				return nil, err
			}
			content = append(content, keyNode, valNode)
		}
		return &yaml.Node{Kind: yaml.MappingNode, Content: content}, nil
	case reflect.Pointer:
		_, ok := val.Interface().(*bencodextype.Dictionary)
		if ok {
			dict := val.Interface().(*bencodextype.Dictionary)
			content := make([]*yaml.Node, 0)
			for _, key := range dict.Keys() {
				keyNode, err := convertBencodexDataToYamlNode(key)
				if err != nil {
					return nil, err
				}
				valNode, err := convertBencodexDataToYamlNode(dict.Get(key))
				if err != nil {
					return nil, err
				}
				content = append(content, keyNode, valNode)
			}
			return &yaml.Node{Kind: yaml.MappingNode, Content: content}, nil
		}
		_, ok = val.Interface().(*big.Int)
		if ok {
			return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: val.Interface().(*big.Int).String()}, nil
		}
		return nil, fmt.Errorf("convertToBencodexYamlData: unsupported type")
	default:
		return nil, fmt.Errorf("convertToBencodexYamlData: unsupported type")
	}
}

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
