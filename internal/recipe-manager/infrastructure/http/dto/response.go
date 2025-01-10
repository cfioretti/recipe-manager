package dto

import "github.com/google/uuid"

type RecipeAggregateResponse struct {
	Recipe
	SplitIngredients SplitIngredients `json:"splitIngredients"`
}

type Recipe struct {
	Uuid        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Dough       Dough     `json:"dough"`
	Topping     Topping   `json:"topping"`
	Steps       Steps     `json:"steps"`
}

type SplitIngredients struct {
	SplitDough   []Dough   `json:"splitDough"`
	SplitTopping []Topping `json:"splitTopping"`
}

type Dough struct {
	Flour  float64 `json:"flour"`
	Water  float64 `json:"water"`
	Salt   float64 `json:"salt"`
	EvoOil float64 `json:"evoOil"`
	Yeast  float64 `json:"yeast"`
}

type Topping struct{}

type Steps struct{}
