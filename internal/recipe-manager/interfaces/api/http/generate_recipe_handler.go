package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cfioretti/recipe-manager/internal/recipe-manager/application"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/interfaces/api/http/dto"
)

type GenerateRecipeService interface {
	Handle(context.Context, application.GenerateRecipeCommand) (*domain.RecipeAggregate, error)
}

type GenerateRecipeHandler struct {
	generateRecipeService GenerateRecipeService
}

func NewGenerateRecipeHandler(generateRecipeService GenerateRecipeService) *GenerateRecipeHandler {
	return &GenerateRecipeHandler{generateRecipeService: generateRecipeService}
}

func (h *GenerateRecipeHandler) GenerateRecipeAggregate(ctx *gin.Context) {
	var request dto.GenerateRecipeRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	aggregate, err := h.generateRecipeService.Handle(ctx.Request.Context(), request.ToApplication())
	if err != nil {
		errorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	aggregateResponse := dto.DomainToDTO(*aggregate)
	writeDataResponse(ctx, http.StatusOK, &aggregateResponse)
}
