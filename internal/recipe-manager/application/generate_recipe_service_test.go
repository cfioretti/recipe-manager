package application

import (
	"context"
	"errors"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
)

type MockRecipeGenerationClient struct {
	mock.Mock
}

func (m *MockRecipeGenerationClient) GenerateRecipe(ctx context.Context, command GenerateRecipeCommand) (*GeneratedRecipeDraft, error) {
	args := m.Called(ctx, command)
	draft, _ := args.Get(0).(*GeneratedRecipeDraft)
	return draft, args.Error(1)
}

func TestGenerateRecipeService_Handle(t *testing.T) {
	ctx := context.Background()
	command := GenerateRecipeCommand{
		Mode:   "prompt",
		Prompt: "high hydration",
		Pans: domain.Pans{
			Pans: []domain.Pan{
				{
					Shape: "round",
					Measures: domain.Measures{
						Diameter: intPtr(30),
					},
				},
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		mockGenerationClient := new(MockRecipeGenerationClient)
		mockCalculatorService := new(MockCalculatorService)
		mockBalancerService := new(MockBalancerService)

		draft := &GeneratedRecipeDraft{
			Name:        "AI Pizza",
			Description: "Generated",
			Author:      "AI",
			Dough: map[string]float64{
				"flour": 1000,
				"water": 700,
			},
			Topping: map[string]float64{
				"mozzarella": 300,
			},
		}

		pans := domain.Pans{}
		aggregate := &domain.RecipeAggregate{
			Recipe: domain.Recipe{
				Name: "AI Pizza",
			},
		}

		mockGenerationClient.On("GenerateRecipe", mock.Anything, command).Return(draft, nil)
		mockCalculatorService.On("TotalDoughWeightByPans", mock.Anything, command.Pans).Return(&pans, nil)
		mockBalancerService.On("Balance", mock.Anything, mock.Anything, pans).Return(aggregate, nil)

		service := NewGenerateRecipeService(mockGenerationClient, mockCalculatorService, mockBalancerService)
		result, err := service.Handle(ctx, command)

		assert.NoError(t, err)
		assert.Equal(t, aggregate, result)
	})

	t.Run("invalid mode", func(t *testing.T) {
		service := NewGenerateRecipeService(new(MockRecipeGenerationClient), new(MockCalculatorService), new(MockBalancerService))
		_, err := service.Handle(ctx, GenerateRecipeCommand{
			Mode: "unsupported",
			Pans: domain.Pans{},
		})

		assert.ErrorIs(t, err, ErrInvalidGenerationMode)
	})

	t.Run("generation client error", func(t *testing.T) {
		mockGenerationClient := new(MockRecipeGenerationClient)
		mockCalculatorService := new(MockCalculatorService)
		mockBalancerService := new(MockBalancerService)
		expectedErr := errors.New("generation failure")

		mockGenerationClient.On("GenerateRecipe", mock.Anything, command).Return((*GeneratedRecipeDraft)(nil), expectedErr)

		service := NewGenerateRecipeService(mockGenerationClient, mockCalculatorService, mockBalancerService)
		_, err := service.Handle(ctx, command)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), expectedErr.Error())
	})

	t.Run("invalid draft structure when required dough ingredients are missing", func(t *testing.T) {
		mockGenerationClient := new(MockRecipeGenerationClient)
		mockCalculatorService := new(MockCalculatorService)
		mockBalancerService := new(MockBalancerService)

		draft := &GeneratedRecipeDraft{
			Name:        "AI Pizza",
			Description: "Generated",
			Author:      "AI",
			Dough: map[string]float64{
				"ingredient": 2.5,
			},
			Topping: map[string]float64{
				"mozzarella": 120,
			},
		}

		mockGenerationClient.On("GenerateRecipe", mock.Anything, command).Return(draft, nil)

		service := NewGenerateRecipeService(mockGenerationClient, mockCalculatorService, mockBalancerService)
		_, err := service.Handle(ctx, command)

		assert.ErrorIs(t, err, ErrInvalidDraftStructure)
	})

	t.Run("invalid balanced recipe when non-finite values are returned", func(t *testing.T) {
		mockGenerationClient := new(MockRecipeGenerationClient)
		mockCalculatorService := new(MockCalculatorService)
		mockBalancerService := new(MockBalancerService)

		draft := &GeneratedRecipeDraft{
			Name:        "AI Pizza",
			Description: "Generated",
			Author:      "AI",
			Dough: map[string]float64{
				"flour": 1000,
				"water": 700,
			},
			Topping: map[string]float64{
				"mozzarella": 300,
			},
		}
		pans := domain.Pans{}
		aggregate := &domain.RecipeAggregate{
			Recipe: domain.Recipe{
				Name: "AI Pizza",
				Dough: domain.Dough{
					Ingredients: []domain.Ingredient{
						{Name: "flour", Amount: math.NaN()},
					},
				},
			},
		}

		mockGenerationClient.On("GenerateRecipe", mock.Anything, command).Return(draft, nil)
		mockCalculatorService.On("TotalDoughWeightByPans", mock.Anything, command.Pans).Return(&pans, nil)
		mockBalancerService.On("Balance", mock.Anything, mock.Anything, pans).Return(aggregate, nil)

		service := NewGenerateRecipeService(mockGenerationClient, mockCalculatorService, mockBalancerService)
		_, err := service.Handle(ctx, command)

		assert.ErrorIs(t, err, ErrInvalidBalancedRecipe)
	})

	t.Run("draft with uppercase keys passes semantic validation", func(t *testing.T) {
		mockGenerationClient := new(MockRecipeGenerationClient)
		mockCalculatorService := new(MockCalculatorService)
		mockBalancerService := new(MockBalancerService)

		draft := &GeneratedRecipeDraft{
			Name:        "AI Pizza",
			Description: "Generated",
			Author:      "AI",
			Dough: map[string]float64{
				"Flour": 1000,
				"Water": 700,
			},
			Topping: map[string]float64{
				"Mozzarella": 300,
			},
		}

		pans := domain.Pans{}
		aggregate := &domain.RecipeAggregate{
			Recipe: domain.Recipe{
				Name: "AI Pizza",
			},
		}

		mockGenerationClient.On("GenerateRecipe", mock.Anything, command).Return(draft, nil)
		mockCalculatorService.On("TotalDoughWeightByPans", mock.Anything, command.Pans).Return(&pans, nil)
		mockBalancerService.On("Balance", mock.Anything, mock.Anything, pans).Return(aggregate, nil)

		service := NewGenerateRecipeService(mockGenerationClient, mockCalculatorService, mockBalancerService)
		result, err := service.Handle(ctx, command)

		assert.NoError(t, err)
		assert.Equal(t, aggregate, result)
	})
}

func TestToBakerPercentages(t *testing.T) {
	t.Run("converts gram amounts to baker percentages", func(t *testing.T) {
		dough := map[string]float64{
			"flour": 1000,
			"water": 700,
			"salt":  25,
			"yeast": 3,
		}

		result := toBakerPercentages(dough)

		amounts := map[string]float64{}
		for _, ing := range result {
			amounts[ing.Name] = ing.Amount
		}

		assert.Equal(t, 100.0, amounts["flour"])
		assert.Equal(t, 70.0, amounts["water"])
		assert.Equal(t, 2.5, amounts["salt"])
		assert.Equal(t, 0.3, amounts["yeast"])
	})

	t.Run("handles missing flour by returning raw values", func(t *testing.T) {
		dough := map[string]float64{
			"semolina": 500,
			"water":    350,
		}

		result := toBakerPercentages(dough)

		amounts := map[string]float64{}
		for _, ing := range result {
			amounts[ing.Name] = ing.Amount
		}

		assert.Equal(t, 500.0, amounts["semolina"])
		assert.Equal(t, 350.0, amounts["water"])
	})
}

func TestToDomainRecipe_SetsReferenceArea(t *testing.T) {
	draft := GeneratedRecipeDraft{
		Name:        "Test",
		Description: "Test",
		Author:      "AI",
		Dough:       map[string]float64{"flour": 1000, "water": 700},
		Topping:     map[string]float64{"mozzarella": 300},
	}

	recipe := draft.toDomainRecipe()
	assert.Equal(t, defaultToppingReferenceArea, recipe.Topping.ReferenceArea)
	assert.Equal(t, 100.0, recipe.Dough.Ingredients[0].Amount) // flour = 100%
}

func intPtr(value int) *int {
	return &value
}
