package application

import (
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"

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

func (rs *RecipeService) Handle(recipeUuid uuid.UUID) (*domain.Recipe, error) {
	recipe, err := rs.repository.GetRecipeByUuid(recipeUuid)
	if err != nil {
		return nil, err
	}
	return recipe, nil
}
