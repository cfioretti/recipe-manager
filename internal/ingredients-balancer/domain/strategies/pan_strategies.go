package strategies

import (
	"errors"
	"fmt"
	"math"

	"github.com/cfioretti/recipe-manager/internal/ingredients-balancer/domain"
)

type PanStrategy interface {
	Calculate(measures domain.Measures) (domain.Pan, error)
}

type RoundPanStrategy struct{}
type RectangularPanStrategy struct{}
type SquarePanStrategy struct{}

func GetStrategy(shape string) (PanStrategy, error) {
	switch shape {
	case "round":
		return &RoundPanStrategy{}, nil
	case "square":
		return &SquarePanStrategy{}, nil
	case "rectangular":
		return &RectangularPanStrategy{}, nil
	default:
		return nil, fmt.Errorf("unsupported shape: %s", shape)
	}
}

func (s *RoundPanStrategy) Calculate(measures domain.Measures) (domain.Pan, error) {
	if measures.Diameter == nil {
		return domain.Pan{}, errors.New("diameter is required")
	}

	radius := float64(*measures.Diameter) / 2
	area := math.Pi * radius * radius

	return domain.Pan{
		Shape:       "round",
		Measures:    measures,
		DoughWeight: area / 2,
	}, nil
}

func (s *SquarePanStrategy) Calculate(measures domain.Measures) (domain.Pan, error) {
	if measures.Edge == nil {
		return domain.Pan{}, errors.New("edge is required")
	}

	area := float64(*measures.Edge * *measures.Edge)

	return domain.Pan{
		Shape:       "square",
		Measures:    measures,
		DoughWeight: area / 2,
	}, nil
}

func (s *RectangularPanStrategy) Calculate(measures domain.Measures) (domain.Pan, error) {
	if measures.Width == nil || measures.Length == nil {
		return domain.Pan{}, errors.New("width and length are required")
	}

	area := float64(*measures.Width * *measures.Length)

	return domain.Pan{
		Shape:       "rectangular",
		Measures:    measures,
		DoughWeight: area / 2,
	}, nil
}
