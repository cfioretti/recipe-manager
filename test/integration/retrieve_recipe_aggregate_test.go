package integration

import (
	"context"
	"encoding/json"
	"testing"

	balancerapplication "github.com/cfioretti/recipe-manager/internal/ingredients-balancer/application"
	balancerdomain "github.com/cfioretti/recipe-manager/internal/ingredients-balancer/domain"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/application"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/infrastructure/mysql"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
		balancerapplication.NewDoughCalculatorService(),
		balancerapplication.NewDoughBalancerService(),
	)

	t.Run("Happy Path - retrieve RecipeAggregate successfully", func(t *testing.T) {
		dough := domain.Dough{PercentVariation: -10, Flour: 60, Water: 30, Salt: 5, EvoOil: 3, Yeast: 2}
		doughJSON, err := json.Marshal(dough)
		testRecipe := &domain.Recipe{
			Uuid:        uuid.New(),
			Name:        "Test Recipe",
			Description: "Test Recipe Description",
			Author:      "Test Author",
			Dough:       dough,
		}

		_, err = db.DB.Exec(`DELETE FROM recipes WHERE true`)
		query := `INSERT INTO recipes (uuid, name, description, author, dough) VALUES (?, ?, ?, ?, ?)`
		_, err = db.DB.Exec(query, testRecipe.Uuid, testRecipe.Name, testRecipe.Description, testRecipe.Author, string(doughJSON))
		if err != nil {
			t.Fatal(err)
		}

		pans := balancerdomain.Pans{
			Pans: []balancerdomain.Pan{
				{
					Shape: "round",
					Measures: balancerdomain.Measures{
						Diameter: intPtr(50),
					},
				},
				{
					Shape: "square",
					Measures: balancerdomain.Measures{
						Edge: intPtr(20),
					},
				},
				{
					Shape: "rectangular",
					Measures: balancerdomain.Measures{
						Width:  intPtr(30),
						Length: intPtr(40),
					},
				},
			},
		}

		expectedDough := domain.Dough{PercentVariation: 0, Flour: 962.1, Water: 481.1, Salt: 80.2, EvoOil: 48.1, Yeast: 32.1}

		splitDough1 := domain.Dough{PercentVariation: 0, Flour: 530.1, Water: 265.1, Salt: 44.2, EvoOil: 26.5, Yeast: 17.7}
		splitDough2 := domain.Dough{PercentVariation: 0, Flour: 108, Water: 54, Salt: 9, EvoOil: 5.4, Yeast: 3.6}
		splitDough3 := domain.Dough{PercentVariation: 0, Flour: 324, Water: 162, Salt: 27, EvoOil: 16.2, Yeast: 10.8}
		expectedSplitDough := []domain.Dough{splitDough1, splitDough2, splitDough3}

		result, err := service.Handle(testRecipe.Uuid, pans)

		assert.NoError(t, err)
		assert.Equal(t, testRecipe.Uuid, result.Uuid)
		assert.Equal(t, expectedDough, result.Dough)
		assert.Equal(t, expectedSplitDough, result.SplitIngredients.SplitDough)
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
		dough := domain.Dough{PercentVariation: -10, Flour: 60, Water: 30, Salt: 5, EvoOil: 3, Yeast: 2}
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
