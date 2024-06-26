package bencodextype

import "fmt"

type Dictionary map[string]any

func (d *Dictionary) Set(key any, value any) {
	switch k := key.(type) {
	case string:
		key = fmt.Sprintf("s:%s", k)
		(*d)[key.(string)] = value
	case []byte:
		key = fmt.Sprintf("b:%s", k)
		(*d)[key.(string)] = value
	default:
		panic("key must be a string or a byte array")
	}
}

func (d *Dictionary) Get(key any) any {
	switch k := key.(type) {
	case string:
		key = fmt.Sprintf("s:%s", k)
		return (*d)[key.(string)]
	case []byte:
		key = fmt.Sprintf("b:%s", k)
		return (*d)[key.(string)]
	default:
		panic("key must be a string or a byte slice")
	}
}

func (d *Dictionary) Delete(key any) {
	switch k := key.(type) {
	case string:
		key = fmt.Sprintf("s:%s", k)
		delete(*d, key.(string))
	case []byte:
		key = fmt.Sprintf("b:%s", k)
		delete(*d, key.(string))
	default:
		panic("key must be a string or a byte slice")
	}
}

func (d *Dictionary) Contains(key any) bool {
	switch k := key.(type) {
	case string:
		key = fmt.Sprintf("s:%s", k)
		_, ok := (*d)[key.(string)]
		return ok
	case []byte:
		key = fmt.Sprintf("b:%s", k)
		_, ok := (*d)[key.(string)]
		return ok
	default:
		panic("key must be a string or a byte slice")
	}
}

func (d *Dictionary) Keys() []any {
	keys := make([]any, 0)
	for key := range *d {
		switch key[:2] {
		case "s:":
			keys = append(keys, key[2:])
		case "b:":
			keys = append(keys, []byte(key[2:]))
		default:
			panic("dictionary contains invalid key")
		}
	}
	return keys
}

func (d *Dictionary) Values() []any {
	values := make([]any, 0)
	for _, value := range *d {
		values = append(values, value)
	}
	return values
}

func (d *Dictionary) Length() int {
	return len(*d)
}

// Dictionary can be converted to a map if all keys are strings
func (d *Dictionary) CanConvertToMap() bool {
	for key := range *d {
		switch key[:2] {
		case "s:":
			_, ok := (*d)[key].(*Dictionary)
			if ok {
				if !(*(*d)[key].(*Dictionary)).CanConvertToMap() {
					return false
				}
			}
			continue
		default:
			return false
		}
	}
	return true
}

// Use CanConvertToMap() to check if the dictionary can be converted to a map
func (d *Dictionary) ConvertToMap() map[string]any {
	if !d.CanConvertToMap() {
		panic("dictionary cannot be converted to map")
	}
	m := make(map[string]any)
	for key, value := range *d {
		_, ok := value.(*Dictionary)
		if ok {
			value = value.(*Dictionary).ConvertToMap()
		}
		m[key[2:]] = value
	}
	return m
}

func NewDictionary() *Dictionary {
	d := make(Dictionary)
	return &d
}

func NewDictionaryFromMap(m map[string]any) *Dictionary {
	d := make(Dictionary)
	for key, value := range m {
		d.Set(key, value)
	}
	return &d
}
