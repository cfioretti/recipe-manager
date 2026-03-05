package dto

import (
	"strconv"

	"github.com/cfioretti/recipe-manager/internal/recipe-manager/application"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
)

type PanRequest struct {
	Pans []Pan `json:"pans" binding:"required,dive"`
}

type Pan struct {
	Shape    string   `json:"shape" binding:"required,oneof=round square rectangular"`
	Measures Measures `json:"measures" binding:"required"`
}

type Measures struct {
	Diameter string `json:"diameter,omitempty" binding:"required_if=Shape round"`
	Edge     string `json:"edge,omitempty" binding:"required_if=Shape square"`
	Width    string `json:"width,omitempty" binding:"required_if=Shape rectangular"`
	Length   string `json:"length,omitempty" binding:"required_if=Shape rectangular"`
}

func (r PanRequest) ToDomain() domain.Pans {
	return pansToDomain(r.Pans)
}

type GenerateRecipeRequest struct {
	Mode   string `json:"mode" binding:"required,oneof=random prompt"`
	Prompt string `json:"prompt,omitempty"`
	Pans   []Pan  `json:"pans" binding:"required,dive"`
}

func (r GenerateRecipeRequest) ToApplication() application.GenerateRecipeCommand {
	return application.GenerateRecipeCommand{
		Mode:   r.Mode,
		Prompt: r.Prompt,
		Pans:   pansToDomain(r.Pans),
	}
}

func pansToDomain(requestPans []Pan) domain.Pans {
	pans := make([]domain.Pan, len(requestPans))
	for i, p := range requestPans {
		pan := domain.Pan{
			Shape:    p.Shape,
			Measures: domain.Measures{},
		}

		switch p.Shape {
		case "round":
			diam, _ := strconv.Atoi(p.Measures.Diameter)
			pan.Measures.Diameter = &diam
		case "square":
			edge, _ := strconv.Atoi(p.Measures.Edge)
			pan.Measures.Edge = &edge
		case "rectangular":
			width, _ := strconv.Atoi(p.Measures.Width)
			length, _ := strconv.Atoi(p.Measures.Length)
			pan.Measures.Width = &width
			pan.Measures.Length = &length
		}

		pans[i] = pan
	}

	return domain.Pans{Pans: pans}
}
