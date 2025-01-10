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
	t.Skip() // todo - add whole controller tests
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

	t.Run("retrieve RecipeAggregate successfully", func(t *testing.T) {
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

		result, err := service.Handle(testRecipe.Uuid, balancerdomain.Pans{})

		assert.NoError(t, err)
		assert.Equal(t, testRecipe.Name, result.Name)
		assert.Equal(t, testRecipe.Dough, result.Dough)
	})
}
