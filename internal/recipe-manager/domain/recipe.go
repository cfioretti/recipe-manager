package domain

import "github.com/google/uuid"

type Recipe struct {
	RecipeUuid uuid.UUID
	Dough      Dough
	Topping    Topping
	Steps      Steps
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
