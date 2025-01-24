package domain

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
