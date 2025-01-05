package mysql

import (
	"database/sql"
	"encoding/json"

	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
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
		if errors.Is(err, sql.ErrNoRows) {
			return &domain.Recipe{}, nil
		}
		return nil, err
	}

	err = json.Unmarshal([]byte(doughJSON), &response.Dough)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
