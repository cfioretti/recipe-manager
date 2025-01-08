package mysql

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetRecipeByUuid(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMySqlRecipeRepository(db)
	newUuid := uuid.New()

	t.Run("should return recipe successfully when found", func(t *testing.T) {
		dough := domain.Dough{PercentVariation: -10, Flour: 60, Water: 30, Salt: 5, EvoOil: 3, Yeast: 2}
		expectedRecipe := &domain.Recipe{
			Id:          1,
			Uuid:        newUuid,
			Name:        "Test Recipe",
			Description: "Test Recipe Description",
			Author:      "Test Author",
			Dough:       dough,
		}
		doughJSON, err := json.Marshal(dough)
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "uuid", "name", "description", "author", "dough"}).
			AddRow(expectedRecipe.Id, expectedRecipe.Uuid, expectedRecipe.Name, expectedRecipe.Description,
				expectedRecipe.Author, string(doughJSON))
		mock.ExpectQuery(`SELECT id, uuid, name, description, author, dough FROM recipes WHERE uuid = ?`).
			WithArgs(newUuid).
			WillReturnRows(rows)

		recipe, err := repo.GetRecipeByUuid(newUuid)

		assert.NoError(t, err)
		assert.Equal(t, expectedRecipe, recipe)
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
