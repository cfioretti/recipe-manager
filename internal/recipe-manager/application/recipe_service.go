package application

import (
	calculatordomain "github.com/cfioretti/recipe-manager/internal/ingredients-balancer/domain"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"

	"github.com/google/uuid"
)

type RecipeRepository interface {
	GetRecipeByUuid(uuid.UUID) (*domain.Recipe, error)
}

type RecipeService struct {
	repository RecipeRepository
	calculator CalculatorService
	balancer   BalancerService
}

type CalculatorService interface {
	TotalDoughWeightByPans([]byte) (*calculatordomain.Pans, error)
}

type BalancerService interface {
	Balance(domain.Recipe, calculatordomain.Pans) (*domain.RecipeAggregate, error)
}

func NewRecipeService(repository RecipeRepository, calculator CalculatorService, balancer BalancerService) *RecipeService {
	return &RecipeService{
		repository: repository,
		calculator: calculator,
		balancer:   balancer,
	}
}

func (rs *RecipeService) Handle(recipeUuid uuid.UUID, request []byte) (*domain.RecipeAggregate, error) {
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
