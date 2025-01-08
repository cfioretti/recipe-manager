package application

import (
	"errors"
	"math"

	"github.com/cfioretti/recipe-manager/internal/ingredients-balancer/domain"
	recipedomain "github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
)

const totalPercentage = 100

type DoughBalancerService struct{}

func NewDoughBalancerService() *DoughBalancerService {
	return &DoughBalancerService{}
}

func (dbs DoughBalancerService) Balance(recipe recipedomain.Recipe, pans domain.Pans) (*recipedomain.RecipeAggregate, error) {
	if pans.TotalDoughWeight <= 0 || recipe.Dough.Flour <= 0 {
		return nil, errors.New("invalid dough weight")
	}

	percentVariation := pans.TotalDoughWeight * recipe.Dough.PercentVariation / 100
	conversionRatio := (pans.TotalDoughWeight + percentVariation) / totalPercentage
	balancedDough := recipedomain.Dough{
		Flour:  round(recipe.Dough.Flour * conversionRatio),
		Water:  round(recipe.Dough.Water * conversionRatio),
		Salt:   round(recipe.Dough.Salt * conversionRatio),
		EvoOil: round(recipe.Dough.EvoOil * conversionRatio),
		Yeast:  round(recipe.Dough.Yeast * conversionRatio),
	}

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
		ratio := pan.DoughWeight / totalDoughWeight
		splitDough := recipedomain.Dough{
			Flour:  round(totalDough.Flour * ratio),
			Water:  round(totalDough.Water * ratio),
			Salt:   round(totalDough.Salt * ratio),
			EvoOil: round(totalDough.EvoOil * ratio),
			Yeast:  round(totalDough.Yeast * ratio),
		}
		splitDoughs = append(splitDoughs, splitDough)
	}

	return splitDoughs
}

func round(num float64) float64 {
	return math.Round(num*10) / 10
}
