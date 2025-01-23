package mysql

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
)

func TestGetRecipeByUuid(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMySqlRecipeRepository(db)
	newUuid := uuid.New()

	t.Run("should return recipe successfully when found", func(t *testing.T) {
		doughJSON := `{"salt": 5, "flour": 60, "water": 30, "evoOil": 3, "yeast": 2, "percentVariation": -10}`
		toppingJSON := `{"referenceArea": 1200, "mozzarellaCheese": 250, "tomatoPuree": 250, "basil": 10, "evoOil": 10, "parmesanCheese": 20}`
		expectedRecipe := &domain.Recipe{
			Id:          1,
			Uuid:        newUuid,
			Name:        "Test Recipe",
			Description: "Test Recipe Description",
			Author:      "Test Author",
			Dough: domain.Dough{
				PercentVariation: -10,
				Ingredients: []domain.Ingredient{
					{Name: "salt", Amount: 5},
					{Name: "flour", Amount: 60},
					{Name: "water", Amount: 30},
					{Name: "evoOil", Amount: 3},
					{Name: "yeast", Amount: 2},
				},
			},
			Topping: domain.Topping{
				Ingredients: []domain.Ingredient{
					{Name: "mozzarellaCheese", Amount: 250},
					{Name: "tomatoPuree", Amount: 250},
					{Name: "basil", Amount: 10},
					{Name: "evoOil", Amount: 10},
					{Name: "parmesanCheese", Amount: 20},
				},
			},
		}

		rows := sqlmock.NewRows([]string{"id", "uuid", "name", "description", "author", "dough", "topping"}).
			AddRow(expectedRecipe.Id, expectedRecipe.Uuid, expectedRecipe.Name, expectedRecipe.Description,
				expectedRecipe.Author, doughJSON, toppingJSON)
		mock.ExpectQuery(`SELECT id, uuid, name, description, author, dough, topping FROM recipes WHERE uuid = ?`).
			WithArgs(newUuid).
			WillReturnRows(rows)

		recipe, err := repo.GetRecipeByUuid(newUuid)

		assert.NoError(t, err)
		assert.Equal(t, expectedRecipe.Uuid, recipe.Uuid)
		assert.Equal(t, expectedRecipe.Name, recipe.Name)
		assert.Equal(t, expectedRecipe.Author, recipe.Author)
		assert.ElementsMatch(t, expectedRecipe.Dough.Ingredients, recipe.Dough.Ingredients)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when recipe is not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, uuid, name, description, author, dough, topping FROM recipes WHERE uuid = ?`).
			WithArgs(newUuid).
			WillReturnError(sql.ErrNoRows)

		recipe, err := repo.GetRecipeByUuid(newUuid)

		assert.Error(t, err)
		assert.Nil(t, recipe)
	})

	t.Run("should return error on DB failure", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, uuid, name, description, author, dough, topping FROM recipes WHERE uuid = ?`).
			WithArgs(newUuid).
			WillReturnError(sql.ErrConnDone)

		recipe, err := repo.GetRecipeByUuid(newUuid)

		assert.Error(t, err)
		assert.Nil(t, recipe)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestParseIngredients(t *testing.T) {
	tests := []struct {
		name          string
		jsonStr       string
		specialField  string
		expectedIngr  []domain.Ingredient
		expectedValue float64
		expectError   bool
	}{
		{
			name:         "Valid dough ingredients",
			jsonStr:      `{"salt": 5, "flour": 60, "water": 30, "evoOil": 3, "yeast": 2, "percentVariation": -10}`,
			specialField: "percentVariation",
			expectedIngr: []domain.Ingredient{
				{Name: "flour", Amount: 60},
				{Name: "water", Amount: 30},
				{Name: "salt", Amount: 5},
				{Name: "evoOil", Amount: 3},
				{Name: "yeast", Amount: 2},
			},
			expectedValue: -10,
			expectError:   false,
		},
		{
			name:         "Valid topping ingredients",
			jsonStr:      `{"referenceArea": 1200, "mozzarellaCheese": 250, "tomatoPuree": 250, "basil": 10, "evoOil": 10, "parmesanCheese": 20}`,
			specialField: "referenceArea",
			expectedIngr: []domain.Ingredient{
				{Name: "mozzarellaCheese", Amount: 250},
				{Name: "tomatoPuree", Amount: 250},
				{Name: "parmesanCheese", Amount: 20},
				{Name: "basil", Amount: 10},
				{Name: "evoOil", Amount: 10},
			},
			expectedValue: 1200,
			expectError:   false,
		},
		{
			name:         "Invalid JSON",
			jsonStr:      `{"salt": 5, "flour": "sixty"}`,
			specialField: "percentVariation",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ingredients, specialValue, err := parseIngredients(tt.jsonStr, tt.specialField)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.ElementsMatch(t, tt.expectedIngr, ingredients)
				assert.Equal(t, tt.expectedValue, specialValue)
			}
		})
	}
}
