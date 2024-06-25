package bencodex

import (
	"encoding/base64"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

func writeToWriter(w io.Writer, b []byte) error {
	l := len(b)
	for t := 0; t < l; {
		n, err := w.Write(b[t:])
		if err != nil {
			return err
		}
		t += n
	}
	return nil
}

func readFromReader(r io.Reader, b *[]byte) error {
	n, err := r.Read(*b)
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	} else if n < 1024 {
		*b = (*b)[:n]
	} else {
		for {
			ab := make([]byte, 64)
			n, err = r.Read(ab)
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			*b = append(*b, ab[:n]...)
		}
	}
	return nil
}

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
		// String not found. We add it to return slice
		if !found {
			diff = append(diff, s1)
		}
	}

	return diff
}

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
		dict := make(map[any]any)
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
