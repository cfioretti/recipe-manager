package infrastructure

import (
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

func (m *MockRecipeService) Handle(recipeUuid uuid.UUID) *domain.RecipeAggregate {
	args := m.Called(recipeUuid)
	return args.Get(0).(*domain.RecipeAggregate)
}

func TestRetrieveRecipe(t *testing.T) {
	recipeUuid := uuid.New()
	mockService := new(MockRecipeService)

	t.Run("HTTP Status 200 on Success", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = append(ctx.Params, gin.Param{Key: "uuid", Value: recipeUuid.String()})
		expectedResponse := domain.RecipeAggregate{
			Recipe: domain.Recipe{RecipeUuid: recipeUuid},
		}
		mockService.On("Handle", recipeUuid).Return(&expectedResponse)
		controller := NewRecipeController(mockService)
		controller.RetrieveRecipe(ctx)

		assert.Equal(t, 200, ctx.Writer.Status())
	})

	t.Run("Panic with wrong UUID", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = append(ctx.Params, gin.Param{Key: "uuid", Value: "WRONG UUID"})
		controller := NewRecipeController(mockService)

		assert.Panics(t, func() {
			controller.RetrieveRecipe(ctx)
		})
	})
}
