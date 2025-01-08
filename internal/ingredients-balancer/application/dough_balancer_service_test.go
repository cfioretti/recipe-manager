package application

import (
	"testing"

	"github.com/cfioretti/recipe-manager/internal/ingredients-balancer/domain"
	recipedomain "github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBalance(t *testing.T) {
	tests := []struct {
		name    string
		recipe  recipedomain.Recipe
		pans    domain.Pans
		want    *recipedomain.RecipeAggregate
		wantErr bool
	}{
		{
			name: "valid recipe and pans",
			recipe: recipedomain.Recipe{
				Id:   1,
				Uuid: uuid.New(),
				Name: "Test Recipe",
				Dough: recipedomain.Dough{
					Flour:  55.7,
					Water:  41.6,
					Salt:   1.1,
					EvoOil: 1.1,
					Yeast:  0.5,
				},
			},
			pans: domain.Pans{
				TotalDoughWeight: 1000,
				Pans: []domain.Pan{
					{
						Shape:       "round",
						DoughWeight: 500,
					},
					{
						Shape:       "round",
						DoughWeight: 500,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid total dough weight",
			recipe: recipedomain.Recipe{
				Dough: recipedomain.Dough{
					Flour: 55.7,
				},
			},
			pans: domain.Pans{
				TotalDoughWeight: 0,
			},
			wantErr: true,
		},
		{
			name: "invalid flour percentage",
			recipe: recipedomain.Recipe{
				Dough: recipedomain.Dough{
					Flour: 0,
				},
			},
			pans: domain.Pans{
				TotalDoughWeight: 1000,
			},
			wantErr: true,
		},
	}

	balancer := NewDoughBalancerService()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := balancer.Balance(tt.recipe, tt.pans)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)

			totalWeight := result.Dough.Flour +
				result.Dough.Water +
				result.Dough.Salt +
				result.Dough.EvoOil +
				result.Dough.Yeast

			assert.InDelta(t, tt.pans.TotalDoughWeight, totalWeight, 0.1)

			assert.InDelta(t, tt.recipe.Dough.Flour/100*tt.pans.TotalDoughWeight, result.Dough.Flour, 0.1)
			assert.InDelta(t, tt.recipe.Dough.Water/100*tt.pans.TotalDoughWeight, result.Dough.Water, 0.1)
		})
	}
}

func TestCalculateSplitDoughs(t *testing.T) {
	t.Run("multiple pans with proportional weights", func(t *testing.T) {
		totalDough := recipedomain.Dough{
			Flour:  1000.0,
			Water:  700.0,
			Salt:   20.0,
			EvoOil: 50.0,
			Yeast:  5.0,
		}

		pans := domain.Pans{
			TotalDoughWeight: 1000.0,
			Pans: []domain.Pan{
				{DoughWeight: 500.0},
				{DoughWeight: 300.0},
				{DoughWeight: 200.0},
			},
		}

		expected := []recipedomain.Dough{
			{Flour: 500.0, Water: 350.0, Salt: 10.0, EvoOil: 25.0, Yeast: 2.5},
			{Flour: 300.0, Water: 210.0, Salt: 6.0, EvoOil: 15.0, Yeast: 1.5},
			{Flour: 200.0, Water: 140.0, Salt: 4.0, EvoOil: 10.0, Yeast: 1.0},
		}

		result := calculateSplitDoughs(totalDough, pans)

		if len(result) != len(expected) {
			t.Fatalf("expected %d doughs, got %d", len(expected), len(result))
		}

		for i, dough := range result {
			if dough != expected[i] {
				t.Errorf("dough %d mismatch: expected %+v, got %+v", i, expected[i], dough)
			}
		}
	})

	t.Run("single pan", func(t *testing.T) {
		totalDough := recipedomain.Dough{
			Flour:  1000.0,
			Water:  700.0,
			Salt:   20.0,
			EvoOil: 50.0,
			Yeast:  5.0,
		}

		pans := domain.Pans{
			TotalDoughWeight: 1000.0,
			Pans: []domain.Pan{
				{DoughWeight: 1000.0},
			},
		}

		expected := []recipedomain.Dough{
			{Flour: 1000.0, Water: 700.0, Salt: 20.0, EvoOil: 50.0, Yeast: 5.0},
		}

		result := calculateSplitDoughs(totalDough, pans)

		if len(result) != len(expected) {
			t.Fatalf("expected %d doughs, got %d", len(expected), len(result))
		}

		if result[0] != expected[0] {
			t.Errorf("mismatch: expected %+v, got %+v", expected[0], result[0])
		}
	})

	t.Run("pans with zero weight", func(t *testing.T) {
		totalDough := recipedomain.Dough{
			Flour:  1000.0,
			Water:  700.0,
			Salt:   20.0,
			EvoOil: 50.0,
			Yeast:  5.0,
		}

		pans := domain.Pans{
			TotalDoughWeight: 1000.0,
			Pans: []domain.Pan{
				{DoughWeight: 0.0},
				{DoughWeight: 1000.0},
			},
		}

		expected := []recipedomain.Dough{
			{Flour: 0.0, Water: 0.0, Salt: 0.0, EvoOil: 0.0, Yeast: 0.0},
			{Flour: 1000.0, Water: 700.0, Salt: 20.0, EvoOil: 50.0, Yeast: 5.0},
		}

		result := calculateSplitDoughs(totalDough, pans)

		if len(result) != len(expected) {
			t.Fatalf("expected %d doughs, got %d", len(expected), len(result))
		}

		for i, dough := range result {
			if dough != expected[i] {
				t.Errorf("dough %d mismatch: expected %+v, got %+v", i, expected[i], dough)
			}
		}
	})

	t.Run("empty pans", func(t *testing.T) {
		totalDough := recipedomain.Dough{
			Flour:  1000.0,
			Water:  700.0,
			Salt:   20.0,
			EvoOil: 50.0,
			Yeast:  5.0,
		}

		pans := domain.Pans{
			TotalDoughWeight: 0,
			Pans:             []domain.Pan{},
		}

		expected := []recipedomain.Dough{}

		result := calculateSplitDoughs(totalDough, pans)

		if len(result) != len(expected) {
			t.Fatalf("expected %d doughs, got %d", len(expected), len(result))
		}
	})
}

func TestRound(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  float64
	}{
		{
			name:  "round up",
			input: 10.56,
			want:  10.6,
		},
		{
			name:  "round down",
			input: 10.54,
			want:  10.5,
		},
		{
			name:  "exact decimal",
			input: 10.50,
			want:  10.5,
		},
		{
			name:  "zero",
			input: 0.0,
			want:  0.0,
		},
		{
			name:  "negative number",
			input: -10.56,
			want:  -10.6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := round(tt.input)
			assert.Equal(t, tt.want, result)
		})
	}
}
