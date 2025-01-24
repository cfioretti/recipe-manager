package application

import (
	"errors"
	"math"

	"github.com/cfioretti/recipe-manager/internal/ingredients-balancer/domain"
	recipedomain "github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
)

const totalPercentage = 100

type IngredientsBalancerService struct{}

func NewIngredientsBalancerService() *IngredientsBalancerService {
	return &IngredientsBalancerService{}
}

func (bs IngredientsBalancerService) Balance(recipe recipedomain.Recipe, pans domain.Pans) (*recipedomain.RecipeAggregate, error) {
	if pans.TotalArea <= 0 || getFirstIngredientAmount(recipe.Dough.Ingredients) <= 0 {
		return nil, errors.New("invalid dough weight")
	}

	totalDoughWeight := pans.TotalArea / 2
	doughPercentVariation := totalDoughWeight * recipe.Dough.PercentVariation / 100
	doughConversionRatio := (totalDoughWeight + doughPercentVariation) / totalPercentage
	balancedDough := recipedomain.Dough{
		PercentVariation: recipe.Dough.PercentVariation,
		Ingredients:      balanceIngredients(recipe.Dough.Ingredients, doughConversionRatio),
	}

	toppingConversionRatio := pans.TotalArea / recipe.Topping.ReferenceArea
	balancedTopping := recipedomain.Topping{
		ReferenceArea: recipe.Topping.ReferenceArea,
		Ingredients:   balanceIngredients(recipe.Topping.Ingredients, toppingConversionRatio),
	}

	recipeAggregate := &recipedomain.RecipeAggregate{
		Recipe: recipe,
		SplitIngredients: recipedomain.SplitIngredients{
			SplitDough:   calculateSplitDoughs(balancedDough, pans),
			SplitTopping: []recipedomain.Topping{},
		},
	}
	recipeAggregate.Dough = balancedDough
	recipeAggregate.Topping = balancedTopping

	return recipeAggregate, nil
}

func calculateSplitDoughs(totalDough recipedomain.Dough, pans domain.Pans) []recipedomain.Dough {
	var splitDoughs []recipedomain.Dough

	for _, pan := range pans.Pans {
		ratio := pan.Area / pans.TotalArea

		splitDough := recipedomain.Dough{
			Name:        pan.Name,
			Ingredients: make([]recipedomain.Ingredient, len(totalDough.Ingredients)),
		}
		splitDough.Ingredients = balanceIngredients(totalDough.Ingredients, ratio)

		splitDoughs = append(splitDoughs, splitDough)
	}

	return splitDoughs
}

func balanceIngredients(ingredients []recipedomain.Ingredient, ratio float64) []recipedomain.Ingredient {
	balancedIngredients := make([]recipedomain.Ingredient, len(ingredients))
	for i, ingredient := range ingredients {
		balancedIngredients[i] = recipedomain.Ingredient{
			Name:   ingredient.Name,
			Amount: round(ingredient.Amount * ratio),
		}
	}
	return balancedIngredients
}

func getFirstIngredientAmount(ingredients []recipedomain.Ingredient) float64 {
	if len(ingredients) == 0 {
		return 0
	}
	return ingredients[0].Amount
}

func round(num float64) float64 {
	return math.Round(num*10) / 10
}
