package util

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"reflect"
	"strconv"

	"github.com/sircoon4/bencodex-go/bencodextype"
	"gopkg.in/yaml.v3"
)

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
