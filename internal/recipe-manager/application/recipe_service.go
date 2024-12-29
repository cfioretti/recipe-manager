package application

import (
	"recipe-manager/internal/recipe-manager/domain"

	"github.com/google/uuid"
)

type RecipeRepository interface {
	GetRecipeByUuid(uuid.UUID) *domain.Recipe
}

type RecipeService struct {
	repository RecipeRepository
}

func NewRecipeService(repository RecipeRepository) *RecipeService {
	return &RecipeService{repository: repository}
}

func (rs *RecipeService) Handle(recipeUuid uuid.UUID) *domain.RecipeAggregate {
	recipe := rs.repository.GetRecipeByUuid(recipeUuid)
	return &domain.RecipeAggregate{Recipe: *recipe}
}
