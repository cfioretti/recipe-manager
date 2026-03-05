package http

import (
	"bytes"
	"context"
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/cfioretti/recipe-manager/internal/recipe-manager/application"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
)

type MockGenerateRecipeService struct {
	mock.Mock
}

func (m *MockGenerateRecipeService) Handle(ctx context.Context, command application.GenerateRecipeCommand) (*domain.RecipeAggregate, error) {
	args := m.Called(ctx, command)
	result, _ := args.Get(0).(*domain.RecipeAggregate)
	return result, args.Error(1)
}

func TestGenerateRecipeAggregate(t *testing.T) {
	t.Run("HTTP Status 200 on success", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		body := `{
			"mode":"random",
			"pans":[{"shape":"round","measures":{"diameter":"100"}}]
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/recipes/generate", bytes.NewBuffer([]byte(body)))
		ctx.Request.Header.Set("Content-Type", "application/json")

		aggregate := &domain.RecipeAggregate{
			Recipe: domain.Recipe{Name: "Generated Pizza"},
		}
		mockService := new(MockGenerateRecipeService)
		mockService.On("Handle", mock.Anything, mock.Anything).Return(aggregate, nil)

		handler := NewGenerateRecipeHandler(mockService)
		handler.GenerateRecipeAggregate(ctx)

		assert.Equal(t, http.StatusOK, ctx.Writer.Status())
	})

	t.Run("HTTP Status 400 on validation error", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request = httptest.NewRequest(http.MethodPost, "/recipes/generate", bytes.NewBuffer([]byte(`{}`)))
		ctx.Request.Header.Set("Content-Type", "application/json")

		mockService := new(MockGenerateRecipeService)
		handler := NewGenerateRecipeHandler(mockService)
		handler.GenerateRecipeAggregate(ctx)

		assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
	})

	t.Run("HTTP Status 400 on service error", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		body := `{
			"mode":"prompt",
			"prompt":"crispy pizza",
			"pans":[{"shape":"round","measures":{"diameter":"100"}}]
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/recipes/generate", bytes.NewBuffer([]byte(body)))
		ctx.Request.Header.Set("Content-Type", "application/json")

		mockService := new(MockGenerateRecipeService)
		mockService.On("Handle", mock.Anything, mock.Anything).Return((*domain.RecipeAggregate)(nil), errors.New("generation error"))

		handler := NewGenerateRecipeHandler(mockService)
		handler.GenerateRecipeAggregate(ctx)

		assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
	})

	t.Run("HTTP Status 502 when payload is not serializable", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		body := `{
			"mode":"random",
			"pans":[{"shape":"round","measures":{"diameter":"100"}}]
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/recipes/generate", bytes.NewBuffer([]byte(body)))
		ctx.Request.Header.Set("Content-Type", "application/json")

		aggregate := &domain.RecipeAggregate{
			Recipe: domain.Recipe{
				Name: "Generated Pizza",
				Dough: domain.Dough{
					Ingredients: []domain.Ingredient{
						{Name: "flour", Amount: math.NaN()},
					},
				},
			},
		}
		mockService := new(MockGenerateRecipeService)
		mockService.On("Handle", mock.Anything, mock.Anything).Return(aggregate, nil)

		handler := NewGenerateRecipeHandler(mockService)
		handler.GenerateRecipeAggregate(ctx)

		assert.Equal(t, http.StatusBadGateway, ctx.Writer.Status())
		assert.Contains(t, recorder.Body.String(), "failed to serialize response payload")
	})
}
