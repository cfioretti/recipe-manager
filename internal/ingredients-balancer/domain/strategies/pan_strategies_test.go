package strategies

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cfioretti/recipe-manager/internal/ingredients-balancer/domain"
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
		measures domain.Measures
		wantArea float64
		wantErr  bool
	}{
		{
			name:     "round pan",
			strategy: &RoundPanStrategy{},
			measures: domain.Measures{Diameter: intPtr(20)},
			wantArea: 157.07963267948966,
			wantErr:  false,
		},
		{
			name:     "square pan",
			strategy: &SquarePanStrategy{},
			measures: domain.Measures{Edge: intPtr(20)},
			wantArea: 200,
			wantErr:  false,
		},
		{
			name:     "rectangular pan",
			strategy: &RectangularPanStrategy{},
			measures: domain.Measures{Width: intPtr(20), Length: intPtr(30)},
			wantArea: 300,
			wantErr:  false,
		},
		{
			name:     "invalid measures",
			strategy: &RoundPanStrategy{},
			measures: domain.Measures{},
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

func intPtr(i int) *int {
	return &i
}
