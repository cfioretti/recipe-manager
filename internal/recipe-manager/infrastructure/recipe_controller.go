package infrastructure

import (
	"net/http"

	"recipe-manager/internal/recipe-manager/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RecipeHandler interface {
	Handle(uuid.UUID) *domain.Recipe
}

type RecipeController struct {
	recipeHandler RecipeHandler
}

func NewRecipeController(recipeHandler RecipeHandler) *RecipeController {
	return &RecipeController{recipeHandler: recipeHandler}
}

func (rc *RecipeController) RetrieveRecipe(ctx *gin.Context) {
	recipeUuid := uuid.MustParse(ctx.Param("uuid"))
	result := rc.recipeHandler.Handle(recipeUuid)
	ctx.JSON(http.StatusOK, result)
}
