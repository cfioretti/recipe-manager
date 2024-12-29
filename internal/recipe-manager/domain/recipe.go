package domain

import "github.com/google/uuid"

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

type RecipeAggregate struct {
	Recipe
	SplitIngredients SplitIngredients
}

type SplitIngredients struct {
	SplitDough   []Dough
	SplitTopping []Topping
}

type Dough struct{}

type Topping struct{}

type Steps struct{}
