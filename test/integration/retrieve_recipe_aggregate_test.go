package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/cfioretti/recipe-manager/internal/recipe-manager/application"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/infrastructure/mysql"
)

type StubCalculatorService struct {
	TotalDoughWeightByPansFunc func(ctx context.Context, pans domain.Pans) (*domain.Pans, error)
}

func (s *StubCalculatorService) TotalDoughWeightByPans(ctx context.Context, pans domain.Pans) (*domain.Pans, error) {
	return s.TotalDoughWeightByPansFunc(ctx, pans)
}

type StubBalancerService struct {
	BalanceFunc func(ctx context.Context, recipe domain.Recipe, pans domain.Pans) (*domain.RecipeAggregate, error)
}

func (s *StubBalancerService) Balance(ctx context.Context, recipe domain.Recipe, pans domain.Pans) (*domain.RecipeAggregate, error) {
	return s.BalanceFunc(ctx, recipe, pans)
}

func TestRecipeIntegration(t *testing.T) {
	ctx := context.Background()
	db, err := SetupTestDb(t)
	if err != nil {
		t.Fatal(err)
	}
	defer func(db *TestDatabase, ctx context.Context) {
		_ = db.Cleanup(ctx)
	}(db, ctx)

	stubCalculator := createStubCalculatorService()
	stubBalancer := createStubBalancerService()

	service := application.NewRecipeService(
		mysql.NewMySqlRecipeRepository(db.DB),
		stubCalculator,
		stubBalancer,
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
		topping := domain.Topping{
			ReferenceArea: 1200,
			Ingredients: []domain.Ingredient{
				{Name: "peeledTomatoes", Amount: 300},
				{Name: "mozzarellaCheese", Amount: 250},
				{Name: "basil", Amount: 15},
				{Name: "evoOil", Amount: 15},
				{Name: "parmesanCheese", Amount: 20},
			},
		}
		testRecipe := &domain.Recipe{
			Uuid:        uuid.New(),
			Name:        "Test Recipe",
			Description: "Test Recipe Description",
			Author:      "Test Author",
			Dough:       dough,
			Topping:     topping,
		}

		stringDoughJSON := `{"salt": 5, "flour": 60, "water": 30, "evoOil": 3, "yeast": 2, "percentVariation": -10}`
		stringToppingJSON := `{"referenceArea": 1200, "basil": 15, "evoOil": 15, "parmesanCheese": 20, "peeledTomatoes": 300, "mozzarellaCheese": 250}`
		_, err = db.DB.Exec(`DELETE FROM recipes WHERE true`)
		query := `INSERT INTO recipes (uuid, name, description, author, dough, topping) VALUES (?, ?, ?, ?, ?, ?)`
		_, err = db.DB.Exec(query, testRecipe.Uuid, testRecipe.Name, testRecipe.Description, testRecipe.Author, stringDoughJSON, stringToppingJSON)
		if err != nil {
			t.Fatal(err)
		}

		pans := domain.Pans{
			Pans: []domain.Pan{
				{
					Shape: "round",
					Measures: domain.Measures{
						Diameter: intPtr(50),
					},
				},
				{
					Shape: "square",
					Measures: domain.Measures{
						Edge: intPtr(20),
					},
				},
				{
					Shape: "rectangular",
					Measures: domain.Measures{
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
			Name:             "round 50 cm",
			PercentVariation: 0,
			Ingredients: []domain.Ingredient{
				{Name: "flour", Amount: 530.1},
				{Name: "water", Amount: 265.1},
				{Name: "salt", Amount: 44.2},
				{Name: "evoOil", Amount: 26.5},
				{Name: "yeast", Amount: 17.7},
			}}
		splitDough2 := domain.Dough{
			Name:             "square 20 cm",
			PercentVariation: 0,
			Ingredients: []domain.Ingredient{
				{Name: "flour", Amount: 108},
				{Name: "water", Amount: 54},
				{Name: "salt", Amount: 9},
				{Name: "evoOil", Amount: 5.4},
				{Name: "yeast", Amount: 3.6},
			}}
		splitDough3 := domain.Dough{
			Name:             "rectangular 30 x 40 cm",
			PercentVariation: 0,
			Ingredients: []domain.Ingredient{
				{Name: "flour", Amount: 324},
				{Name: "water", Amount: 162},
				{Name: "salt", Amount: 27},
				{Name: "evoOil", Amount: 16.2},
				{Name: "yeast", Amount: 10.8},
			}}
		expectedSplitDough := []domain.Dough{splitDough1, splitDough2, splitDough3}
		expectedTopping := domain.Topping{
			ReferenceArea: 1200,
			Ingredients: []domain.Ingredient{
				{Name: "peeledTomatoes", Amount: 890.9},
				{Name: "mozzarellaCheese", Amount: 742.4},
				{Name: "parmesanCheese", Amount: 59.4},
				{Name: "basil", Amount: 44.5},
				{Name: "evoOil", Amount: 44.5},
			},
		}

		result, err := service.Handle(ctx, testRecipe.Uuid, pans)

		assert.NoError(t, err)
		assert.Equal(t, testRecipe.Uuid, result.Uuid)
		assert.Equal(t, expectedDough, result.Recipe.Dough)
		assert.Equal(t, expectedSplitDough, result.SplitIngredients.SplitDough)
		assert.Equal(t, expectedTopping, result.Recipe.Topping)
	})

	t.Run("Error - Recipe not found in repository", func(t *testing.T) {
		nonExistentUuid := uuid.New()

		pans := domain.Pans{
			Pans: []domain.Pan{
				{
					Shape: "round",
					Measures: domain.Measures{
						Diameter: intPtr(50),
					},
				},
			},
		}

		result, err := service.Handle(ctx, nonExistentUuid, pans)

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
		doughJSON, _ := json.Marshal(dough)
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

		pans := domain.Pans{
			Pans: []domain.Pan{
				{
					Shape: "triangle",
					Measures: domain.Measures{
						Diameter: nil,
					},
				},
			},
		}

		result, err := service.Handle(ctx, testRecipe.Uuid, pans)

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

func createStubCalculatorService() *StubCalculatorService {
	return &StubCalculatorService{
		TotalDoughWeightByPansFunc: func(ctx context.Context, pans domain.Pans) (*domain.Pans, error) {
			for _, pan := range pans.Pans {
				if pan.Shape == "triangle" {
					return nil, fmt.Errorf("unsupported shape: %s", pan.Shape)
				}
			}

			totalArea := 0.0
			for _, pan := range pans.Pans {
				totalArea += pan.Area
			}
			pans.TotalArea = totalArea
			return &pans, nil
		},
	}
}

func createStubBalancerService() *StubBalancerService {
	return &StubBalancerService{
		BalanceFunc: func(ctx context.Context, recipe domain.Recipe, pans domain.Pans) (*domain.RecipeAggregate, error) {
			splitDough1 := domain.Dough{
				Name:             "round 50 cm",
				PercentVariation: 0,
				Ingredients: []domain.Ingredient{
					{Name: "flour", Amount: 530.1},
					{Name: "water", Amount: 265.1},
					{Name: "salt", Amount: 44.2},
					{Name: "evoOil", Amount: 26.5},
					{Name: "yeast", Amount: 17.7},
				}}
			splitDough2 := domain.Dough{
				Name:             "square 20 cm",
				PercentVariation: 0,
				Ingredients: []domain.Ingredient{
					{Name: "flour", Amount: 108},
					{Name: "water", Amount: 54},
					{Name: "salt", Amount: 9},
					{Name: "evoOil", Amount: 5.4},
					{Name: "yeast", Amount: 3.6},
				}}
			splitDough3 := domain.Dough{
				Name:             "rectangular 30 x 40 cm",
				PercentVariation: 0,
				Ingredients: []domain.Ingredient{
					{Name: "flour", Amount: 324},
					{Name: "water", Amount: 162},
					{Name: "salt", Amount: 27},
					{Name: "evoOil", Amount: 16.2},
					{Name: "yeast", Amount: 10.8},
				}}

			updatedRecipe := recipe
			updatedRecipe.Dough = domain.Dough{
				PercentVariation: -10,
				Ingredients: []domain.Ingredient{
					{Name: "flour", Amount: 962.1},
					{Name: "water", Amount: 481.1},
					{Name: "salt", Amount: 80.2},
					{Name: "evoOil", Amount: 48.1},
					{Name: "yeast", Amount: 32.1},
				},
			}
			updatedRecipe.Topping = domain.Topping{
				ReferenceArea: 1200,
				Ingredients: []domain.Ingredient{
					{Name: "peeledTomatoes", Amount: 890.9},
					{Name: "mozzarellaCheese", Amount: 742.4},
					{Name: "parmesanCheese", Amount: 59.4},
					{Name: "basil", Amount: 44.5},
					{Name: "evoOil", Amount: 44.5},
				},
			}

			return &domain.RecipeAggregate{
				Recipe: updatedRecipe,
				SplitIngredients: domain.SplitIngredients{
					SplitDough:   []domain.Dough{splitDough1, splitDough2, splitDough3},
					SplitTopping: []domain.Topping{updatedRecipe.Topping},
				},
			}, nil
		},
	}
}
