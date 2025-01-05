package application

import (
	"errors"
	"testing"

	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRecipeRepository struct {
	mock.Mock
}

func (m *MockRecipeRepository) GetRecipeByUuid(recipeUuid uuid.UUID) (*domain.Recipe, error) {
	args := m.Called(recipeUuid)
	return args.Get(0).(*domain.Recipe), args.Error(1)
}

func TestHandle(t *testing.T) {
	recipeUuid := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockRecipeRepository)
		recipe := domain.Recipe{Uuid: recipeUuid}
		mockRepo.On("GetRecipeByUuid", recipeUuid).Return(&recipe, nil)

		service := NewRecipeService(mockRepo)
		result, _ := service.Handle(recipeUuid)

		expectedResponse := domain.RecipeAggregate{
			Recipe: recipe,
		}
		assert.Equal(t, expectedResponse, *result)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := new(MockRecipeRepository)
		mockRepo.On("GetRecipeByUuid", mock.Anything).Return(&domain.Recipe{}, errors.New("DB ERROR"))

		service := NewRecipeService(mockRepo)
		_, err := service.Handle(recipeUuid)

		assert.Error(t, err)
	})
}
