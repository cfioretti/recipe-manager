package application

import (
	"errors"
	"math"

	"github.com/cfioretti/recipe-manager/internal/ingredients-balancer/domain"
	rdomain "github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
)

const totalPercentage = 100

type IngredientsBalancerService struct{}

func NewIngredientsBalancerService() *IngredientsBalancerService {
	return &IngredientsBalancerService{}
}

func (bs IngredientsBalancerService) Balance(recipe rdomain.Recipe, pans domain.Pans) (*rdomain.RecipeAggregate, error) {
	if pans.TotalArea <= 0 || getFirstIngredientAmount(recipe.Dough.Ingredients) <= 0 {
		return nil, errors.New("invalid dough weight")
	}

	totalDoughWeight := pans.TotalArea / 2
	doughPercentVariation := totalDoughWeight * recipe.Dough.PercentVariation / 100
	doughConversionRatio := (totalDoughWeight + doughPercentVariation) / totalPercentage
	balancedDough := rdomain.Dough{
		PercentVariation: recipe.Dough.PercentVariation,
		Ingredients:      balanceIngredients(recipe.Dough.Ingredients, doughConversionRatio),
	}

	toppingConversionRatio := pans.TotalArea / recipe.Topping.ReferenceArea
	balancedTopping := rdomain.Topping{
		ReferenceArea: recipe.Topping.ReferenceArea,
		Ingredients:   balanceIngredients(recipe.Topping.Ingredients, toppingConversionRatio),
	}

	recipeAggregate := &rdomain.RecipeAggregate{
		Recipe: recipe,
		SplitIngredients: rdomain.SplitIngredients{
			SplitDough:   calculateSplitDoughs(balancedDough, pans),
			SplitTopping: []rdomain.Topping{},
		},
	}
	recipeAggregate.Dough = balancedDough
	recipeAggregate.Topping = balancedTopping

	return recipeAggregate, nil
}

func calculateSplitDoughs(totalDough rdomain.Dough, pans domain.Pans) []rdomain.Dough {
	var splitDoughs []rdomain.Dough

	for _, pan := range pans.Pans {
		ratio := pan.Area / pans.TotalArea

		splitDough := rdomain.Dough{
			Name:        pan.Name,
			Ingredients: make([]rdomain.Ingredient, len(totalDough.Ingredients)),
		}
		splitDough.Ingredients = balanceIngredients(totalDough.Ingredients, ratio)

		splitDoughs = append(splitDoughs, splitDough)
	}

	return splitDoughs
}

func balanceIngredients(ingredients []rdomain.Ingredient, ratio float64) []rdomain.Ingredient {
	balancedIngredients := make([]rdomain.Ingredient, len(ingredients))
	for i, ingredient := range ingredients {
		balancedIngredients[i] = rdomain.Ingredient{
			Name:   ingredient.Name,
			Amount: round(ingredient.Amount * ratio),
		}
	}
	return balancedIngredients
}

func getFirstIngredientAmount(ingredients []rdomain.Ingredient) float64 {
	if len(ingredients) == 0 {
		return 0
	}
	return ingredients[0].Amount
}

func round(num float64) float64 {
	return math.Round(num*10) / 10
}
