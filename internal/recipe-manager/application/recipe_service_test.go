package application

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
)

type MockRecipeRepository struct {
	mock.Mock
}

func (m *MockRecipeRepository) GetRecipeByUuid(recipeUuid uuid.UUID) (*domain.Recipe, error) {
	args := m.Called(recipeUuid)
	return args.Get(0).(*domain.Recipe), args.Error(1)
}

type MockCalculatorService struct {
	mock.Mock
}

func (m *MockCalculatorService) TotalDoughWeightByPans(params domain.Pans) (*domain.Pans, error) {
	args := m.Called(params)
	return args.Get(0).(*domain.Pans), args.Error(1)
}

type MockBalancerService struct {
	mock.Mock
}

func (m *MockBalancerService) Balance(recipe domain.Recipe, pans domain.Pans) (*domain.RecipeAggregate, error) {
	args := m.Called(recipe, pans)
	return args.Get(0).(*domain.RecipeAggregate), args.Error(1)
}

func TestHandle(t *testing.T) {
	recipeUuid := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockRecipeRepository := new(MockRecipeRepository)
		recipe := domain.Recipe{Uuid: recipeUuid}
		mockRecipeRepository.On("GetRecipeByUuid", recipeUuid).Return(&recipe, nil)
		pans := domain.Pans{}
		mockCalculatorService := new(MockCalculatorService)
		mockCalculatorService.On("TotalDoughWeightByPans", mock.Anything).Return(&pans, nil)
		recipeAggregate := domain.RecipeAggregate{Recipe: recipe}
		mockBalancerService := new(MockBalancerService)
		mockBalancerService.On("Balance", recipe, pans).Return(&recipeAggregate, nil)

		service := NewRecipeService(mockRecipeRepository, mockCalculatorService, mockBalancerService)
		result, _ := service.Handle(recipeUuid, pans)

		assert.Equal(t, recipeAggregate, *result)
	})

	t.Run("calculator service error", func(t *testing.T) {
		mockCalculatorService := new(MockCalculatorService)
		calculatorError := errors.New("calculator error")
		mockCalculatorService.On("TotalDoughWeightByPans", mock.Anything).Return((*domain.Pans)(nil), calculatorError)

		service := NewRecipeService(new(MockRecipeRepository), mockCalculatorService, new(MockBalancerService))
		result, err := service.Handle(recipeUuid, domain.Pans{})

		assert.Nil(t, result)
		assert.Equal(t, calculatorError, err)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRecipeRepository := new(MockRecipeRepository)
		repositoryError := errors.New("repository error")
		mockRecipeRepository.On("GetRecipeByUuid", recipeUuid).Return((*domain.Recipe)(nil), repositoryError)

		mockCalculatorService := new(MockCalculatorService)
		pans := domain.Pans{}
		mockCalculatorService.On("TotalDoughWeightByPans", mock.Anything).Return(&pans, nil)

		service := NewRecipeService(mockRecipeRepository, mockCalculatorService, new(MockBalancerService))
		result, err := service.Handle(recipeUuid, pans)

		assert.Nil(t, result)
		assert.Equal(t, repositoryError, err)
	})

	t.Run("balancer service error", func(t *testing.T) {
		mockRecipeRepository := new(MockRecipeRepository)
		recipe := domain.Recipe{Uuid: recipeUuid}
		mockRecipeRepository.On("GetRecipeByUuid", recipeUuid).Return(&recipe, nil)

		mockCalculatorService := new(MockCalculatorService)
		pans := domain.Pans{}
		mockCalculatorService.On("TotalDoughWeightByPans", mock.Anything).Return(&pans, nil)

		mockBalancerService := new(MockBalancerService)
		balancerError := errors.New("balancer error")
		mockBalancerService.On("Balance", mock.Anything, mock.Anything).Return((*domain.RecipeAggregate)(nil), balancerError)

		service := NewRecipeService(mockRecipeRepository, mockCalculatorService, mockBalancerService)
		result, err := service.Handle(recipeUuid, pans)

		assert.Nil(t, result)
		assert.Equal(t, balancerError, err)
	})
}
