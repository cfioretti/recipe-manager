package application

import (
	"recipe-manager/internal/recipe-manager/domain"

	"github.com/google/uuid"
)

type RecipeRepository interface {
	GetRecipeByUuid(uuid.UUID) (*domain.Recipe, error)
}

type RecipeService struct {
	repository RecipeRepository
}

func NewRecipeService(repository RecipeRepository) *RecipeService {
	return &RecipeService{repository: repository}
}

func (rs *RecipeService) Handle(recipeUuid uuid.UUID) *domain.RecipeAggregate {
	recipe, err := rs.repository.GetRecipeByUuid(recipeUuid)
	if err != nil {
		return nil
	}
	return &domain.RecipeAggregate{Recipe: *recipe}
}
