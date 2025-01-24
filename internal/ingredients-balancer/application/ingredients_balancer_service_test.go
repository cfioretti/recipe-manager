package application

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/cfioretti/recipe-manager/internal/ingredients-balancer/domain"
	recipedomain "github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
)

func TestBalance(t *testing.T) {
	tests := []struct {
		name             string
		recipe           recipedomain.Recipe
		pans             domain.Pans
		totalDoughWeight float64
		wantErr          bool
	}{
		{
			name: "valid recipe and pans without percent variation",
			recipe: recipedomain.Recipe{
				Id:   1,
				Uuid: uuid.New(),
				Name: "Test Recipe",
				Dough: recipedomain.Dough{
					PercentVariation: 0,
					Ingredients: []recipedomain.Ingredient{
						{Name: "flour", Amount: 55.7},
						{Name: "water", Amount: 41.6},
						{Name: "salt", Amount: 1.1},
						{Name: "evoOil", Amount: 1.1},
						{Name: "yeast", Amount: 0.5},
					},
				},
			},
			pans: domain.Pans{
				TotalArea: 1000,
				Pans: []domain.Pan{
					{
						Shape: "round",
						Area:  500,
					},
					{
						Shape: "round",
						Area:  500,
					},
				},
			},
			totalDoughWeight: 500,
			wantErr:          false,
		},
		{
			name: "valid recipe and pans with percent variation",
			recipe: recipedomain.Recipe{
				Id:   1,
				Uuid: uuid.New(),
				Name: "Test Recipe",
				Dough: recipedomain.Dough{
					PercentVariation: 10,
					Ingredients: []recipedomain.Ingredient{
						{Name: "flour", Amount: 55.7},
						{Name: "water", Amount: 41.6},
						{Name: "salt", Amount: 1.1},
						{Name: "evoOil", Amount: 1.1},
						{Name: "yeast", Amount: 0.5},
					},
				},
			},
			pans: domain.Pans{
				TotalArea: 2000,
				Pans: []domain.Pan{
					{
						Shape: "square",
						Area:  1000,
					},
					{
						Shape: "square",
						Area:  1000,
					},
				},
			},
			totalDoughWeight: 1100,
			wantErr:          false,
		},
		{
			name: "invalid total dough weight",
			recipe: recipedomain.Recipe{
				Dough: recipedomain.Dough{
					Ingredients: []recipedomain.Ingredient{
						{Name: "flour", Amount: 55.7},
					},
				},
			},
			pans: domain.Pans{
				TotalArea: 0,
			},
			wantErr: true,
		},
		{
			name: "empty ingredients",
			recipe: recipedomain.Recipe{
				Dough: recipedomain.Dough{
					Ingredients: []recipedomain.Ingredient{},
				},
			},
			pans: domain.Pans{
				TotalArea: 1000,
			},
			wantErr: true,
		},
	}

	balancer := NewIngredientsBalancerService()

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

			totalWeight := 0.0
			for _, ing := range result.Dough.Ingredients {
				totalWeight += ing.Amount
			}
			assert.InDelta(t, tt.totalDoughWeight, totalWeight, 0.1)

			firstIngredientRatio := getFirstIngredientAmount(tt.recipe.Dough.Ingredients) / 100
			expectedAmount := firstIngredientRatio * tt.totalDoughWeight
			actualAmount := getFirstIngredientAmount(result.Dough.Ingredients)
			assert.InDelta(t, expectedAmount, actualAmount, 0.1)
		})
	}
}

func TestCalculateSplitDoughs(t *testing.T) {
	t.Run("multiple pans with proportional weights", func(t *testing.T) {
		totalDough := recipedomain.Dough{
			Ingredients: []recipedomain.Ingredient{
				{Name: "flour", Amount: 1000.0},
				{Name: "water", Amount: 700.0},
				{Name: "salt", Amount: 20.0},
				{Name: "evoOil", Amount: 50.0},
				{Name: "yeast", Amount: 5.0},
			},
		}

		pans := domain.Pans{
			TotalArea: 2000.0,
			Pans: []domain.Pan{
				{Area: 1000.0},
				{Area: 600.0},
				{Area: 400.0},
			},
		}

		result := calculateSplitDoughs(totalDough, pans)
		assert.Len(t, result, len(pans.Pans))

		for i, pan := range pans.Pans {
			ratio := pan.Area / pans.TotalArea
			for j, ingredient := range totalDough.Ingredients {
				expectedAmount := round(ingredient.Amount * ratio)
				assert.InDelta(t, expectedAmount, result[i].Ingredients[j].Amount, 0.1)
			}
		}
	})

	t.Run("single pan", func(t *testing.T) {
		totalDough := recipedomain.Dough{
			Ingredients: []recipedomain.Ingredient{
				{Name: "flour", Amount: 1000.0},
				{Name: "water", Amount: 700.0},
			},
		}

		pans := domain.Pans{
			TotalArea: 1000.0,
			Pans: []domain.Pan{
				{Area: 1000.0},
			},
		}

		result := calculateSplitDoughs(totalDough, pans)
		assert.Len(t, result, 1)

		for i, ingredient := range totalDough.Ingredients {
			assert.Equal(t, ingredient.Amount, result[0].Ingredients[i].Amount)
		}
	})

	t.Run("empty pans", func(t *testing.T) {
		totalDough := recipedomain.Dough{
			Ingredients: []recipedomain.Ingredient{
				{Name: "flour", Amount: 1000.0},
			},
		}

		pans := domain.Pans{
			TotalArea: 0,
			Pans:      []domain.Pan{},
		}

		result := calculateSplitDoughs(totalDough, pans)
		assert.Empty(t, result)
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
