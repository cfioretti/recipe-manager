package domain

import (
	"math"

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
	Name             string
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

func mapDoughListToDTO(doughList []Dough) []dto.SplitDough {
	dtoList := make([]dto.SplitDough, len(doughList))
	for i, d := range doughList {
		totalDoughWeight := calculateTotal(d)
		dtoList[i] = dto.SplitDough{
			Dough: dto.DoughResponse{
				Total: totalDoughWeight,
				Dough: dto.Dough{
					Flour:  d.Flour,
					Water:  d.Water,
					Salt:   d.Salt,
					EvoOil: d.EvoOil,
					Yeast:  d.Yeast,
				},
			},
			Shape: d.Name,
		}
	}
	return dtoList
}

func calculateTotal(d Dough) float64 {
	totalWeight := d.Flour + d.Water + d.Salt + d.EvoOil + d.Yeast
	return math.Round(totalWeight*10) / 10
}
