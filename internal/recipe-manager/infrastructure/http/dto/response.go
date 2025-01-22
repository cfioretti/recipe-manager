package dto

import (
	"math"

	"github.com/google/uuid"

	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
)

type RecipeAggregateResponse struct {
	Recipe
	SplitIngredients SplitIngredients `json:"splitIngredients"`
}

type Recipe struct {
	Uuid        uuid.UUID     `json:"uuid"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Author      string        `json:"author"`
	Dough       DoughResponse `json:"dough"`
	Topping     Topping       `json:"topping"`
	Steps       Steps         `json:"steps"`
}

type SplitIngredients struct {
	SplitDough   []SplitDough `json:"splitDough"`
	SplitTopping []Topping    `json:"splitTopping"`
}

type SplitDough struct {
	Shape string        `json:"shape"`
	Dough DoughResponse `json:"dough"`
}

type DoughResponse struct {
	Total float64 `json:"total"`
	Dough
}

type Dough struct {
	Ingredients []Ingredient
}

type Ingredient struct {
	Name   string
	Amount float64
}

type Topping struct{}

type Steps struct{}

func DomainToDTO(r domain.RecipeAggregate) RecipeAggregateResponse {
	return RecipeAggregateResponse{
		Recipe: Recipe{
			Uuid:        r.Recipe.Uuid,
			Name:        r.Recipe.Name,
			Description: r.Recipe.Description,
			Author:      r.Recipe.Author,
			Dough: DoughResponse{
				Total: calculateTotal(r.Dough),
				Dough: Dough{
					Ingredients: mapIngredientsToDTO(r.Recipe.Dough.Ingredients),
				},
			},
			Topping: Topping{},
			Steps:   Steps{},
		},
		SplitIngredients: SplitIngredients{
			SplitDough:   mapDoughListToDTO(r.SplitIngredients.SplitDough),
			SplitTopping: []Topping{},
		},
	}
}

func mapDoughListToDTO(doughList []domain.Dough) []SplitDough {
	dtoList := make([]SplitDough, len(doughList))
	for i, d := range doughList {
		totalDoughWeight := calculateTotal(d)
		dtoList[i] = SplitDough{
			Dough: DoughResponse{
				Total: totalDoughWeight,
				Dough: Dough{
					Ingredients: mapIngredientsToDTO(d.Ingredients),
				},
			},
			Shape: d.Name,
		}
	}
	return dtoList
}

func mapIngredientsToDTO(ingredients []domain.Ingredient) []Ingredient {
	dtoIngredients := make([]Ingredient, len(ingredients))
	for i, ing := range ingredients {
		dtoIngredients[i] = Ingredient{
			Name:   ing.Name,
			Amount: ing.Amount,
		}
	}
	return dtoIngredients
}

func calculateTotal(d domain.Dough) float64 {
	totalWeight := 0.0
	for _, ingredient := range d.Ingredients {
		totalWeight += ingredient.Amount
	}
	return math.Round(totalWeight*10) / 10
}
