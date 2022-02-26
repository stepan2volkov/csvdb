package parser

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
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
					value:     "where",
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
					value:     "where",
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
					value:     "where",
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
					value:     "where",
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
					value:     "where",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.reader)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
