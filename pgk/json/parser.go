package json

import (
	"errors"
	"fmt"
	"strconv"
)

type jsonParser struct {
	lexer *Lexer
}

func newJsonParser(lexer *Lexer) *jsonParser {
	return &jsonParser{lexer: lexer}
}

func ParseJson(json string) (interface{}, error) {
	parser := newJsonParser(NewLexer(json))
	result, err := parser.parseValue()
	if err != nil {
		return nil, err
	}

	// check for any remaining tokens
	if parser.lexer.NextToken().Kind != TokenKindEOF {
		return nil, parser.invalidTokenError()
	}
	return result, nil
}

func (p *jsonParser) parseValue() (any, error) {
	for {
		token := p.lexer.NextToken()
		switch token.Kind {
		case TokenKindEOF:
			return nil, errors.New("EOF: end of file")
		case TokenKindInvalid:
			return nil, p.invalidTokenError()
		case TokenKindNull:
			return nil, nil
		case TokenKindBoolean:
			switch token.Value {
			case "true":
				return true, nil
			case "false":
				return false, nil
			default:
				return nil, p.invalidTokenError()
			}
		case TokenKindNumber:
			return strconv.ParseFloat(token.Value, 64)
		case TokenKindString:
			return unquoteString(token.Value), nil
		case TokenKindBraceOpen:
			return p.parseObject()
		case TokenKindBracketOpen:
			return p.parseArray()
		case TokenKindColon:
			return nil, p.invalidTokenError()
		case TokenKindComma:
			return nil, p.invalidTokenError()
		case TokenKindBraceClose:
			return nil, p.invalidTokenError()
		case TokenKindBracketClose:
			return nil, p.invalidTokenError()
		}
	}
}

func (p *jsonParser) parseArray() (any, error) {
	obj := make([]any, 0)
	for {
		value, err := p.parseValue()
		if err != nil {
			if p.lexer.currentToken.Kind == TokenKindBracketClose {
				return obj, nil
			}
		}
		obj = append(obj, value)

		switch p.lexer.NextToken().Kind {
		case TokenKindBracketClose:
			return obj, nil
		case TokenKindComma:
			continue
		default:
			return nil, p.invalidTokenError()
		}
	}

}
func (p *jsonParser) parseObject() (any, error) {
	obj := make(map[string]any)
	for {
		keyToken := p.lexer.NextToken()
		switch keyToken.Kind {
		case TokenKindBraceClose:
			return obj, nil
		case TokenKindString:
			key := unquoteString(keyToken.Value)
			if p.lexer.NextToken().Kind != TokenKindColon {
				return nil, p.invalidTokenError()
			}
			value, err := p.parseValue()
			if err != nil {
				return nil, err
			}
			obj[key] = value

			switch p.lexer.NextToken().Kind {
			case TokenKindBraceClose:
				return obj, nil
			case TokenKindComma:
				continue
			default:
				return nil, p.invalidTokenError()
			}
		default:
			return nil, p.invalidTokenError()
		}
	}
}

func (p *jsonParser) invalidTokenError() error {
	return fmt.Errorf("unexpected token: type= %s -> value= `%v`", p.lexer.currentToken.Kind.toString(), p.lexer.currentToken.Value)
}
