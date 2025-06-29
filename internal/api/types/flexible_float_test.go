package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlexibleFloat_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      float64
		expectErr bool
	}{
		{
			name:      "valid float number",
			input:     `123.45`,
			want:      123.45,
			expectErr: false,
		},
		{
			name:      "valid float string",
			input:     `"123.45"`,
			want:      123.45,
			expectErr: false,
		},
		{
			name:      "integer number",
			input:     `42`,
			want:      42.0,
			expectErr: false,
		},
		{
			name:      "integer string",
			input:     `"42"`,
			want:      42.0,
			expectErr: false,
		},
		{
			name:      "invalid string",
			input:     `"abc"`,
			expectErr: true,
		},
		{
			name:      "invalid json type (object)",
			input:     `{"foo": "bar"}`,
			expectErr: true,
		},
		{
			name:      "invalid json type (array)",
			input:     `[1,2,3]`,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var f FlexibleFloat
			err := json.Unmarshal([]byte(tt.input), &f)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.InDelta(t, tt.want, float64(f), 0.00001)
			}
		})
	}
}
