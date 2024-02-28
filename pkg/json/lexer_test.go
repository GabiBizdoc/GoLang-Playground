package json

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLexer_NextToken(t *testing.T) {
	const input = `  123  "456" null   true false "false" "true"  "trf`
	tks := strings.Fields(input)
	lexer := NewLexer(input)
	cmpTokens(t, lexer, tks)
}

func TestLexer_NextTokenUsingJson(t *testing.T) {
	const input = `{"glossary":{"title":"example glossary","GlossDiv":{"title":"S","GlossList":{"GlossEntry":
{"ID":"SGML","SortAs":"SGML","GlossTerm":"Standard Generalized Markup Language","Acronym":"SGML","Abbrev":
"ISO 8879:1986","GlossDef":{"para":"A meta-markup language, used to create markup languages such as DocBook.",
"GlossSeeAlso":["GML","XML"]},"GlossSee":"markup"}}}}}
`
	tokens := []string{"{", "\"glossary\"", ":", "{", "\"title\"", ":", "\"example glossary\"", ",",
		"\"GlossDiv\"", ":", "{", "\"title\"", ":", "\"S\"", ",", "\"GlossList\"", ":",
		"{", "\"GlossEntry\"", ":", "{", "\"ID\"", ":", "\"SGML\"", ",", "\"SortAs\"", ":",
		"\"SGML\"", ",", "\"GlossTerm\"", ":", "\"Standard Generalized Markup Language\"", ",",
		"\"Acronym\"", ":", "\"SGML\"", ",", "\"Abbrev\"", ":",
		"\"ISO 8879:1986\"", ",", "\"GlossDef\"", ":", "{", "\"para\"", ":",
		"\"A meta-markup language, used to create markup languages such as DocBook.\"", ",",
		"\"GlossSeeAlso\"", ":", "[", "\"GML\"", ",", "\"XML\"", "]", "}", ",",
		"\"GlossSee\"", ":", "\"markup\"", "}", "}", "}", "}", "}", "EOF"}

	lexer := NewLexer(input)

	cmpTokens(t, lexer, tokens)
	//previewTokens(NewLexer(input), true)
}

func TestLexer_NextTokenSpecialCases(t *testing.T) {
	const input = `  {"kk":-0.e-23} {"kind": [true, false, 3.4e2, -1, {"key": "value"}]}`
	lexer := NewLexer(input)
	tokens := []string{"{", "\"kk\"", ":", "-0.e-23", "}",
		"{", "\"kind\"", ":", "[", "true", ",", "false", ",", "3.4e2", ",", "-1", ",",
		"{", "\"key\"", ":", "\"value\"", "}", "]", "}", "EOF"}
	cmpTokens(t, lexer, tokens)
	//previewTokens(NewLexer(input), false)
}

func TestLexer_NextTokenSpecialCases2(t *testing.T) {
	const input = `":,"`
	lexer := NewLexer(input)
	tokens := []string{"\":,\"", "EOF"}
	cmpTokens(t, lexer, tokens)
	//previewTokens(NewLexer(input), false)
}

func TestLexer_NextTokenEscape(t *testing.T) {
	input := `"hi \"Steve\" how are you?"`
	lexer := NewLexer(input)
	tokens := []string{input, "EOF"}
	cmpTokens(t, lexer, tokens)
	//previewTokens(NewLexer(input), false)
}

func TestLexer_NextTokenBoolean(t *testing.T) {
	const input = `true`
	lexer := NewLexer(input)
	tokens := []string{"true", "EOF"}
	cmpTokens(t, lexer, tokens)
	//previewTokens(NewLexer(input), false)
}
func TestLexer_NextTokenNumber(t *testing.T) {
	const input = `123`
	lexer := NewLexer(input)
	tokens := []string{"123", "EOF"}
	cmpTokens(t, lexer, tokens)
	//previewTokens(NewLexer(input), false)
}

func cmpTokens(t *testing.T, lexer *Lexer, tokens []string) {
	for _, expected := range tokens {
		token := lexer.NextToken()
		if token.Value != expected {
			t.Errorf("FAIL: expected: `%s` but got `%s`. type: %s", expected, token.Value, token.Kind.toString())
		} else {
			//t.Logf("OK__: `%s` type: %s", token.Value, token.Kind.toString())
		}
	}

	next := lexer._getNextToken()
	if next.Kind != TokenKindEOF {
		t.Errorf("too make tokens. expected %s but got %s", TokenKindEOF.toString(), next.Kind.toString())
	}
}

func readLocalFile(file string) string {
	largeJsonFile, err := filepath.Abs(filepath.Join("tests", file))
	if err != nil {
		panic(err)
	}
	body, err := os.ReadFile(largeJsonFile)
	if err != nil {
		panic(err)
	}
	return string(body)
}

func previewTokens(lexer *Lexer, asVariable bool) {
	var sb strings.Builder
	if asVariable {
		sb.WriteString("tokens := []string{")
	}
	for {
		token := lexer.NextToken()
		if asVariable {
			sb.WriteString("\"")
		}

		if asVariable && token.Kind == TokenKindString {
			sb.WriteString("\\\"")
			sb.WriteString(unquoteString(token.Value))
			sb.WriteString("\\\"")
		} else {
			sb.WriteString(token.Value)
		}
		if asVariable {
			sb.WriteString("\"")
			sb.WriteString(",")
		}
		sb.WriteRune(' ')

		if token.Kind == TokenKindEOF {
			break
		}
	}
	if asVariable {
		sb.WriteString("}")
	}
	log.Println(sb.String())
}
