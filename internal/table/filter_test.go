package table

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mergeIndexes(t *testing.T) {
	tests := []struct {
		name   string
		first  []int
		second []int
		want   []int
	}{
		{
			name:   "merging without duplicates",
			first:  []int{0, 2, 4, 6},
			second: []int{1, 3, 5},
			want:   []int{0, 1, 2, 3, 4, 5, 6},
		},
		{
			name:   "merging with duplicates",
			first:  []int{0, 1, 2, 4, 6},
			second: []int{1, 3, 5},
			want:   []int{0, 1, 2, 3, 4, 5, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeIndexes(tt.first, tt.second)
			assert.Equal(t, tt.want, got)
		})
	}
}
