package application

import (
	"github.com/google/uuid"

	bdomain "github.com/cfioretti/recipe-manager/internal/ingredients-balancer/domain"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
)

type RecipeRepository interface {
	GetRecipeByUuid(uuid.UUID) (*domain.Recipe, error)
}

type CalculatorService interface {
	TotalDoughWeightByPans(bdomain.Pans) (*bdomain.Pans, error)
}

type BalancerService interface {
	Balance(domain.Recipe, bdomain.Pans) (*domain.RecipeAggregate, error)
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

func (rs *RecipeService) Handle(recipeUuid uuid.UUID, request bdomain.Pans) (*domain.RecipeAggregate, error) {
	pans, calculatorError := rs.calculator.TotalDoughWeightByPans(request)
	if calculatorError != nil {
		return nil, calculatorError
	}

	recipe, err := rs.repository.GetRecipeByUuid(recipeUuid)
	if err != nil {
		return nil, err
	}

	response, balancerError := rs.balancer.Balance(*recipe, *pans)
	if balancerError != nil {
		return nil, balancerError
	}

	return response, nil
}
