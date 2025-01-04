package http

import (
	"errors"
	"net/http/httptest"
	"testing"

	"recipe-manager/internal/recipe-manager/domain"

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

func TestRetrieveRecipeAggregate(t *testing.T) {
	recipeUuid := uuid.New()

	t.Run("HTTP Status 200 on Success", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = append(ctx.Params, gin.Param{Key: "uuid", Value: recipeUuid.String()})
		expectedResponse := domain.RecipeAggregate{
			Recipe: domain.Recipe{Uuid: recipeUuid},
		}
		mockService := new(MockRecipeService)
		mockService.On("Handle", recipeUuid).Return(&expectedResponse, nil)
		controller := NewRecipeController(mockService)

		controller.RetrieveRecipeAggregate(ctx)

		assert.Equal(t, 200, ctx.Writer.Status())
	})

	t.Run("HTTP Status 400 on Failure", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = append(ctx.Params, gin.Param{Key: "uuid", Value: recipeUuid.String()})
		expectedResponse := domain.RecipeAggregate{
			Recipe: domain.Recipe{Uuid: recipeUuid},
		}
		mockService := new(MockRecipeService)
		mockService.On("Handle", recipeUuid).Return(&expectedResponse, errors.New("ERROR"))
		controller := NewRecipeController(mockService)

		controller.RetrieveRecipeAggregate(ctx)

		assert.Equal(t, 500, ctx.Writer.Status())
	})

	t.Run("Panic with wrong UUID", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = append(ctx.Params, gin.Param{Key: "uuid", Value: "WRONG UUID"})
		mockService := new(MockRecipeService)

		controller := NewRecipeController(mockService)

		assert.Panics(t, func() {
			controller.RetrieveRecipeAggregate(ctx)
		})
	})
}
