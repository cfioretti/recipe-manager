package http

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRecipeService struct {
	mock.Mock
}

func (m *MockRecipeService) Handle(recipeUuid uuid.UUID, requestBody []byte) (*domain.RecipeAggregate, error) {
	args := m.Called(recipeUuid, requestBody)
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

		recipeAggregate := domain.RecipeAggregate{Recipe: domain.Recipe{Uuid: recipeUuid}}
		mockRecipeService := new(MockRecipeService)
		mockRecipeService.On("Handle", recipeUuid, mock.Anything).Return(&recipeAggregate, nil)

		controller := NewRecipeController(mockRecipeService)
		controller.RetrieveRecipeAggregate(ctx)

		assert.Equal(t, 200, ctx.Writer.Status())
	})

	t.Run("HTTP Status 400 on handler error", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = append(ctx.Params, gin.Param{Key: "uuid", Value: recipeUuid.String()})
		ctx.Request = httptest.NewRequest(
			http.MethodPost,
			"/recipes/"+recipeUuid.String()+"/aggregate",
			bytes.NewBuffer([]byte("body")),
		)

		recipeAggregate := domain.RecipeAggregate{Recipe: domain.Recipe{Uuid: recipeUuid}}
		mockRecipeService := new(MockRecipeService)
		mockRecipeService.On("Handle", recipeUuid, mock.Anything).Return(&recipeAggregate, errors.New("ERROR"))

		controller := NewRecipeController(mockRecipeService)
		controller.RetrieveRecipeAggregate(ctx)

		assert.Equal(t, 400, ctx.Writer.Status())
	})

	t.Run("Panic with wrong UUID", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = append(ctx.Params, gin.Param{Key: "uuid", Value: "WRONG UUID"})
		mockRecipeService := new(MockRecipeService)

		controller := NewRecipeController(mockRecipeService)

		assert.Panics(t, func() {
			controller.RetrieveRecipeAggregate(ctx)
		})
	})
}
