package http

import (
	"io"
	"net/http"

	calculatordomain "github.com/cfioretti/recipe-manager/internal/ingredients-balancer/domain"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RecipeHandler interface {
	Handle(uuid.UUID) (*domain.RecipeAggregate, error)
}

type Calculator interface {
	TotalDoughWeightByPans([]byte) (*calculatordomain.Pans, error)
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
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			err.Error(),
		)
	}
	_, validationError := rc.calculator.TotalDoughWeightByPans(body)
	if validationError != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			validationError.Error(),
		)
	}
	recipe, err := rc.recipeHandler.Handle(recipeUuid)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			err.Error(),
		)
	}
	ctx.JSON(
		http.StatusOK,
		recipe,
	)
}
