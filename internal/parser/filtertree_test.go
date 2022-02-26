package parser

import (
	"testing"

	"github.com/stepan2volkov/csvdb/internal/table"
	"github.com/stretchr/testify/assert"
)

func TestMakeFilters(t *testing.T) {
	tests := []struct {
		name    string
		tokens  []Token
		want    table.LogicalFilter
		wantErr error
	}{
		{
			name: "single condition",
			tokens: []Token{
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
			want: table.DummyFilter{
				Filter: table.CompareFilter{
					FieldName: "col_1",
					Op:        table.CompareOperationTypeMore,
					Val:       2.0,
				},
			},
		},
		{
			name: "two conditions",
			tokens: []Token{
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
			want: table.AndFilter{
				Left: table.DummyFilter{
					Filter: table.CompareFilter{
						FieldName: "age",
						Op:        table.CompareOperationTypeMore,
						Val:       18.0,
					},
				},
				Right: table.DummyFilter{
					Filter: table.CompareFilter{
						FieldName: "fullname",
						Op:        table.CompareOperationTypeEqual,
						Val:       "Mike Smith",
					},
				},
			},
		},
		{
			name: "three conditions",
			tokens: []Token{
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
			want: table.AndFilter{
				Left: table.OrFilter{
					Left: table.DummyFilter{
						Filter: table.CompareFilter{
							FieldName: "salary",
							Op:        table.CompareOperationTypeMore,
							Val:       15000.99,
						},
					},
					Right: table.DummyFilter{
						Filter: table.CompareFilter{
							FieldName: "age",
							Op:        table.CompareOperationTypeMore,
							Val:       18.0,
						},
					},
				},
				Right: table.DummyFilter{
					Filter: table.CompareFilter{
						FieldName: "fullname",
						Op:        table.CompareOperationTypeEqual,
						Val:       "Mike Smith",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MakeFilters(tt.tokens)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
