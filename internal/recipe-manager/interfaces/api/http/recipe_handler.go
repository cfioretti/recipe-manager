package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	bdomain "github.com/cfioretti/recipe-manager/internal/ingredients-balancer/domain"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/interfaces/api/http/dto"
)

type RecipeService interface {
	Handle(uuid.UUID, bdomain.Pans) (*domain.RecipeAggregate, error)
}

type RecipeHandler struct {
	recipeService RecipeService
}

func NewRecipeHandler(recipeService RecipeService) *RecipeHandler {
	return &RecipeHandler{recipeService: recipeService}
}

func (rc *RecipeHandler) RetrieveRecipeAggregate(ctx *gin.Context) {
	recipeUuid, uuidErr := uuid.Parse(ctx.Param("uuid"))
	if uuidErr != nil {
		errorResponse(ctx, http.StatusBadRequest, "invalid UUID")
		return
	}

	var requestBody dto.PanRequest
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	recipe, err := rc.recipeService.Handle(recipeUuid, requestBody.ToDomain())
	if err != nil {
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	aggregateResponse := dto.DomainToDTO(*recipe)
	ctx.JSON(
		http.StatusOK,
		gin.H{"data": &aggregateResponse},
	)
}

func errorResponse(ctx *gin.Context, statusCode int, errorMsg string) {
	ctx.AbortWithStatusJSON(
		statusCode,
		gin.H{"error": errorMsg},
	)
}
