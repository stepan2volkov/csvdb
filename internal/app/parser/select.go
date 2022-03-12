package parser

import (
	"fmt"

	"github.com/stepan2volkov/csvdb/internal/app/scanner"
	"github.com/stepan2volkov/csvdb/internal/app/table"
	"github.com/stepan2volkov/csvdb/internal/app/table/operation"
)

const (
	KeywordSelect = "select"
	KeywordFrom   = "from"
	KeywordWhere  = "where"
)

type SelectStmt struct {
	Fields    []string
	AllField  bool
	Tablename string
	Filter    table.LogicalOperation
}

func MakeSelectStmt(tokens []scanner.Token) (SelectStmt, error) {
	var builder selectStmtBuilder
	for _, token := range tokens {
		if err := builder.append(token); err != nil {
			return SelectStmt{}, err
		}
	}
	return builder.build()
}

type selectStmtBuilder struct {
	lastKeyword string
	allFields   bool
	fields      []string
	tablename   string
	conditions  []scanner.Token
}

func (b *selectStmtBuilder) build() (SelectStmt, error) {
	var filter table.LogicalOperation = operation.DummyValueOperation{
		CompareOperation: table.CompareValueOperation{
			Type: table.CompareOperationTypeDummy,
		},
	}
	if len(b.conditions) > 0 {
		newFilter, err := makeWhere(b.conditions)
		if err != nil {
			return SelectStmt{}, err
		}
		filter = newFilter
	}

	return SelectStmt{
		Fields:    b.fields,
		AllField:  b.allFields,
		Tablename: b.tablename,
		Filter:    filter,
	}, nil
}

func (b *selectStmtBuilder) append(token scanner.Token) error {
	if token.Type() == scanner.TokenTypeKeyword {
		value := token.Value().(string)
		switch value {
		case KeywordSelect:
			if b.lastKeyword != "" {
				return fmt.Errorf("select should be the first word")
			}
		case KeywordFrom:
			if b.lastKeyword != KeywordSelect {
				return fmt.Errorf("from section should be after select")
			}
			if b.fields == nil && !b.allFields {
				return fmt.Errorf("fields should be specified after select")
			}
		case KeywordWhere:
			if b.lastKeyword != KeywordFrom {
				return fmt.Errorf("where section should be after from")
			}
			if b.tablename == "" {
				return fmt.Errorf("tablename should be specified after from")
			}
		}
		b.lastKeyword = value
		return nil
	}
	if token.Type() == scanner.TokenTypeID {
		switch b.lastKeyword {
		case KeywordSelect:
			if b.allFields {
				return fmt.Errorf("invalid format of select stmt")
			}
			value := token.Value().(string)
			if value == "*" {
				b.allFields = true
				return nil
			}
			b.fields = append(b.fields, value)
		case KeywordFrom:
			if b.tablename != "" {
				return fmt.Errorf("tablename should be specified once")
			}
			b.tablename = token.Value().(string)
		case KeywordWhere:
			b.conditions = append(b.conditions, token)
		default:
			return fmt.Errorf("select should be the first word")
		}
		return nil
	}
	if b.lastKeyword != KeywordWhere {
		return fmt.Errorf("invalid format of select stmt")
	}
	b.conditions = append(b.conditions, token)
	return nil
}
