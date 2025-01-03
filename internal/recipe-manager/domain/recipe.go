package domain

import "github.com/google/uuid"

type Recipe struct {
	Id          int       `json:"id"`
	Uuid        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Dough       Dough     `json:"dough"`
	Topping     Topping   `json:"topping"`
	Steps       Steps     `json:"steps"`
}

type RecipeAggregate struct {
	Recipe
	SplitIngredients SplitIngredients
}

type SplitIngredients struct {
	SplitDough   []Dough
	SplitTopping []Topping
}

type Dough struct {
	Total  float64 `json:"total"`
	Flour  float64 `json:"flour"`
	Water  float64 `json:"water"`
	Salt   float64 `json:"salt"`
	EvoOil float64 `json:"evoOil"`
	Yeast  float64 `json:"yeast"`
}

type Topping struct{}

type Steps struct{}
