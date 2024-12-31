package mysql

import (
	"database/sql"
	"testing"

	"recipe-manager/internal/recipe-manager/domain"

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
		expectedRecipe := &domain.Recipe{
			Id:     1,
			Uuid:   newUuid,
			Name:   "Test Recipe",
			Author: "Test Author",
		}
		rows := sqlmock.NewRows([]string{"id", "uuid", "name", "author"}).
			AddRow(expectedRecipe.Id, expectedRecipe.Uuid, expectedRecipe.Name, expectedRecipe.Author)
		mock.ExpectQuery("SELECT (.+) FROM recipes").
			WithArgs(newUuid).
			WillReturnRows(rows)

		recipe, err := repo.GetRecipeByUuid(newUuid)

		assert.NoError(t, err)
		assert.Equal(t, expectedRecipe, recipe)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return empty recipe when not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM recipes").
			WithArgs(newUuid).
			WillReturnError(sql.ErrNoRows)

		recipe, err := repo.GetRecipeByUuid(newUuid)

		assert.NoError(t, err)
		assert.Equal(t, &domain.Recipe{}, recipe)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error on DB failure", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM recipes").
			WithArgs(newUuid).
			WillReturnError(sql.ErrConnDone)

		recipe, err := repo.GetRecipeByUuid(newUuid)

		assert.Error(t, err)
		assert.Nil(t, recipe)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
