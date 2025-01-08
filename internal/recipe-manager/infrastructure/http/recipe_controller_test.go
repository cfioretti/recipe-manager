package http

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	calculatordomain "github.com/cfioretti/recipe-manager/internal/ingredients-balancer/domain"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRecipeService struct {
	mock.Mock
}

func (m *MockRecipeService) Handle(recipeUuid uuid.UUID) (*domain.Recipe, error) {
	args := m.Called(recipeUuid)
	return args.Get(0).(*domain.Recipe), args.Error(1)
}

type MockCalculatorService struct {
	mock.Mock
}

func (m *MockCalculatorService) TotalDoughWeightByPans(params []byte) (*calculatordomain.Pans, error) {
	args := m.Called(params)
	return args.Get(0).(*calculatordomain.Pans), args.Error(1)
}

type MockBalancerService struct {
	mock.Mock
}

func (m *MockBalancerService) Balance(recipe domain.Recipe, pans calculatordomain.Pans) (*domain.RecipeAggregate, error) {
	args := m.Called(recipe, pans)
	return args.Get(0).(*domain.RecipeAggregate), args.Error(1)
}

func TestRetrieveRecipeAggregate(t *testing.T) {
	recipeUuid := uuid.New()

	t.Run("HTTP Status 200 on success", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = append(ctx.Params, gin.Param{Key: "uuid", Value: recipeUuid.String()})
		ctx.Request = httptest.NewRequest(
			http.MethodPost,
			"/recipes/"+recipeUuid.String()+"/aggregate",
			bytes.NewBuffer([]byte("body")),
		)

		recipe := domain.Recipe{Uuid: recipeUuid}
		mockRecipeService := new(MockRecipeService)
		mockRecipeService.On("Handle", recipeUuid).Return(&recipe, nil)
		pans := calculatordomain.Pans{}
		mockCalculatorService := new(MockCalculatorService)
		mockCalculatorService.On("TotalDoughWeightByPans", mock.Anything).Return(&pans, nil)
		recipeAggregate := domain.RecipeAggregate{Recipe: recipe}
		mockBalancerService := new(MockBalancerService)
		mockBalancerService.On("Balance", mock.Anything, mock.Anything).Return(&recipeAggregate, nil)

		controller := NewRecipeController(mockRecipeService, mockCalculatorService, mockBalancerService)
		controller.RetrieveRecipeAggregate(ctx)

		assert.Equal(t, 200, ctx.Writer.Status())
	})

	t.Run("HTTP Status 400 on validation error", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = append(ctx.Params, gin.Param{Key: "uuid", Value: recipeUuid.String()})
		ctx.Request = httptest.NewRequest(
			http.MethodPost,
			"/recipes/"+recipeUuid.String()+"/aggregate",
			bytes.NewBuffer([]byte("WRONG BODY")),
		)

		recipe := domain.Recipe{Uuid: recipeUuid}
		mockRecipeService := new(MockRecipeService)
		mockRecipeService.On("Handle", recipeUuid).Return(&recipe, nil)
		pans := calculatordomain.Pans{}
		mockCalculatorService := new(MockCalculatorService)
		mockCalculatorService.On("TotalDoughWeightByPans", mock.Anything).Return(&pans, errors.New("ERROR"))
		recipeAggregate := domain.RecipeAggregate{Recipe: recipe}
		mockBalancerService := new(MockBalancerService)
		mockBalancerService.On("Balance", mock.Anything, mock.Anything).Return(&recipeAggregate, nil)

		controller := NewRecipeController(mockRecipeService, mockCalculatorService, mockBalancerService)
		controller.RetrieveRecipeAggregate(ctx)

		assert.Equal(t, 400, ctx.Writer.Status())
	})

	t.Run("HTTP Status 400 on handler error", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = append(ctx.Params, gin.Param{Key: "uuid", Value: recipeUuid.String()})
		ctx.Request = httptest.NewRequest(
			http.MethodPost,
			"/recipes/"+recipeUuid.String()+"/aggregate",
			bytes.NewBuffer([]byte("body")),
		)

		recipe := domain.Recipe{Uuid: recipeUuid}
		mockRecipeService := new(MockRecipeService)
		mockRecipeService.On("Handle", recipeUuid).Return(&recipe, errors.New("ERROR"))
		pans := calculatordomain.Pans{}
		mockCalculatorService := new(MockCalculatorService)
		mockCalculatorService.On("TotalDoughWeightByPans", mock.Anything).Return(&pans, nil)
		recipeAggregate := domain.RecipeAggregate{Recipe: recipe}
		mockBalancerService := new(MockBalancerService)
		mockBalancerService.On("Balance", mock.Anything, mock.Anything).Return(&recipeAggregate, nil)

		controller := NewRecipeController(mockRecipeService, mockCalculatorService, mockBalancerService)
		controller.RetrieveRecipeAggregate(ctx)

		assert.Equal(t, 400, ctx.Writer.Status())
	})

	t.Run("Panic with wrong UUID", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = append(ctx.Params, gin.Param{Key: "uuid", Value: "WRONG UUID"})
		mockRecipeService := new(MockRecipeService)
		mockCalculatorService := new(MockCalculatorService)
		mockBalancerService := new(MockBalancerService)

		controller := NewRecipeController(mockRecipeService, mockCalculatorService, mockBalancerService)

		assert.Panics(t, func() {
			controller.RetrieveRecipeAggregate(ctx)
		})
	})
}
