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
