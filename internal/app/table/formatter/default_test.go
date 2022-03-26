package formatter

import (
	"context"
	"testing"

	"github.com/stepan2volkov/csvdb/internal/app/table"
	"github.com/stepan2volkov/csvdb/internal/app/table/value"
	"github.com/stretchr/testify/assert"
)

func TestDefaultFormatter_Format(t *testing.T) {
	tests := []struct {
		name        string
		t           table.Table
		want        string
		expectedErr error
	}{
		{
			name: "pretty print test",
			t: table.NewTable("sales", []table.Column{
				{
					Field: table.Field{
						Name: "Region",
						Type: table.FieldTypeString,
					},
					Values: []table.Value{
						value.NewStringValue("Africa"),
						value.NewStringValue("USA"),
						value.NewStringValue("England"),
					},
				},
			}),
			want: `┌─────────┐
│ REGION  │
├─────────┤
│ Africa  │
│ USA     │
│ England │
└─────────┘`,
		},
	}

	f := &DefaultFormatter{}
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := f.Format(ctx, tt.t)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
