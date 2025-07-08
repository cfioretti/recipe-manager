package application

import (
	"context"

	"github.com/google/uuid"

	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
)

type RecipeRepository interface {
	GetRecipeByUuid(uuid.UUID) (*domain.Recipe, error)
}

type CalculatorService interface {
	TotalDoughWeightByPans(context.Context, domain.Pans) (*domain.Pans, error)
}

type BalancerService interface {
	Balance(context.Context, domain.Recipe, domain.Pans) (*domain.RecipeAggregate, error)
}

type RecipeService struct {
	repository RecipeRepository
	calculator CalculatorService
	balancer   BalancerService
}

func NewRecipeService(repository RecipeRepository, calculator CalculatorService, balancer BalancerService) *RecipeService {
	return &RecipeService{
		repository: repository,
		calculator: calculator,
		balancer:   balancer,
	}
}

func (rs *RecipeService) Handle(ctx context.Context, recipeUuid uuid.UUID, request domain.Pans) (*domain.RecipeAggregate, error) {
	pans, calculatorError := rs.calculator.TotalDoughWeightByPans(ctx, request)
	if calculatorError != nil {
		return nil, calculatorError
	}

	recipe, err := rs.repository.GetRecipeByUuid(recipeUuid)
	if err != nil {
		return nil, err
	}

	response, balancerError := rs.balancer.Balance(ctx, *recipe, *pans)
	if balancerError != nil {
		return nil, balancerError
	}

	return response, nil
}
