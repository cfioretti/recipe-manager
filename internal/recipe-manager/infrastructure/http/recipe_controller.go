package http

import (
	"net/http"

	"recipe-manager/internal/recipe-manager/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RecipeHandler interface {
	Handle(uuid.UUID) (*domain.RecipeAggregate, error)
}

type Calculator interface {
	TotalPansWeight(gin.Params) (int, error)
}

type RecipeController struct {
	recipeHandler RecipeHandler
	calculator    Calculator
}

func NewRecipeController(recipeHandler RecipeHandler, calculator Calculator) *RecipeController {
	return &RecipeController{
		recipeHandler: recipeHandler,
		calculator:    calculator,
	}
}

func (rc *RecipeController) RetrieveRecipeAggregate(ctx *gin.Context) {
	recipeUuid := uuid.MustParse(ctx.Param("uuid"))
	_, validationError := rc.calculator.TotalPansWeight(ctx.Params)
	if validationError != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			validationError.Error(),
		)
	}
	result, err := rc.recipeHandler.Handle(recipeUuid)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			err.Error(),
		)
	}
	ctx.JSON(
		http.StatusOK,
		result,
	)
}
