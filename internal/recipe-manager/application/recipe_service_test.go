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

func (m *MockRecipeRepository) GetRecipeByUuid(recipeUuid uuid.UUID) (*domain.Recipe, error) {
	args := m.Called(recipeUuid)
	return args.Get(0).(*domain.Recipe), nil
}

func TestHandle(t *testing.T) {
	mockRepo := new(MockRecipeRepository)
	recipeUuid := uuid.New()

	t.Run("Success", func(t *testing.T) {
		recipe := domain.Recipe{Uuid: recipeUuid}
		mockRepo.On("GetRecipeByUuid", recipeUuid).Return(&recipe)

		service := NewRecipeService(mockRepo)
		result := service.Handle(recipeUuid)

		expectedResponse := domain.RecipeAggregate{
			Recipe: recipe,
		}
		assert.Equal(t, expectedResponse, *result)
	})
}
