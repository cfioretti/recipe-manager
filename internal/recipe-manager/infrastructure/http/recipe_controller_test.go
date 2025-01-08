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

func (m *MockRecipeService) Handle(recipeUuid uuid.UUID) (*domain.RecipeAggregate, error) {
	args := m.Called(recipeUuid)
	return args.Get(0).(*domain.RecipeAggregate), args.Error(1)
}

type MockCalculatorService struct {
	mock.Mock
}

func (m *MockCalculatorService) TotalDoughWeightByPans(params []byte) (*calculatordomain.Pans, error) {
	args := m.Called(params)
	return args.Get(0).(*calculatordomain.Pans), args.Error(1)
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

		expectedResponse := domain.RecipeAggregate{
			Recipe: domain.Recipe{Uuid: recipeUuid},
		}
		mockRecipeService := new(MockRecipeService)
		mockRecipeService.On("Handle", recipeUuid).Return(&expectedResponse, nil)
		pans := calculatordomain.Pans{}
		mockCalculatorService := new(MockCalculatorService)
		mockCalculatorService.On("TotalDoughWeightByPans", mock.Anything).Return(&pans, nil)

		controller := NewRecipeController(mockRecipeService, mockCalculatorService)
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

		expectedResponse := domain.RecipeAggregate{
			Recipe: domain.Recipe{Uuid: recipeUuid},
		}
		mockRecipeService := new(MockRecipeService)
		mockRecipeService.On("Handle", recipeUuid).Return(&expectedResponse, nil)
		pans := calculatordomain.Pans{}
		mockCalculatorService := new(MockCalculatorService)
		mockCalculatorService.On("TotalDoughWeightByPans", mock.Anything).Return(&pans, errors.New("ERROR"))

		controller := NewRecipeController(mockRecipeService, mockCalculatorService)
		controller.RetrieveRecipeAggregate(ctx)

		assert.Equal(t, 400, ctx.Writer.Status())
	})

	t.Run("HTTP Status 500 on dB failure", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = append(ctx.Params, gin.Param{Key: "uuid", Value: recipeUuid.String()})
		ctx.Request = httptest.NewRequest(
			http.MethodPost,
			"/recipes/"+recipeUuid.String()+"/aggregate",
			bytes.NewBuffer([]byte("body")),
		)

		expectedResponse := domain.RecipeAggregate{
			Recipe: domain.Recipe{Uuid: recipeUuid},
		}
		mockRecipeService := new(MockRecipeService)
		mockRecipeService.On("Handle", recipeUuid).Return(&expectedResponse, errors.New("ERROR"))
		pans := calculatordomain.Pans{}
		mockCalculatorService := new(MockCalculatorService)
		mockCalculatorService.On("TotalDoughWeightByPans", mock.Anything).Return(&pans, nil)

		controller := NewRecipeController(mockRecipeService, mockCalculatorService)
		controller.RetrieveRecipeAggregate(ctx)

		assert.Equal(t, 500, ctx.Writer.Status())
	})

	t.Run("Panic with wrong UUID", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = append(ctx.Params, gin.Param{Key: "uuid", Value: "WRONG UUID"})
		mockRecipeService := new(MockRecipeService)
		mockCalculatorService := new(MockCalculatorService)

		controller := NewRecipeController(mockRecipeService, mockCalculatorService)

		assert.Panics(t, func() {
			controller.RetrieveRecipeAggregate(ctx)
		})
	})
}
