package integration

import (
	"context"
	"testing"

	"recipe-manager/internal/recipe-manager/application"
	"recipe-manager/internal/recipe-manager/domain"
	"recipe-manager/internal/recipe-manager/infrastructure/mysql"

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

	service := application.NewRecipeService(mysql.NewMySqlRecipeRepository(db.DB))

	t.Run("retrieve RecipeAggregate successfully", func(t *testing.T) {
		testRecipe := &domain.Recipe{
			Uuid:   uuid.New(),
			Name:   "Test Recipe",
			Author: "Test Author",
		}

		_, err = db.DB.Exec(`DELETE FROM recipes WHERE true`)
		query := `INSERT INTO recipes (uuid, name, author) VALUES (?, ?, ?)`
		_, err = db.DB.Exec(query, testRecipe.Uuid, testRecipe.Name, testRecipe.Author)
		if err != nil {
			t.Fatal(err)
		}

		result, err := service.Handle(testRecipe.Uuid)

		assert.NoError(t, err)
		assert.Equal(t, testRecipe.Name, result.Name)
		assert.Equal(t, testRecipe.Author, result.Author)
	})

	t.Run("recipe not found returns empty RecipeAggregate", func(t *testing.T) {
		_, err = db.DB.Exec(`DELETE FROM recipes WHERE true`)
		if err != nil {
			t.Fatal(err)
		}

		result, err := service.Handle(uuid.New())

		assert.NoError(t, err)
		assert.Equal(t, &domain.RecipeAggregate{}, result)
	})
}
