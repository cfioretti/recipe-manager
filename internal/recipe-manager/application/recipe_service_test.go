package application

import (
	"testing"

	"recipe-manager/internal/recipe-manager/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRecipeRepository struct {
	mock.Mock
}

func (m *MockRecipeRepository) GetRecipe(recipeUuid uuid.UUID) *domain.RecipeAggregate {
	args := m.Called(recipeUuid)
	return args.Get(0).(*domain.RecipeAggregate)
}

func TestHandle(t *testing.T) {
	mockRepo := new(MockRecipeRepository)
	recipeUuid := uuid.New()

	t.Run("Success", func(t *testing.T) {
		expectedResponse := domain.RecipeAggregate{
			Recipe: domain.Recipe{RecipeUuid: recipeUuid},
		}
		mockRepo.On("GetRecipe", recipeUuid).Return(&expectedResponse)
		service := NewRecipeService(mockRepo)
		result := service.Handle(recipeUuid)

		assert.Equal(t, expectedResponse, *result)
	})
}
