package integration

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	balancerapplication "github.com/cfioretti/recipe-manager/internal/ingredients-balancer/application"
	balancerdomain "github.com/cfioretti/recipe-manager/internal/ingredients-balancer/domain"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/application"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/infrastructure/mysql"
)

func TestRecipeIntegration(t *testing.T) {
	ctx := context.Background()
	db, err := SetupTestDb(t)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Cleanup(ctx)

	service := application.NewRecipeService(
		mysql.NewMySqlRecipeRepository(db.DB),
		balancerapplication.NewIngredientsCalculatorService(),
		balancerapplication.NewIngredientsBalancerService(),
	)

	t.Run("Happy Path - retrieve RecipeAggregate successfully", func(t *testing.T) {
		dough := domain.Dough{
			PercentVariation: -10,
			Ingredients: []domain.Ingredient{
				{Name: "flour", Amount: 60},
				{Name: "water", Amount: 30},
				{Name: "salt", Amount: 5},
				{Name: "evoOil", Amount: 3},
				{Name: "yeast", Amount: 2},
			},
		}
		testRecipe := &domain.Recipe{
			Uuid:        uuid.New(),
			Name:        "Test Recipe",
			Description: "Test Recipe Description",
			Author:      "Test Author",
			Dough:       dough,
		}

		stringDoughJSON := `{"salt": 5, "flour": 60, "water": 30, "evoOil": 3, "yeast": 2, "percentVariation": -10}`
		_, err = db.DB.Exec(`DELETE FROM recipes WHERE true`)
		query := `INSERT INTO recipes (uuid, name, description, author, dough) VALUES (?, ?, ?, ?, ?)`
		_, err = db.DB.Exec(query, testRecipe.Uuid, testRecipe.Name, testRecipe.Description, testRecipe.Author, stringDoughJSON)
		if err != nil {
			t.Fatal(err)
		}

		shape1 := "round"
		shape2 := "square"
		shape3 := "rectangular"
		pans := balancerdomain.Pans{
			Pans: []balancerdomain.Pan{
				{
					Shape: shape1,
					Measures: balancerdomain.Measures{
						Diameter: intPtr(50),
					},
				},
				{
					Shape: shape2,
					Measures: balancerdomain.Measures{
						Edge: intPtr(20),
					},
				},
				{
					Shape: shape3,
					Measures: balancerdomain.Measures{
						Width:  intPtr(30),
						Length: intPtr(40),
					},
				},
			},
		}

		expectedDough := domain.Dough{
			PercentVariation: -10,
			Ingredients: []domain.Ingredient{
				{Name: "flour", Amount: 962.1},
				{Name: "water", Amount: 481.1},
				{Name: "salt", Amount: 80.2},
				{Name: "evoOil", Amount: 48.1},
				{Name: "yeast", Amount: 32.1},
			},
		}

		splitDough1 := domain.Dough{
			Name:             shape1,
			PercentVariation: 0,
			Ingredients: []domain.Ingredient{
				{Name: "flour", Amount: 530.1},
				{Name: "water", Amount: 265.1},
				{Name: "salt", Amount: 44.2},
				{Name: "evoOil", Amount: 26.5},
				{Name: "yeast", Amount: 17.7},
			}}
		splitDough2 := domain.Dough{
			Name:             shape2,
			PercentVariation: 0,
			Ingredients: []domain.Ingredient{
				{Name: "flour", Amount: 108},
				{Name: "water", Amount: 54},
				{Name: "salt", Amount: 9},
				{Name: "evoOil", Amount: 5.4},
				{Name: "yeast", Amount: 3.6},
			}}
		splitDough3 := domain.Dough{
			Name:             shape3,
			PercentVariation: 0,
			Ingredients: []domain.Ingredient{
				{Name: "flour", Amount: 324},
				{Name: "water", Amount: 162},
				{Name: "salt", Amount: 27},
				{Name: "evoOil", Amount: 16.2},
				{Name: "yeast", Amount: 10.8},
			}}
		expectedSplitDough := []domain.Dough{splitDough1, splitDough2, splitDough3}

		result, err := service.Handle(testRecipe.Uuid, pans)

		assert.NoError(t, err)
		assert.Equal(t, testRecipe.Uuid, result.Uuid)
		assert.ElementsMatch(t, expectedDough.Ingredients, result.Recipe.Dough.Ingredients)
		for i, d := range result.SplitIngredients.SplitDough {
			assert.ElementsMatch(t, expectedSplitDough[i].Ingredients, d.Ingredients)
		}
	})

	t.Run("Error - Recipe not found in repository", func(t *testing.T) {
		nonExistentUuid := uuid.New()

		pans := balancerdomain.Pans{
			Pans: []balancerdomain.Pan{
				{
					Shape: "round",
					Measures: balancerdomain.Measures{
						Diameter: intPtr(50),
					},
				},
			},
		}

		result, err := service.Handle(nonExistentUuid, pans)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("Error - Invalid pans data -> Unsupported shape", func(t *testing.T) {
		dough := domain.Dough{
			PercentVariation: -10,
			Ingredients: []domain.Ingredient{
				{Name: "flour", Amount: 60},
				{Name: "water", Amount: 30},
				{Name: "salt", Amount: 5},
				{Name: "evoOil", Amount: 3},
				{Name: "yeast", Amount: 2},
			},
		}
		doughJSON, err := json.Marshal(dough)
		testRecipe := &domain.Recipe{
			Uuid:        uuid.New(),
			Name:        "Test Recipe",
			Description: "Test Recipe Description",
			Author:      "Test Author",
			Dough:       dough,
		}

		query := `INSERT INTO recipes (uuid, name, description, author, dough) VALUES (?, ?, ?, ?, ?)`
		_, err = db.DB.Exec(query, testRecipe.Uuid, testRecipe.Name, testRecipe.Description, testRecipe.Author, string(doughJSON))
		if err != nil {
			t.Fatal(err)
		}

		pans := balancerdomain.Pans{
			Pans: []balancerdomain.Pan{
				{
					Shape: "triangle",
					Measures: balancerdomain.Measures{
						Diameter: nil,
					},
				},
			},
		}

		result, err := service.Handle(testRecipe.Uuid, pans)

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}
