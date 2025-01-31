package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	balancerdomain "github.com/cfioretti/recipe-manager/internal/ingredients-balancer/domain"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/interfaces/api/http/dto"
)

type RecipeService interface {
	Handle(uuid.UUID, balancerdomain.Pans) (*domain.RecipeAggregate, error)
}

type RecipeHandler struct {
	recipeService RecipeService
}

func NewRecipeHandler(recipeService RecipeService) *RecipeHandler {
	return &RecipeHandler{recipeService: recipeService}
}

func (rc *RecipeHandler) RetrieveRecipeAggregate(ctx *gin.Context) {
	recipeUuid := uuid.MustParse(ctx.Param("uuid"))
	var requestBody dto.PanRequest
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	recipe, err := rc.recipeService.Handle(recipeUuid, requestBody.ToDomain())
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
	}
	aggregateResponse := dto.DomainToDTO(*recipe)
	ctx.JSON(
		http.StatusOK,
		gin.H{"data": &aggregateResponse},
	)
}
