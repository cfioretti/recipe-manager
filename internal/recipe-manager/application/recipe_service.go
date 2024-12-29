package application

import (
	"recipe-manager/internal/recipe-manager/domain"

	"github.com/google/uuid"
)

type RecipeRepositoryInterface interface {
	GetRecipe(uuid.UUID) *domain.RecipeAggregate
}

type RecipeService struct {
	repository RecipeRepositoryInterface
}

func NewRecipeService(repository RecipeRepositoryInterface) *RecipeService {
	return &RecipeService{repository: repository}
}

func (s *RecipeService) Handle(recipeUuid uuid.UUID) *domain.RecipeAggregate {
	return s.repository.GetRecipe(recipeUuid)
}
