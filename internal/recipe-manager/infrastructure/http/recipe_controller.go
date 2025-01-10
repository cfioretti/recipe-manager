package http

import (
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/infrastructure/http/dto"
	"net/http"

	balancerdomain "github.com/cfioretti/recipe-manager/internal/ingredients-balancer/domain"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RecipeHandler interface {
	Handle(uuid.UUID, balancerdomain.Pans) (*domain.RecipeAggregate, error)
}

type RecipeController struct {
	recipeHandler RecipeHandler
}

func NewRecipeController(recipeHandler RecipeHandler) *RecipeController {
	return &RecipeController{recipeHandler: recipeHandler}
}

func (rc *RecipeController) RetrieveRecipeAggregate(ctx *gin.Context) {
	recipeUuid := uuid.MustParse(ctx.Param("uuid"))
	var requestBody dto.PanRequest
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	recipe, err := rc.recipeHandler.Handle(recipeUuid, requestBody.ToDomain())
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
	}
	ctx.JSON(
		http.StatusOK,
		&recipe,
	)
}
