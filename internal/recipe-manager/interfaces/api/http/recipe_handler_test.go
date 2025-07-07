package http

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
)

type MockRecipeService struct {
	mock.Mock
}

func (m *MockRecipeService) Handle(recipeUuid uuid.UUID, requestBody domain.Pans) (*domain.RecipeAggregate, error) {
	args := m.Called(recipeUuid, requestBody)
	return args.Get(0).(*domain.RecipeAggregate), args.Error(1)
}

func TestRetrieveRecipeAggregate(t *testing.T) {
	recipeUuid := uuid.New()

	t.Run("HTTP Status 200 on success", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = append(ctx.Params, gin.Param{Key: "uuid", Value: recipeUuid.String()})
		body := `{"pans": [{"shape": "round","measures": {"diameter": "100"}}]}`
		ctx.Request = httptest.NewRequest(
			http.MethodPost,
			"/recipes/"+recipeUuid.String()+"/aggregate",
			bytes.NewBuffer([]byte(body)),
		)
		ctx.Request.Header.Set("Content-Type", "application/json")

		recipeAggregate := domain.RecipeAggregate{Recipe: domain.Recipe{Uuid: recipeUuid}}
		mockRecipeService := new(MockRecipeService)
		mockRecipeService.On("Handle", recipeUuid, mock.Anything).Return(&recipeAggregate, nil)

		handler := NewRecipeHandler(mockRecipeService)
		handler.RetrieveRecipeAggregate(ctx)

		assert.Equal(t, 200, ctx.Writer.Status())
	})

	t.Run("HTTP Status 400 on validation error", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = append(ctx.Params, gin.Param{Key: "uuid", Value: recipeUuid.String()})
		ctx.Request = httptest.NewRequest(
			http.MethodPost,
			"/recipes/"+recipeUuid.String()+"/aggregate",
			bytes.NewBuffer([]byte("ERROR BODY")),
		)

		mockRecipeService := new(MockRecipeService)
		handler := NewRecipeHandler(mockRecipeService)
		handler.RetrieveRecipeAggregate(ctx)

		assert.Equal(t, 400, ctx.Writer.Status())
	})

	t.Run("HTTP Status 400 on handler error", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = append(ctx.Params, gin.Param{Key: "uuid", Value: recipeUuid.String()})
		body := `{"pans": [{"shape": "round","measures": {"diameter": "100"}}]}`
		ctx.Request = httptest.NewRequest(
			http.MethodPost,
			"/recipes/"+recipeUuid.String()+"/aggregate",
			bytes.NewBuffer([]byte(body)),
		)

		recipeAggregate := domain.RecipeAggregate{Recipe: domain.Recipe{Uuid: recipeUuid}}
		mockRecipeService := new(MockRecipeService)
		mockRecipeService.On("Handle", recipeUuid, mock.Anything).Return(&recipeAggregate, errors.New("ERROR"))

		handler := NewRecipeHandler(mockRecipeService)
		handler.RetrieveRecipeAggregate(ctx)

		assert.Equal(t, 400, ctx.Writer.Status())
	})

	t.Run("HTTP Status 400 with wrong UUID", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = append(ctx.Params, gin.Param{Key: "uuid", Value: "WRONG UUID"})
		mockRecipeService := new(MockRecipeService)

		handler := NewRecipeHandler(mockRecipeService)

		defer func() {
			if r := recover(); r != nil {
				assert.Contains(t, r, "invalid UUID")
			}
		}()
		handler.RetrieveRecipeAggregate(ctx)
	})
}
