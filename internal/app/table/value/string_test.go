package value

import (
	"fmt"
	"testing"

	"github.com/stepan2volkov/csvdb/internal/app/table"
	"github.com/stretchr/testify/assert"
)

func TestStringValue_Compare(t *testing.T) {
	tests := []struct {
		name string
		val1 string
		op   table.CompareOperationType
		val2 string
		want bool
	}{
		{
			name: "string is equal",
			val1: "hello",
			op:   table.CompareOperationTypeEqual,
			val2: "hello",
			want: true,
		},
		{
			name: "string is not equal",
			val1: "hello",
			op:   table.CompareOperationTypeEqual,
			val2: "world",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewStringValue(tt.val1)
			got, _ := v.Compare(tt.val2, tt.op)
			assert.Equal(t, tt.want, got, fmt.Sprintf("'%v' %v '%v'", tt.val1, tt.op, tt.val2))
		})
	}
}
