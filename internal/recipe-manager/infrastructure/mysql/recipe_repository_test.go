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
		}

		rows := sqlmock.NewRows([]string{"id", "uuid", "name", "description", "author", "dough"}).
			AddRow(expectedRecipe.Id, expectedRecipe.Uuid, expectedRecipe.Name, expectedRecipe.Description,
				expectedRecipe.Author, doughJSON)
		mock.ExpectQuery(`SELECT id, uuid, name, description, author, dough FROM recipes WHERE uuid = ?`).
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
		mock.ExpectQuery(`SELECT id, uuid, name, description, author, dough FROM recipes WHERE uuid = ?`).
			WithArgs(newUuid).
			WillReturnError(sql.ErrNoRows)

		recipe, err := repo.GetRecipeByUuid(newUuid)

		assert.Error(t, err)
		assert.Nil(t, recipe)
	})

	t.Run("should return error on DB failure", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, uuid, name, description, author, dough FROM recipes WHERE uuid = ?`).
			WithArgs(newUuid).
			WillReturnError(sql.ErrConnDone)

		recipe, err := repo.GetRecipeByUuid(newUuid)

		assert.Error(t, err)
		assert.Nil(t, recipe)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
