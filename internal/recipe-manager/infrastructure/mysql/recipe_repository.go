package mysql

import (
	"database/sql"
	"encoding/json"
	"slices"

	"github.com/google/uuid"

	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
)

type MySqlRecipeRepository struct {
	db *sql.DB
}

func NewMySqlRecipeRepository(db *sql.DB) *MySqlRecipeRepository {
	return &MySqlRecipeRepository{db: db}
}

func (rr MySqlRecipeRepository) GetRecipeByUuid(recipeUuid uuid.UUID) (*domain.Recipe, error) {
	var response domain.Recipe
	var doughJSON string

	query := `SELECT id, uuid, name, description, author, dough FROM recipes WHERE uuid = ?`
	err := rr.db.QueryRow(query, recipeUuid).Scan(
		&response.Id,
		&response.Uuid,
		&response.Name,
		&response.Description,
		&response.Author,
		&doughJSON,
	)
	if err != nil {
		return nil, err
	}

	var doughMap map[string]float64
	err = json.Unmarshal([]byte(doughJSON), &doughMap)
	if err != nil {
		return nil, err
	}

	percentVariation := doughMap["percentVariation"]
	delete(doughMap, "percentVariation")

	var ingredients []domain.Ingredient
	for key, value := range doughMap {
		ingredients = append(ingredients, domain.Ingredient{
			Name:   key,
			Amount: value,
		})
	}

	slices.SortFunc(ingredients, func(a, b domain.Ingredient) int {
		switch {
		case a.Amount > b.Amount:
			return -1
		case a.Amount == b.Amount:
			if a.Name < b.Name {
				return -1
			}
			fallthrough
		default:
			return 1
		}
	})

	response.Dough = domain.Dough{
		PercentVariation: percentVariation,
		Ingredients:      ingredients,
	}

	return &response, nil
}
