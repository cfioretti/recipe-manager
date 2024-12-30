package mysql

import (
	"database/sql"

	"recipe-manager/internal/recipe-manager/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Repository struct {
	db *sql.DB
}

func NewMysqlRecipeRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (rr Repository) GetRecipeByUuid(recipeUuid uuid.UUID) (*domain.Recipe, error) {
	var response domain.Recipe

	query := `
		SELECT id, uuid, name, author
		FROM recipes
		WHERE uuid = ?
	`
	err := rr.db.QueryRow(query, recipeUuid).Scan(
		&response.Id,
		&response.Uuid,
		&response.Name,
		&response.Author,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &domain.Recipe{}, nil
		}
		return nil, err
	}

	return &response, nil
}
