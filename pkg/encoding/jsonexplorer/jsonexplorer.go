package jsonexplorer

import (
	"errors"
	"fmt"
	"math"
	"reflect"
)

type Number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

type JSONValue interface {
	bool | float64 | string | []any | map[string]any | Number
}

type JSONExplorer struct {
	data any
	_err error
}

func NewJSONExplorer(data any) JSONExplorer {
	return JSONExplorer{data: data}
}

func ValueOf[T JSONValue](simpleJson JSONExplorer) (x T, err error) {
	if simpleJson._err != nil {
		return x, simpleJson._err
	}

	switch data := simpleJson.data.(type) {
	case T:
		return data, nil
	default:
		switch reflect.ValueOf(x).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:

			data2 := simpleJson.data.(float64)

			if data2 == math.Trunc(data2) {
				value := reflect.ValueOf(data2).Convert(reflect.TypeOf(x))
				return convertValue(value).(T), err
			}
			return x, fmt.Errorf("Cannot saftly convert type %T to %T\n", simpleJson.data, x)

		default:
			return x, fmt.Errorf("Cannot convert type %T to %T\n", simpleJson.data, x)
		}
	}
}

func convertValue(v reflect.Value) any {
	k := v.Kind()
	switch k {
	case reflect.Int:
		return int(v.Int())
	case reflect.Int8:
		return int8(v.Int())
	case reflect.Int16:
		return int16(v.Int())
	case reflect.Int32:
		return int32(v.Int())
	case reflect.Int64:
		return v.Int()
	case reflect.Uint:
		return uint(v.Uint())
	case reflect.Uint8:
		return uint8(v.Uint())
	case reflect.Uint16:
		return uint16(v.Uint())
	case reflect.Uint32:
		return uint32(v.Uint())
	case reflect.Uint64:
		return v.Uint()
	default:
		panic("invalid type " + v.Kind().String())
	}
}

func (j JSONExplorer) At(x int) JSONExplorer {
	if j._err != nil {
		return j
	}
	if x < 0 {
		j._err = fmt.Errorf("invalid index %d: index must be non-negative", x)
		return j
	}
	switch value := j.data.(type) {
	case []any:
		if x < len(value) {
			j.data = value[x]
		} else {
			j._err = fmt.Errorf("index out of range: %d (slice length: %d)", x, len(value))
		}
	default:
		j._err = fmt.Errorf("cannot index into type= %T", reflect.TypeOf(j.data).Kind())
	}
	return j
}

func (j JSONExplorer) Field(key string) JSONExplorer {
	if j._err != nil {
		return j
	}
	switch value := j.data.(type) {
	case map[string]any:
		if next, ok := value[key]; ok {
			j.data = next
		} else {
			j._err = fmt.Errorf("key %s not found in object", key)
		}
	default:
		j._err = fmt.Errorf("key: %s not found in type: %T", key, reflect.TypeOf(j.data).Kind().String())
	}
	return j
}

// Traverse traverses the JSONExplorer object and retrieves the value at the specified keys.
// It accepts only string and int keys.
func (j JSONExplorer) Traverse(keys ...any) JSONExplorer {
	for _, key := range keys {
		if j._err != nil {
			break
		}

		switch value := key.(type) {
		case int:
			j = j.At(value)
		case string:
			j = j.Field(value)
		default:
			j._err = fmt.Errorf("invalid type for key: %v", key)
		}
	}
	return j
}

func (j JSONExplorer) TraverseToKey(key string) JSONExplorer {
	if j._err != nil {
		return j
	}

	switch data := j.data.(type) {
	case []any:
		for _, item := range data {
			if found := NewJSONExplorer(item).TraverseToKey(key); found._err == nil {
				return found
			}
		}
	case map[string]any:
		for k, item := range data {
			explorer := NewJSONExplorer(item)
			if k == key {
				return explorer
			}
			if found := explorer.TraverseToKey(key); found._err == nil {
				return found
			}
		}
	}
	j._err = errors.New("TraverseToKey not found")
	return j
}

func (j JSONExplorer) Value() (any, error) {
	return j.data, j._err
}
