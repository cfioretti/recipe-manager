package application

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	balancerdomain "github.com/cfioretti/recipe-manager/internal/ingredients-balancer/domain"
)

func TestTotalDoughWeightByPans(t *testing.T) {
	tests := []struct {
		name     string
		input    balancerdomain.Pans
		wantArea float64
		wantErr  bool
	}{
		{
			name: "success with multiple pans",
			input: balancerdomain.Pans{
				Pans: []balancerdomain.Pan{
					{
						Shape: "round",
						Measures: balancerdomain.Measures{
							Diameter: intPtr(20),
						},
					},
					{
						Shape: "square",
						Measures: balancerdomain.Measures{
							Edge: intPtr(20),
						},
					},
				},
			},
			wantArea: (math.Pi*100 + 400) / 2,
			wantErr:  false,
		},
		{
			name: "invalid shape",
			input: balancerdomain.Pans{
				Pans: []balancerdomain.Pan{
					{
						Shape: "triangle",
						Measures: balancerdomain.Measures{
							Width:  intPtr(20),
							Length: intPtr(30),
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid measures",
			input: balancerdomain.Pans{
				Pans: []balancerdomain.Pan{
					{
						Shape: "round",
						Measures: balancerdomain.Measures{
							Diameter: nil,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "empty pans",
			input: balancerdomain.Pans{
				Pans: []balancerdomain.Pan{},
			},
			wantArea: 0,
			wantErr:  false,
		},
	}

	calculator := NewDoughCalculatorService()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculator.TotalDoughWeightByPans(tt.input)

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

func intPtr(value int) *int {
	return &value
}
