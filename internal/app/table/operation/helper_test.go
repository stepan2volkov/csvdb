package operation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeIndexes(t *testing.T) {
	tests := []struct {
		name   string
		first  []int
		second []int
		want   []int
	}{
		{
			name:   "sets without intersection",
			first:  []int{1, 2, 3, 4},
			second: []int{5, 6, 7, 8},
			want:   []int{1, 2, 3, 4, 5, 6, 7, 8},
		},
		{
			name:   "sets with intersection",
			first:  []int{1, 2, 5, 6},
			second: []int{5, 6, 7, 8},
			want:   []int{1, 2, 5, 6, 7, 8},
		},
		{
			name:   "saving order",
			first:  []int{1, 3, 5, 7},
			second: []int{2, 4, 6, 8},
			want:   []int{1, 2, 3, 4, 5, 6, 7, 8},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeIndexes(tt.first, tt.second)
			assert.Equal(t, tt.want, got)
		})
	}
}
