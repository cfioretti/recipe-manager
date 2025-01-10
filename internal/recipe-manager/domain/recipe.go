package domain

import (
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/infrastructure/http/dto"

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
	PercentVariation float64
	Flour            float64
	Water            float64
	Salt             float64
	EvoOil           float64
	Yeast            float64
}

type Topping struct{}

type Steps struct{}

func (r RecipeAggregate) ToDTO() dto.RecipeAggregateResponse {
	return dto.RecipeAggregateResponse{
		Recipe: dto.Recipe{
			Uuid:        r.Recipe.Uuid,
			Name:        r.Recipe.Name,
			Description: r.Recipe.Description,
			Author:      r.Recipe.Author,
			Dough: dto.Dough{
				Flour:  r.Recipe.Dough.Flour,
				Water:  r.Recipe.Dough.Water,
				Salt:   r.Recipe.Dough.Salt,
				EvoOil: r.Recipe.Dough.EvoOil,
				Yeast:  r.Recipe.Dough.Yeast,
			},
			Topping: dto.Topping{},
			Steps:   dto.Steps{},
		},
		SplitIngredients: dto.SplitIngredients{
			SplitDough:   mapDoughListToDTO(r.SplitIngredients.SplitDough),
			SplitTopping: []dto.Topping{},
		},
	}
}

func mapDoughListToDTO(doughList []Dough) []dto.Dough {
	dtoList := make([]dto.Dough, len(doughList))
	for i, d := range doughList {
		dtoList[i] = dto.Dough{
			Flour:  d.Flour,
			Water:  d.Water,
			Salt:   d.Salt,
			EvoOil: d.EvoOil,
			Yeast:  d.Yeast,
		}
	}
	return dtoList
}
