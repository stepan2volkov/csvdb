package value

import (
	"fmt"
	"testing"

	"github.com/stepan2volkov/csvdb/internal/app/table"
	"github.com/stretchr/testify/assert"
)

func TestNumberValue_Compare(t *testing.T) {
	tests := []struct {
		name string
		val1 string
		op   table.CompareOperationType
		val2 float64
		want bool
	}{
		{
			name: "number is not more",
			val1: "10",
			op:   table.CompareOperationTypeMore,
			val2: 15.0,
			want: false,
		},
		{
			name: "number is more",
			val1: "10",
			op:   table.CompareOperationTypeMore,
			val2: 5.0,
			want: true,
		},
		{
			name: "number is equal",
			val1: "10",
			op:   table.CompareOperationTypeEqual,
			val2: 10.0,
			want: true,
		},
		{
			name: "number is not less",
			val1: "10",
			op:   table.CompareOperationTypeLess,
			val2: 5.0,
			want: false,
		},
		{
			name: "number is less",
			val1: "10",
			op:   table.CompareOperationTypeLess,
			val2: 15.0,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, _ := NewNumberValue(tt.val1)
			got, _ := v.Compare(tt.val2, tt.op)
			assert.Equal(t, tt.want, got, fmt.Sprintf("%v %v %v", tt.val1, tt.op, tt.val2))
		})
	}
}
