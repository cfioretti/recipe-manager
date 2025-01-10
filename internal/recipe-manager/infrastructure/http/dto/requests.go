package dto

import (
	"strconv"

	"github.com/cfioretti/recipe-manager/internal/ingredients-balancer/domain"
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
	pans := make([]domain.Pan, len(r.Pans))
	for i, p := range r.Pans {
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

	return domain.Pans{
		Pans: pans,
	}
}
