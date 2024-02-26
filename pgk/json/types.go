package json

import "strings"

type TokenKind int8

const (
	TokenKindInvalid TokenKind = iota
	TokenKindNull
	TokenKindBoolean
	TokenKindNumber
	TokenKindString
	TokenKindBraceOpen
	TokenKindBraceClose
	TokenKindBracketOpen
	TokenKindBracketClose
	TokenKindColon
	TokenKindComma
	TokenKindEOF
)

func (t TokenKind) toString() string {
	tokens := strings.TrimSpace(`
	TokenKindInvalid
	TokenKindNull
	TokenKindBoolean
	TokenKindNumber
	TokenKindString
	TokenKindBraceOpen
	TokenKindBraceClose
	TokenKindBracketOpen
	TokenKindBracketClose
	TokenKindColon
	TokenKindComma
	TokenKindEOF
`)
	return strings.TrimSpace(strings.Split(tokens, "\n")[t])
}

type Token struct {
	Kind  TokenKind
	Value string
}

func NewToken(kind TokenKind, value string) Token {
	return Token{Kind: kind, Value: value}
}

func NewTokenFromRune(kind TokenKind, r rune) Token {
	return Token{Kind: kind, Value: string(r)}
}
