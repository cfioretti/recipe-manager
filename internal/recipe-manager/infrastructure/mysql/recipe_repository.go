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
	var doughJSON, toppingJSON string

	query := `SELECT id, uuid, name, description, author, dough, topping FROM recipes WHERE uuid = ?`
	err := rr.db.QueryRow(query, recipeUuid).Scan(
		&response.Id,
		&response.Uuid,
		&response.Name,
		&response.Description,
		&response.Author,
		&doughJSON,
		&toppingJSON,
	)
	if err != nil {
		return nil, err
	}

	response.Dough, err = parseDough(doughJSON)
	if err != nil {
		return nil, err
	}

	response.Topping, err = parseTopping(toppingJSON)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func parseDough(doughJSON string) (domain.Dough, error) {
	ingredients, percentVariation, err := parseIngredients(doughJSON, "percentVariation")
	if err != nil {
		return domain.Dough{}, err
	}
	return domain.Dough{
		PercentVariation: percentVariation,
		Ingredients:      ingredients,
	}, nil
}

func parseTopping(toppingJSON string) (domain.Topping, error) {
	ingredients, referenceArea, err := parseIngredients(toppingJSON, "referenceArea")
	if err != nil {
		return domain.Topping{}, err
	}
	return domain.Topping{
		ReferenceArea: referenceArea,
		Ingredients:   ingredients,
	}, nil
}

func parseIngredients(jsonStr string, specialField string) ([]domain.Ingredient, float64, error) {
	var dataMap map[string]float64
	err := json.Unmarshal([]byte(jsonStr), &dataMap)
	if err != nil {
		return nil, 0, err
	}

	specialValue := dataMap[specialField]
	delete(dataMap, specialField)

	var ingredients []domain.Ingredient
	for key, value := range dataMap {
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

	return ingredients, specialValue, nil
}
