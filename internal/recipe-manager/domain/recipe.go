package domain

import (
	"github.com/google/uuid"
)

type RecipeAggregate struct {
	Recipe
	SplitIngredients SplitIngredients
}

type Recipe struct {
	Id          int
	Uuid        uuid.UUID
	Name        string
	Description string
	Author      string
	Dough       Dough
	Topping     Topping
	Steps       Steps
}

type SplitIngredients struct {
	SplitDough   []Dough
	SplitTopping []Topping
}

type Dough struct {
	Name             string
	PercentVariation float64
	Ingredients      []Ingredient
}

type Topping struct {
	Name          string
	ReferenceArea float64
	Ingredients   []Ingredient
}

type Ingredient struct {
	Name   string
	Amount float64
}

type Steps struct{}
