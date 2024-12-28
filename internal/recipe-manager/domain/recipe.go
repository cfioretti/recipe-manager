package domain

import "github.com/google/uuid"

type Recipe struct {
	RecipeUuid uuid.UUID
	Total      TotalIngredients
	Split      SplitIngredients
	Steps      Steps
}

type TotalIngredients struct {
	Dough   Dough
	Topping Topping
}

type SplitIngredients struct {
	Dough   []Dough
	Topping []Topping
}

type Dough struct{}

type Topping struct{}

type Steps struct{}
