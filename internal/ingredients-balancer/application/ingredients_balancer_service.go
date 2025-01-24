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
	if pans.TotalDoughWeight <= 0 || getFirstIngredientAmount(recipe.Dough.Ingredients) <= 0 {
		return nil, errors.New("invalid dough weight")
	}

	percentVariation := pans.TotalDoughWeight * recipe.Dough.PercentVariation / 100
	conversionRatio := (pans.TotalDoughWeight + percentVariation) / totalPercentage
	balancedDough := recipedomain.Dough{
		PercentVariation: recipe.Dough.PercentVariation,
		Ingredients:      make([]recipedomain.Ingredient, len(recipe.Dough.Ingredients)),
	}
	balancedDough.Ingredients = balanceIngredients(recipe.Dough.Ingredients, conversionRatio)

	splitDoughs := calculateSplitDoughs(balancedDough, pans)

	recipeAggregate := &recipedomain.RecipeAggregate{
		Recipe: recipe,
		SplitIngredients: recipedomain.SplitIngredients{
			SplitDough:   splitDoughs,
			SplitTopping: []recipedomain.Topping{},
		},
	}
	recipeAggregate.Dough = balancedDough

	return recipeAggregate, nil
}

func calculateSplitDoughs(totalDough recipedomain.Dough, pans domain.Pans) []recipedomain.Dough {
	var splitDoughs []recipedomain.Dough

	totalDoughWeight := pans.TotalDoughWeight
	for _, pan := range pans.Pans {
		ratio := (pan.Area / 2) / totalDoughWeight

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
