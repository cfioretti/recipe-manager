package http

import (
	"io"
	"net/http"

	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RecipeHandler interface {
	Handle(uuid.UUID, []byte) (*domain.RecipeAggregate, error)
}

type RecipeController struct {
	recipeHandler RecipeHandler
}

func NewRecipeController(recipeHandler RecipeHandler) *RecipeController {
	return &RecipeController{recipeHandler: recipeHandler}
}

func (rc *RecipeController) RetrieveRecipeAggregate(ctx *gin.Context) {
	recipeUuid := uuid.MustParse(ctx.Param("uuid"))
	requestBody, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			err.Error(),
		)
	}
	recipe, err := rc.recipeHandler.Handle(recipeUuid, requestBody)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			err.Error(),
		)
	}
	ctx.JSON(
		http.StatusOK,
		&recipe,
	)
}
