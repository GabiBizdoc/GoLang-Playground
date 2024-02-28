package json

import (
	"encoding/json"
	"testing"
)

func BenchmarkLexer_NextToken(b *testing.B) {
	tokenCount := 0
	expected := 2507032

	body := readLocalFile("large-file.json")
	lexer := NewLexer(body)
	for lexer.NextToken().Kind != TokenKindEOF {
		tokenCount += 1
	}
	if tokenCount != expected {
		b.Fatalf("tokenCount: expected %d but got %d", expected, tokenCount)
	}
}

// BenchmarkJSONModule-8   	1000000000	         0.1880 ns/op
func BenchmarkJSONModule(b *testing.B) {
	body := readLocalFile("large-file.json")
	b.ResetTimer()

	var data any
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		panic(err)
	}
}

// BenchmarkParseJson-8   	1000000000	         0.2974 ns/op +58%
func BenchmarkParseJson(b *testing.B) {
	body := readLocalFile("large-file.json")
	b.ResetTimer()

	_, err := ParseJson(body)
	if err != nil {
		panic(err)
	}
}

func getFiles() []string {
	files := make([]string, 0)
	for i := 0; i < 50; i++ {
		files = append(files, "large-file.json")
		files = append(files, "small-file.json")
		files = append(files, "something.json")
	}
	return files
}

// BenchmarkParseJsonMany-8   	       1	15150813916 ns/op +45%
func BenchmarkParseJsonMany(b *testing.B) {
	for _, filename := range getFiles() {
		body := readLocalFile(filename)

		_, err := ParseJson(body)
		if err != nil {
			panic(err)
		}
	}
}

// BenchmarkJSONModuleMany-8   	       1	10475346333 ns/op
func BenchmarkJSONModuleMany(b *testing.B) {
	for _, filename := range getFiles() {
		body := readLocalFile(filename)

		var data any
		err := json.Unmarshal([]byte(body), &data)
		if err != nil {
			panic(err)
		}
	}
}
