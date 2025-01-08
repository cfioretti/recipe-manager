package application

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTotalDoughWeightByPans(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantArea float64
		wantErr  bool
	}{
		{
			name: "success with multiple pans",
			input: `{
                "pans": [
                    {"shape": "round", "measures": {"diameter": "20"}},
                    {"shape": "square", "measures": {"edge": "20"}}
                ]
            }`,
			wantArea: (math.Pi*100 + 400) / 2,
			wantErr:  false,
		},
		{
			name: "invalid shape",
			input: `{
                "pans": [{
                    "shape": "triangle",
                    "measures": {"base": "20", "height": "30"}
                }]
            }`,
			wantErr: true,
		},
		{
			name: "invalid measures",
			input: `{
                "pans": [{
                    "shape": "round",
                    "measures": {"diameter": "invalid"}
                }]
            }`,
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			input:   `{invalid json}`,
			wantErr: true,
		},
	}

	calculator := NewDoughCalculatorService()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculator.TotalDoughWeightByPans([]byte(tt.input))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.InDelta(t, tt.wantArea, result.TotalDoughWeight, 0.001)
		})
	}
}
