package json

import (
	"fmt"
	"unicode"
)

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}
func isWhitespace(r rune) bool {
	switch r {
	case '\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0:
		return true
	}
	return false
}
func isNumberStart(r rune) bool {
	return r == '-' || isDigit(r)
}
func isObjectStart(r rune) bool {
	return r == '{'
}
func isObjectEnd(r rune) bool {
	return r == '}'
}
func isArrayStart(r rune) bool {
	return r == '['
}
func isArrayEnd(r rune) bool {
	return r == ']'
}
func isBoolStart(r rune) bool {
	return r == 'f' || r == 't'
}
func isNullStart(r rune) bool {
	return r == 'n'
}
func isEscape(r rune) bool {
	return r == '\\'
}
func isQuote(r rune) bool {
	return r == '"'
}
func isColon(r rune) bool {
	return r == ':'
}
func isComma(r rune) bool {
	return r == ','
}

type Lexer struct {
	input        []rune
	pos          int
	c            rune
	done         bool
	currentToken Token
}

func NewLexer(s string) *Lexer {
	l := &Lexer{input: []rune(s)}
	l.advance()
	return l
}

func (l *Lexer) NextToken() Token {
	l.currentToken = l._getNextToken()
	l.advance()
	return l.currentToken
}

func (l *Lexer) _getNextToken() Token {
	l.skipWhiteSpace()
	switch {
	case l.done:
		return NewToken(TokenKindEOF, "EOF")
	case isObjectStart(l.c):
		return NewTokenFromRune(TokenKindBraceOpen, l.c)
	case isObjectEnd(l.c):
		return NewTokenFromRune(TokenKindBraceClose, l.c)
	case isArrayStart(l.c):
		return NewTokenFromRune(TokenKindBracketOpen, l.c)
	case isArrayEnd(l.c):
		return NewTokenFromRune(TokenKindBracketClose, l.c)
	case isQuote(l.c):
		return NewToken(TokenKindString, l.readString())
	case isColon(l.c):
		return NewTokenFromRune(TokenKindColon, l.c)
	case isComma(l.c):
		return NewTokenFromRune(TokenKindComma, l.c)
	case isBoolStart(l.c):
		return NewToken(TokenKindBoolean, l.readWord())
	case isNullStart(l.c):
		return NewToken(TokenKindNull, l.readWord())
	case isNumberStart(l.c):
		return NewToken(TokenKindNumber, l.readNumber())
	default:
		return NewToken(TokenKindInvalid, fmt.Sprintf("Invalid token `%c` at position %d", l.c, l.pos))
	}
}

func (l *Lexer) skipWhiteSpace() {
	for !l.done && unicode.IsSpace(l.c) {
		l.advance()
	}
}
func (l *Lexer) advance() {
	if l.pos >= len(l.input) {
		l.done = true
		return
	}
	l.c = l.input[l.pos]

	l.pos += 1
}
func (l *Lexer) prev() {
	if l.pos > 0 {
		l.pos -= 1
	}
}
func (l *Lexer) readString() string {
	var position = l.pos
	l.advance()

	escaped := false
	for !l.done && !(isQuote(l.c) && !escaped) {
		if isEscape(l.c) {
			escaped = !escaped
		} else {
			escaped = false
		}
		l.advance()
	}

	return string(l.input[position-1 : l.pos])
}

func (l *Lexer) readWord() string {
	var position = l.pos

WORD:
	for !l.done && !isWhitespace(l.c) {
		switch l.c {
		case ',', '{', '}', '[', ']':
			break WORD
		}
		l.advance()
	}
	if !l.done {
		l.prev()
	}
	return string(l.input[position-1 : l.pos])
}

func (l *Lexer) readNumber() string {
	var position = l.pos

	for !l.done && (isDigit(l.c) || l.c == '.' || l.c == 'E' || l.c == 'e' || l.c == '-') {
		l.advance()
	}
	if !l.done {
		l.prev()
	}
	return string(l.input[position-1 : l.pos])
}
