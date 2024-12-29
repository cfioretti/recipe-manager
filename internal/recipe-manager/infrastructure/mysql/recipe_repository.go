package mysql

import (
	"recipe-manager/internal/recipe-manager/domain"

	"github.com/google/uuid"
)

type Repository struct{}

func NewMysqlRecipeRepository() *Repository {
	return &Repository{}
}

func (rr Repository) GetRecipeByUuid(recipeUuid uuid.UUID) *domain.Recipe {
	return &domain.Recipe{
		Id:     1,
		Uuid:   recipeUuid,
		Name:   "Margherita",
		Author: "PizzaMaker",
	}
}
