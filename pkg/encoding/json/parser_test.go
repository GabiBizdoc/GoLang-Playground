package json

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestParseJson(t *testing.T) {
	type TestCase struct {
		In  string
		Out any
	}
	cases := []TestCase{
		{In: "true", Out: true},
		{In: "-3.e-2", Out: -3.e-2},
		{In: "false", Out: false},
		{In: "null", Out: nil},
		{In: "[]", Out: []any{}},
		{In: "[true]", Out: []any{true}},
		{In: "[false]", Out: []any{false}},
		{In: `["hello"]`, Out: []any{"hello"}},
		{In: `["hello", "world"]`, Out: []any{"hello", "world"}},
		{In: `["hello", "world", "!", true]`, Out: []any{"hello", "world", "!", true}},
		{In: "{}", Out: map[string]any{}},
		{In: `{"a":"b"}`, Out: map[string]any{"a": "b"}},
		{In: `[{"hello":"world"}]`, Out: []any{map[string]any{"hello": "world"}}},
		{In: `"hello\r\n world"`, Out: "hello\r\n world"},
	}
	for _, testCase := range cases {
		value, err := ParseJson(testCase.In)
		if err != nil {
			t.Fail()
			t.Logf("FAIL: input %s error: %s", testCase.In, err.Error())
		}

		if !isEqual(value, testCase.Out) {
			t.Logf("FAIL: input: %s expected %#v but got %#v", testCase.In, testCase.Out, value)
			t.Fail()
		}
	}
	//value, err := ParseJson(`{"a":"b"}"`)
	//t.Log(value, err)
}

func TestParseJson_WithLargeObjects(t *testing.T) {
	for _, filename := range []string{"small-file.json", "large-file.json"} {
		body := readLocalFile(filename)
		value, err := ParseJson(body)
		if err != nil {
			t.Error(err)
		}
		var expected any
		err = json.Unmarshal([]byte(body), &expected)
		if err != nil {
			panic(err)
		}

		if !isEqual(expected, value) {
			t.Errorf("FAIL: Failed to parse %s file", filename)
		}
	}

}

func isEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false
	}

	return reflect.DeepEqual(a, b)
}
