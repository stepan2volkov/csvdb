package parser

import (
	"strings"
	"testing"

	"github.com/stepan2volkov/csvdb/internal/app/scanner"
	"github.com/stepan2volkov/csvdb/internal/app/table"
	"github.com/stepan2volkov/csvdb/internal/app/table/operation"
	"github.com/stretchr/testify/assert"
)

func TestMakeWhere(t *testing.T) {
	tests := []struct {
		name    string
		stmt    string
		want    table.LogicalOperation
		wantErr error
	}{
		{
			name: "single condition",
			stmt: "col_1 > 2;",
			want: operation.DummyValueOperation{
				CompareOperation: table.CompareValueOperation{
					ColumnName: "col_1",
					Type:       table.CompareOperationTypeMore,
					Val:        2.0,
				},
			},
		},
		{
			name: "two conditions",
			stmt: "fullname = 'Mike Smith' and age > 18;",
			want: operation.AndOperation{
				Left: operation.DummyValueOperation{
					CompareOperation: table.CompareValueOperation{
						ColumnName: "age",
						Type:       table.CompareOperationTypeMore,
						Val:        18.0,
					},
				},
				Right: operation.DummyValueOperation{
					CompareOperation: table.CompareValueOperation{
						ColumnName: "fullname",
						Type:       table.CompareOperationTypeEqual,
						Val:        "Mike Smith",
					},
				},
			},
		},
		{
			name: "three conditions",
			stmt: "fullname = 'Mike Smith' and (age > 18 or salary > 15000.99;",
			want: operation.AndOperation{
				Left: operation.OrOperation{
					Left: operation.DummyValueOperation{
						CompareOperation: table.CompareValueOperation{
							ColumnName: "salary",
							Type:       table.CompareOperationTypeMore,
							Val:        15000.99,
						},
					},
					Right: operation.DummyValueOperation{
						CompareOperation: table.CompareValueOperation{
							ColumnName: "age",
							Type:       table.CompareOperationTypeMore,
							Val:        18.0,
						},
					},
				},
				Right: operation.DummyValueOperation{
					CompareOperation: table.CompareValueOperation{
						ColumnName: "fullname",
						Type:       table.CompareOperationTypeEqual,
						Val:        "Mike Smith",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := scanner.NewScanner().Scan(strings.NewReader(tt.stmt))
			assert.ErrorIs(t, err, nil)
			got, err := makeWhere(tokens)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
