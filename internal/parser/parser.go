package parser

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	regexpNumber  = regexp.MustCompile(`^[0-9]+(.[0-9]+)?$`)
	regexpID      = regexp.MustCompile(`[a-zA-Z][a-zA-Z0-9_]*`)
	regexpKeyword = regexp.MustCompile(`select|from|where`)
)

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

func NewTokenizer() *Tokenizer {
	return &Tokenizer{
		tokens: make([]Token, 0),
		stack:  make([]Token, 0),
	}
}

type Tokenizer struct {
	tokens []Token
	stack  []Token
}

func (t *Tokenizer) AddToTokens(token Token) error {
	if token.Type() == TokenTypeOpenCurlyBracket {
		fmt.Println("found")
	}
	switch token.Type() {
	case TokenTypeKeyword, TokenTypeID, TokenTypeString, TokenTypeNumber, TokenTypeUnknown:
		t.tokens = append(t.tokens, token)
	default:
		// Помещаем операции с большим или равным приоритетом в список токенов
		i := len(t.stack) - 1
		for i >= 0 && t.stack[i].priority >= token.priority && t.stack[i].Type() != TokenTypeOpenCurlyBracket {
			t.tokens = append(t.tokens, t.stack[i])
			i = i - 1
		}

		// Помещаем операции до открывающейся скобки в список токенов
		if token.Type() == TokenTypeClosedCurlyBracket {
			var openBracketFound bool

			for i >= 0 {
				openBracketFound = t.stack[i].Type() == TokenTypeOpenCurlyBracket
				if openBracketFound {
					i = i - 1
					break
				}
				t.tokens = append(t.tokens, t.stack[i])
				i = i - 1
			}
			if !openBracketFound {
				return fmt.Errorf("not found opened curly bracket")
			}
		}
		t.stack = t.stack[:i+1]
		if token.Type() != TokenTypeClosedCurlyBracket {
			t.stack = append(t.stack, token)
		}

	}
	return nil
}

func (t *Tokenizer) GetTokens() []Token {
	for i := len(t.stack) - 1; i >= 0; i-- {
		t.tokens = append(t.tokens, t.stack[i])
	}
	return t.tokens
}

func Parse(reader io.RuneReader) ([]Token, error) {
	tokenizer := NewTokenizer()

	var hasOpenedQuotationMark bool
	var hasPrevRuneEscape bool

	lexeme := strings.Builder{}

	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			return nil, fmt.Errorf("not found ';'")
		}
		if err != nil {
			return nil, fmt.Errorf("unexpected error when reading statement: %w", err)
		}

		// Конец выражения
		if r == ';' {
			if lexeme.Len() > 0 {
				if err := tokenizer.AddToTokens(ParseTokenType(lexeme.String())); err != nil {
					return nil, err
				}
				lexeme.Reset()
			}
			return tokenizer.GetTokens(), nil
		}

		// ==== Обработка строк

		// Строка закончилась
		if hasOpenedQuotationMark && (r == '\'' && !hasPrevRuneEscape) {
			hasOpenedQuotationMark = false
			if lexeme.Len() > 0 {
				if err := tokenizer.AddToTokens(NewToken(lexeme.String(), TokenTypeString)); err != nil {
					return nil, err
				}
				lexeme.Reset()
			}
			continue
		}
		// Строка продолжается
		if hasOpenedQuotationMark {
			if r == '\\' && !hasPrevRuneEscape {
				hasPrevRuneEscape = true
				continue
			}
			hasPrevRuneEscape = false
			lexeme.WriteRune(r)
			continue
		}
		// Строка началась
		if r == '\'' {
			hasOpenedQuotationMark = true
			continue
		}

		// ==== Обработка знаков сравнения
		if r == '=' || r == '>' || r == '<' {
			tokenType := TokenTypeOpEqual
			switch r {
			case '>':
				tokenType = TokenTypeOpMore
			case '<':
				tokenType = TokenTypeOpLess
			}
			if lexeme.Len() > 0 {
				if err := tokenizer.AddToTokens(ParseTokenType(lexeme.String())); err != nil {
					return nil, err
				}
				lexeme.Reset()
			}
			lexeme.WriteRune(r)
			if err := tokenizer.AddToTokens(NewToken(lexeme.String(), tokenType)); err != nil {
				return nil, err
			}
			lexeme.Reset()
			continue
		}

		// ==== Обработка скобок
		if r == '(' || r == ')' {
			if lexeme.Len() > 0 {
				if err := tokenizer.AddToTokens(ParseTokenType(lexeme.String())); err != nil {
					return nil, err
				}
				lexeme.Reset()
			}
			tokenType := TokenTypeOpenCurlyBracket
			if r == ')' {
				tokenType = TokenTypeClosedCurlyBracket
			}
			lexeme.WriteRune(r)
			if err := tokenizer.AddToTokens(NewToken(lexeme.String(), tokenType)); err != nil {
				return nil, err
			}
			lexeme.Reset()
			continue
		}

		// ==== Обработка символов ключевых слов, наименования полей и таблиц
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '_' {
			lexeme.WriteRune(unicode.ToLower(r))
			continue
		}
		if lexeme.Len() > 0 {
			if err := tokenizer.AddToTokens(ParseTokenType(lexeme.String())); err != nil {
				return nil, err
			}
			lexeme.Reset()
		}
	}
}
