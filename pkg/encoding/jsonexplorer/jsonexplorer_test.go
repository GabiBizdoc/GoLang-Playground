package jsonexplorer_test

import (
	"encoding/json"
	"fmt"
	. "github.com/GabiBizdoc/golang-playground/pkg/encoding/jsonexplorer"
	"os"
	"reflect"
	"testing"
)

var data any

func TestMain(m *testing.M) {
	json.Unmarshal([]byte(`[0, 1,2,3, [4,5,6], "seven", "eight", {"nime": "nine", "ten": true, "n": 42}]`), &data)
	os.Exit(m.Run())
}

func TestJSONExplorer_UnknownValue(t *testing.T) {
	_data := NewJSONExplorer(data)

	value, _ := _data.At(7).Value()
	expected, _ := ValueOf[map[string]any](_data.At(7))

	if !reflect.DeepEqual(value, value) {
		t.Errorf("Expected %v but got %v", expected, value)
	}
}
func TestJSONExplorer_At(t *testing.T) {
	_data := NewJSONExplorer(data)
	expected := data.([]any)[7]

	// accessing data shouldn't change the _data
	_data.At(4)
	value, err := ValueOf[map[string]any](_data.At(7))
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(value, expected) {
		t.Errorf("Expected %v, but got %v", expected, value)
	}
}

func TestJSONExplorer_Field(t *testing.T) {
	const expected = true
	value, err := ValueOf[bool](NewJSONExplorer(data).At(7).Field("ten"))
	if err != nil {
		t.Error(err)
	}
	if value != expected {
		t.Errorf("Expected %v, but got %v", expected, value)
	}
}

func TestJSONExplorer_Traverse(t *testing.T) {
	const expected = 42
	value, err := ValueOf[int8](NewJSONExplorer(data).Traverse(7, "n"))
	if err != nil {
		t.Error(err)
	}
	if value != expected {
		t.Errorf("Expected %v, but got %v", expected, value)
	}
}

func TestJSONExplorer_InvalidPath(t *testing.T) {
	value, err := ValueOf[int](NewJSONExplorer(data).Traverse("invalid", "invalid"))
	if err == nil {
		t.Errorf("Expected an error for invalid path traversal, but got %v", value)
	}
}

func TestJSONExplorer_InvalidKeys(t *testing.T) {
	value, err := ValueOf[int](NewJSONExplorer(any(true)).Traverse("true"))
	fmt.Println(value, err)
	if err == nil {
		t.Errorf("Expected an error for invalid type, but got %v", value)
	}
}

func TestJSONExplorer_TraverseToKey(t *testing.T) {
	value, err := ValueOf[int](NewJSONExplorer(data).TraverseToKey("n"))
	if err != nil {
		t.Errorf("Error traversing to key 'n': %v", err)
	}
	if value != 42 {
		t.Errorf("Expected 42, but got %v ", value)
	}
}
