package repository

import (
	"recipe-manager/internal/recipe-manager/domain"

	"github.com/google/uuid"
)

type RecipeRepository struct{}

func NewRecipeRepository() *RecipeRepository {
	return &RecipeRepository{}
}

func (rr RecipeRepository) GetRecipe(recipeUuid uuid.UUID) *domain.RecipeAggregate {
	recipe := domain.Recipe{
		RecipeUuid: recipeUuid,
	}
	return &domain.RecipeAggregate{Recipe: recipe}
}
