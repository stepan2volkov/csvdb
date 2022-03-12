package parser

import (
	"strings"
	"testing"

	"github.com/stepan2volkov/csvdb/internal/app/scanner"
	"github.com/stepan2volkov/csvdb/internal/app/table"
	"github.com/stepan2volkov/csvdb/internal/app/table/operation"
	"github.com/stretchr/testify/assert"
)

func TestMakeSelectStmt(t *testing.T) {
	tests := []struct {
		name    string
		stmt    string
		want    SelectStmt
		wantErr error
	}{
		{
			name: "with where",
			stmt: "select col_1, col_2 from table where col_1 > 2;",
			want: SelectStmt{
				Fields:    []string{"col_1", "col_2"},
				Tablename: "table",
				Filter: operation.DummyValueOperation{
					CompareOperation: table.CompareValueOperation{
						ColumnName: "col_1",
						Type:       table.CompareOperationTypeMore,
						Val:        2.0,
					},
				},
			},
		},
		{
			name: "select all",
			stmt: "SELECT * FROM table;",
			want: SelectStmt{
				Fields:    nil,
				AllField:  true,
				Tablename: "table",
				Filter: operation.DummyValueOperation{
					CompareOperation: table.CompareValueOperation{
						Type: table.CompareOperationTypeDummy,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := scanner.NewScanner().Scan(strings.NewReader(tt.stmt))
			assert.Equal(t, err, nil)
			got, err := MakeSelectStmt(tokens)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
