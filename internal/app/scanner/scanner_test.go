package scanner

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestScan(t *testing.T) {
	tests := []struct {
		name    string
		reader  io.RuneReader
		want    []Token
		wantErr error
	}{
		{
			name:   "first condition",
			reader: strings.NewReader("WHERE col_1> 2;"),
			want: []Token{
				{
					tokenType: TokenTypeKeyword,
					value:     KeywordWhere,
					priority:  0,
				},
				{
					tokenType: TokenTypeID,
					value:     "col_1",
					priority:  0,
				},
				{
					tokenType: TokenTypeNumber,
					value:     2.0,
					priority:  0,
				},
				{
					tokenType: TokenTypeOpMore,
					value:     ">",
					priority:  3,
				},
			},
		},
		{
			name:   "second condition",
			reader: strings.NewReader(`WHERE fullname ='Mike';`),
			want: []Token{
				{
					tokenType: TokenTypeKeyword,
					value:     KeywordWhere,
					priority:  0,
				},
				{
					tokenType: TokenTypeID,
					value:     "fullname",
					priority:  0,
				},
				{
					tokenType: TokenTypeString,
					value:     "Mike",
					priority:  0,
				},
				{
					tokenType: TokenTypeOpEqual,
					value:     "=",
					priority:  3,
				},
			},
		},
		{
			name:   "two conditions",
			reader: strings.NewReader(`WHERE fullname ='Mike Smith' AND age > 18;`),
			want: []Token{
				{
					tokenType: TokenTypeKeyword,
					value:     KeywordWhere,
					priority:  0,
				},
				{
					tokenType: TokenTypeID,
					value:     "fullname",
					priority:  0,
				},
				{
					tokenType: TokenTypeString,
					value:     "Mike Smith",
					priority:  0,
				},
				{
					tokenType: TokenTypeOpEqual,
					value:     "=",
					priority:  3,
				},
				{
					tokenType: TokenTypeID,
					value:     "age",
					priority:  0,
				},
				{
					tokenType: TokenTypeNumber,
					value:     18.0,
					priority:  0,
				},
				{
					tokenType: TokenTypeOpMore,
					value:     ">",
					priority:  3,
				},
				{
					tokenType: TokenTypeOpAnd,
					value:     "and",
					priority:  2,
				},
			},
		},
		{
			name:   "curly brackets",
			reader: strings.NewReader(`WHERE fullname ='Mike Smith' AND (age > 18 OR salary > 15000.99);`),
			want: []Token{
				{
					tokenType: TokenTypeKeyword,
					value:     KeywordWhere,
					priority:  0,
				},
				{
					tokenType: TokenTypeID,
					value:     "fullname",
					priority:  0,
				},
				{
					tokenType: TokenTypeString,
					value:     "Mike Smith",
					priority:  0,
				},
				{
					tokenType: TokenTypeOpEqual,
					value:     "=",
					priority:  3,
				},
				{
					tokenType: TokenTypeID,
					value:     "age",
					priority:  0,
				},
				{
					tokenType: TokenTypeNumber,
					value:     18.0,
					priority:  0,
				},
				{
					tokenType: TokenTypeOpMore,
					value:     ">",
					priority:  3,
				},
				{
					tokenType: TokenTypeID,
					value:     "salary",
					priority:  0,
				},
				{
					tokenType: TokenTypeNumber,
					value:     15000.99,
					priority:  0,
				},
				{
					tokenType: TokenTypeOpMore,
					value:     ">",
					priority:  3,
				},
				{
					tokenType: TokenTypeOpOr,
					value:     "or",
					priority:  1,
				},
				{
					tokenType: TokenTypeOpAnd,
					value:     "and",
					priority:  2,
				},
			},
		},
		{
			name:   "without curly brackets",
			reader: strings.NewReader(`WHERE fullname ='Mike Smith' OR age > 18 AND salary > 15000.99;`),
			want: []Token{

				{
					tokenType: TokenTypeKeyword,
					value:     KeywordWhere,
					priority:  0,
				},
				{
					tokenType: TokenTypeID,
					value:     "fullname",
					priority:  0,
				},
				{
					tokenType: TokenTypeString,
					value:     "Mike Smith",
					priority:  0,
				},
				{
					tokenType: TokenTypeOpEqual,
					value:     "=",
					priority:  3,
				},
				{
					tokenType: TokenTypeID,
					value:     "age",
					priority:  0,
				},
				{
					tokenType: TokenTypeNumber,
					value:     18.0,
					priority:  0,
				},
				{
					tokenType: TokenTypeOpMore,
					value:     ">",
					priority:  3,
				},
				{
					tokenType: TokenTypeID,
					value:     "salary",
					priority:  0,
				},
				{
					tokenType: TokenTypeNumber,
					value:     15000.99,
					priority:  0,
				},
				{
					tokenType: TokenTypeOpMore,
					value:     ">",
					priority:  3,
				},
				{
					tokenType: TokenTypeOpAnd,
					value:     "and",
					priority:  2,
				},
				{
					tokenType: TokenTypeOpOr,
					value:     "or",
					priority:  1,
				},
			},
		},
		{
			name:   "select stmt",
			reader: strings.NewReader("SELECT col_1, col_2 FROM table WHERE col_1 > 2;"),
			want: []Token{
				{
					tokenType: TokenTypeKeyword,
					value:     KeywordSelect,
					priority:  0,
				},
				{
					tokenType: TokenTypeID,
					value:     "col_1",
					priority:  0,
				},
				{
					tokenType: TokenTypeID,
					value:     "col_2",
					priority:  0,
				},
				{
					tokenType: TokenTypeKeyword,
					value:     KeywordFrom,
					priority:  0,
				},
				{
					tokenType: TokenTypeID,
					value:     "table",
					priority:  0,
				},
				{
					tokenType: TokenTypeKeyword,
					value:     KeywordWhere,
					priority:  0,
				},
				{
					tokenType: TokenTypeID,
					value:     "col_1",
					priority:  0,
				},
				{
					tokenType: TokenTypeNumber,
					value:     2.0,
					priority:  0,
				},
				{
					tokenType: TokenTypeOpMore,
					value:     ">",
					priority:  3,
				},
			},
		},
		{
			name:   "select all",
			reader: strings.NewReader("SELECT * FROM table;"),
			want: []Token{
				{
					tokenType: TokenTypeKeyword,
					value:     KeywordSelect,
					priority:  0,
				},
				{
					tokenType: TokenTypeID,
					value:     "*",
					priority:  0,
				},
				{
					tokenType: TokenTypeKeyword,
					value:     KeywordFrom,
					priority:  0,
				},
				{
					tokenType: TokenTypeID,
					value:     "table",
					priority:  0,
				},
			},
		},
		{
			name:   "id double quotes ",
			reader: strings.NewReader(`SELECT "Col Name" FROM "tableName";`),
			want: []Token{
				{
					tokenType: TokenTypeKeyword,
					value:     KeywordSelect,
					priority:  0,
				},
				{
					tokenType: TokenTypeID,
					value:     "Col Name",
					priority:  0,
				},
				{
					tokenType: TokenTypeKeyword,
					value:     KeywordFrom,
					priority:  0,
				},
				{
					tokenType: TokenTypeID,
					value:     "tableName",
					priority:  0,
				},
			},
		},
	}

	logger, _ := zap.NewDevelopment()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewScanner(logger)
			got, err := parser.Scan(tt.reader)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
