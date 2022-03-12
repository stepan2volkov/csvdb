package scanner

import "strconv"

type TokenType int

const (
	TokenTypeUnknown            TokenType = iota
	TokenTypeKeyword            TokenType = iota
	TokenTypeID                 TokenType = iota
	TokenTypeString             TokenType = iota
	TokenTypeNumber             TokenType = iota
	TokenTypeOpMore             TokenType = iota
	TokenTypeOpLess             TokenType = iota
	TokenTypeOpEqual            TokenType = iota
	TokenTypeOpAnd              TokenType = iota
	TokenTypeOpOr               TokenType = iota
	TokenTypeOpenCurlyBracket   TokenType = iota
	TokenTypeClosedCurlyBracket TokenType = iota
)

func NewToken(value interface{}, tokenType TokenType) Token {
	priority := 0
	switch tokenType {
	case TokenTypeOpenCurlyBracket:
		priority = 4
	case TokenTypeClosedCurlyBracket:
		priority = 4
	case TokenTypeOpMore:
		priority = 3
	case TokenTypeOpLess:
		priority = 3
	case TokenTypeOpEqual:
		priority = 3
	case TokenTypeOpAnd:
		priority = 2
	case TokenTypeOpOr:
		priority = 1
	}

	return Token{
		tokenType: tokenType,
		priority:  priority,
		value:     value,
	}
}

func ParseTokenType(value string) Token {
	if value == "and" {
		return NewToken(value, TokenTypeOpAnd)
	}
	if value == "or" {
		return NewToken(value, TokenTypeOpOr)
	}
	if matched := regexpNumber.MatchString(value); matched {
		// можно игнорировать ошибку, поскольку значение проверено регулярным выражением
		val, _ := strconv.ParseFloat(value, 64)
		return NewToken(val, TokenTypeNumber)
	}
	if matched := regexpKeyword.MatchString(value); matched {
		return NewToken(value, TokenTypeKeyword)
	}
	if matched := regexpID.MatchString(value); matched {
		return NewToken(value, TokenTypeID)
	}

	return NewToken(value, TokenTypeUnknown)
}

type Token struct {
	tokenType TokenType
	priority  int
	value     interface{}
}

func (t *Token) Value() interface{} {
	return t.value
}

func (t *Token) Type() TokenType {
	return t.tokenType
}
