package scanner

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode"

	"go.uber.org/zap"
)

const (
	KeywordSelect = "select"
	KeywordFrom   = "from"
	KeywordWhere  = "where"
)

var (
	regexpNumber  = regexp.MustCompile(`^[0-9]+(.[0-9]+)?$`)
	regexpID      = regexp.MustCompile(`[a-zA-Z][a-zA-Z0-9_]*|\*`)
	regexpKeyword = regexp.MustCompile(`select|from|where`)
)

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

func NewScanner(logger *zap.Logger) *Scanner {
	return &Scanner{
		buf:       strings.Builder{},
		tokenizer: NewTokenizer(),
		logger:    logger,
	}
}

type Scanner struct {
	buf       strings.Builder
	tokenizer *Tokenizer
	logger    *zap.Logger
}

func (p *Scanner) Scan(reader io.RuneReader) ([]Token, error) {
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
			if err = p.flushBuffer(); err != nil {
				return nil, err
			}

			return p.tokenizer.GetTokens(), nil
		}

		// Парсинг строк
		if r == '\'' {
			if err = p.extractString(reader); err != nil {
				return nil, err
			}

			continue
		}

		// Парсинг наименования полей и таблиц
		if r == '"' {
			if err = p.extractID(reader); err != nil {
				return nil, err
			}

			continue
		}

		// ==== Обработка знаков сравнения
		if r == '=' || r == '>' || r == '<' {
			if err = p.handleCompareSign(r); err != nil {
				return nil, err
			}

			continue
		}

		// ==== Обработка скобок
		if r == '(' || r == ')' {
			if err = p.flushBuffer(); err != nil {
				return nil, err
			}
			tokenType := TokenTypeOpenCurlyBracket
			if r == ')' {
				tokenType = TokenTypeClosedCurlyBracket
			}
			p.buf.WriteRune(r)
			if err = p.tokenizer.AddToTokens(NewToken(p.buf.String(), tokenType)); err != nil {
				return nil, err
			}
			p.buf.Reset()

			continue
		}

		// ==== Обработка символов ключевых слов, наименования полей и таблиц
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '_' || r == '*' {
			p.buf.WriteRune(unicode.ToLower(r))

			continue
		}
		if err = p.flushBuffer(); err != nil {
			return nil, err
		}
	}
}

func (p *Scanner) flushBuffer() error {
	if p.buf.Len() > 0 {
		if err := p.tokenizer.AddToTokens(ParseTokenType(p.buf.String())); err != nil {
			return err
		}
		p.buf.Reset()
	}

	return nil
}

func (p *Scanner) extractString(reader io.RuneReader) error {
	var hasPrevRuneEscape bool

	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			return fmt.Errorf("string value should be closed")
		}
		if err != nil {
			return fmt.Errorf("unexpected error: %w", err)
		}

		// Строка закончилась
		if r == '\'' && !hasPrevRuneEscape {
			if p.buf.Len() > 0 {
				if err := p.tokenizer.AddToTokens(NewToken(p.buf.String(), TokenTypeString)); err != nil {
					return err
				}
				p.buf.Reset()
			}

			return nil
		}
		if r == '\\' && !hasPrevRuneEscape {
			hasPrevRuneEscape = true

			continue
		}
		hasPrevRuneEscape = false
		p.buf.WriteRune(r)
	}
}

func (p *Scanner) extractID(reader io.RuneReader) error {
	var hasPrevRuneEscape bool

	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			return fmt.Errorf("keyword should be closed by")
		}
		if err != nil {
			return fmt.Errorf("unexpected error: %w", err)
		}

		// Идентификатор закончен
		if r == '"' && !hasPrevRuneEscape {
			if p.buf.Len() > 0 {
				if err := p.tokenizer.AddToTokens(NewToken(p.buf.String(), TokenTypeID)); err != nil {
					return err
				}
				p.buf.Reset()
			}

			return nil
		}
		if r == '\\' && !hasPrevRuneEscape {
			hasPrevRuneEscape = true

			continue
		}
		hasPrevRuneEscape = false
		p.buf.WriteRune(r)
	}
}

func (p *Scanner) handleCompareSign(r rune) error {
	tokenType := TokenTypeOpEqual
	switch r {
	case '>':
		tokenType = TokenTypeOpMore
	case '<':
		tokenType = TokenTypeOpLess
	}
	if err := p.flushBuffer(); err != nil {
		return err
	}
	p.buf.WriteRune(r)
	if err := p.tokenizer.AddToTokens(NewToken(p.buf.String(), tokenType)); err != nil {
		return err
	}
	p.buf.Reset()

	return nil
}
