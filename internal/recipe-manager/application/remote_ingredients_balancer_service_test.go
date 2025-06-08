package application_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	bdomain "github.com/cfioretti/recipe-manager/internal/ingredients-balancer/domain"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/application"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/infrastructure/grpc/client"
)

func TestBalance(t *testing.T) {
	t.Skip("this test requires a running gRPC server")

	grpcClient, err := client.NewIngredientsBalancerClient("localhost:50052", 5*time.Second)
	require.NoError(t, err)
	defer grpcClient.Close()

	service := application.NewRemoteIngredientsBalancerService(grpcClient)

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
	pans := bdomain.Pans{
		Pans: []bdomain.Pan{
			{
				Shape: "round",
				Measures: bdomain.Measures{
					Diameter: &diameter,
				},
				Name: "round 28 cm",
				Area: 615.75,
			},
		},
		TotalArea: 615.75,
	}

	result, err := service.Balance(recipe, pans)

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
