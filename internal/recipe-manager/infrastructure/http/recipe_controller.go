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
	Handle(uuid.UUID) (*domain.Recipe, error)
}

type Calculator interface {
	TotalDoughWeightByPans([]byte) (*calculatordomain.Pans, error)
}

type Balancer interface {
	Balance(domain.Recipe, calculatordomain.Pans) (*domain.RecipeAggregate, error)
}

type RecipeController struct {
	recipeHandler RecipeHandler
	calculator    Calculator
	balancer      Balancer
}

func NewRecipeController(recipeHandler RecipeHandler, calculator Calculator, balancer Balancer) *RecipeController {
	return &RecipeController{
		recipeHandler: recipeHandler,
		calculator:    calculator,
		balancer:      balancer,
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
	pans, calculatorError := rc.calculator.TotalDoughWeightByPans(body)
	if calculatorError != nil {
		rc.abortWithStatus400(ctx, calculatorError)
	}
	recipe, handlerErr := rc.recipeHandler.Handle(recipeUuid)
	if handlerErr != nil {
		rc.abortWithStatus400(ctx, handlerErr)
	}
	_, balancerError := rc.balancer.Balance(*recipe, *pans)
	if balancerError != nil {
		rc.abortWithStatus400(ctx, balancerError)
	}
	response := domain.RecipeAggregate{Recipe: *recipe}
	ctx.JSON(
		http.StatusOK,
		&response,
	)
}

func (rc *RecipeController) abortWithStatus400(ctx *gin.Context, err error) {
	ctx.AbortWithStatusJSON(
		http.StatusBadRequest,
		err.Error(),
	)
}
