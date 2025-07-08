package application_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/cfioretti/recipe-manager/internal/recipe-manager/application"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
)

type StubIngredientsBalancerClient struct {
	BalanceFunc func(ctx context.Context, recipe domain.Recipe, pans domain.Pans) (*domain.RecipeAggregate, error)
}

func (s *StubIngredientsBalancerClient) Balance(ctx context.Context, recipe domain.Recipe, pans domain.Pans) (*domain.RecipeAggregate, error) {
	return s.BalanceFunc(ctx, recipe, pans)
}

func (s *StubIngredientsBalancerClient) Close() error {
	return nil
}

func TestBalance(t *testing.T) {
	stubClient := createStubIngredientsBalancerClient()
	service := application.NewRemoteIngredientsBalancerService(stubClient)

	recipeID := 1
	recipeUUID := uuid.New()
	recipeName := "Test Pizza"
	recipeDescription := "A test pizza recipe"
	recipeAuthor := "Test Author"

	dough := domain.Dough{
		Name:             "Basic Dough",
		PercentVariation: 0.0,
		Ingredients: []domain.Ingredient{
			{Name: "Flour", Amount: 500.0},
			{Name: "Water", Amount: 350.0},
			{Name: "Salt", Amount: 10.0},
			{Name: "Yeast", Amount: 5.0},
		},
	}

	topping := domain.Topping{
		Name:          "Basic Topping",
		ReferenceArea: 615.75,
		Ingredients: []domain.Ingredient{
			{Name: "Tomato Sauce", Amount: 200.0},
			{Name: "Mozzarella", Amount: 300.0},
			{Name: "Basil", Amount: 10.0},
		},
	}

	steps := domain.Steps{
		RecipeId: recipeID,
		Steps: []domain.Step{
			{Id: 1, StepNumber: 1, Description: "Mix ingredients"},
			{Id: 2, StepNumber: 2, Description: "Knead dough"},
			{Id: 3, StepNumber: 3, Description: "Let rise"},
			{Id: 4, StepNumber: 4, Description: "Add toppings"},
			{Id: 5, StepNumber: 5, Description: "Bake"},
		},
	}

	recipe := domain.Recipe{
		Id:          recipeID,
		Uuid:        recipeUUID,
		Name:        recipeName,
		Description: recipeDescription,
		Author:      recipeAuthor,
		Dough:       dough,
		Topping:     topping,
		Steps:       steps,
	}

	diameter := 28
	pans := domain.Pans{
		Pans: []domain.Pan{
			{
				Shape: "round",
				Measures: domain.Measures{
					Diameter: &diameter,
				},
				Name: "round 28 cm",
				Area: 615.75,
			},
		},
		TotalArea: 615.75,
	}

	result, err := service.Balance(context.Background(), recipe, pans)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, recipe.Id, result.Id)
	assert.Equal(t, recipe.Uuid, result.Uuid)
	assert.Equal(t, recipe.Name, result.Name)
	assert.Equal(t, recipe.Description, result.Description)
	assert.Equal(t, recipe.Author, result.Author)

	assert.GreaterOrEqual(t, len(result.SplitIngredients.SplitDough), 1)
	assert.GreaterOrEqual(t, len(result.SplitIngredients.SplitTopping), 1)
}

func createStubIngredientsBalancerClient() *StubIngredientsBalancerClient {
	return &StubIngredientsBalancerClient{
		BalanceFunc: func(ctx context.Context, recipe domain.Recipe, pans domain.Pans) (*domain.RecipeAggregate, error) {
			result := &domain.RecipeAggregate{
				Recipe: recipe,
				SplitIngredients: domain.SplitIngredients{
					SplitDough: []domain.Dough{
						{
							Name:             "Split Basic Dough",
							PercentVariation: 0.0,
							Ingredients: []domain.Ingredient{
								{Name: "Flour", Amount: 500.0},
								{Name: "Water", Amount: 350.0},
								{Name: "Salt", Amount: 10.0},
								{Name: "Yeast", Amount: 5.0},
							},
						},
					},
					SplitTopping: []domain.Topping{
						{
							Name:          "Split Basic Topping",
							ReferenceArea: 615.75,
							Ingredients: []domain.Ingredient{
								{Name: "Tomato Sauce", Amount: 200.0},
								{Name: "Mozzarella", Amount: 300.0},
								{Name: "Basil", Amount: 10.0},
							},
						},
					},
				},
			}
			return result, nil
		},
	}
}
