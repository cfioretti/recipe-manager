package strategies

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStrategy(t *testing.T) {
	tests := []struct {
		name    string
		shape   string
		wantErr bool
	}{
		{"round shape", "round", false},
		{"square shape", "square", false},
		{"rectangular shape", "rectangular", false},
		{"invalid shape", "triangle", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy, err := GetStrategy(tt.shape)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, strategy)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, strategy)
			}
		})
	}
}

func TestPanStrategies(t *testing.T) {
	tests := []struct {
		name     string
		strategy PanStrategy
		measures map[string]interface{}
		wantArea float64
		wantErr  bool
	}{
		{
			name:     "round pan",
			strategy: &RoundPanStrategy{},
			measures: map[string]interface{}{"diameter": "20"},
			wantArea: 157.07963267948966,
			wantErr:  false,
		},
		{
			name:     "square pan",
			strategy: &SquarePanStrategy{},
			measures: map[string]interface{}{"edge": "20"},
			wantArea: 200,
			wantErr:  false,
		},
		{
			name:     "rectangular pan",
			strategy: &RectangularPanStrategy{},
			measures: map[string]interface{}{"width": "20", "length": "30"},
			wantArea: 300,
			wantErr:  false,
		},
		{
			name:     "invalid measures",
			strategy: &RoundPanStrategy{},
			measures: map[string]interface{}{"diameter": "invalid"},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pan, err := tt.strategy.Calculate(tt.measures)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantArea, pan.DoughWeight)
		})
	}
}
