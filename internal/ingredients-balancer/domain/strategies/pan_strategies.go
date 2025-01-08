package strategies

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/cfioretti/recipe-manager/internal/ingredients-balancer/domain"
)

type PanStrategy interface {
	Calculate(measures map[string]interface{}) (domain.Pan, error)
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

func (s *RoundPanStrategy) Calculate(measures map[string]interface{}) (domain.Pan, error) {
	diameter, err := strconv.Atoi(measures["diameter"].(string))
	if err != nil {
		return domain.Pan{}, errors.New("invalid diameter")
	}

	radius := float64(diameter) / 2
	area := math.Pi * radius * radius

	return domain.Pan{
		Shape: "round",
		Measures: domain.Measures{
			Diameter: &diameter,
		},
		DoughWeight: area / 2,
	}, nil
}

func (s *SquarePanStrategy) Calculate(measures map[string]interface{}) (domain.Pan, error) {
	edge, err := strconv.Atoi(measures["edge"].(string))
	if err != nil {
		return domain.Pan{}, errors.New("invalid edge")
	}

	area := float64(edge * edge)

	return domain.Pan{
		Shape: "square",
		Measures: domain.Measures{
			Edge: &edge,
		},
		DoughWeight: area / 2,
	}, nil
}

func (s *RectangularPanStrategy) Calculate(measures map[string]interface{}) (domain.Pan, error) {
	width, widthErr := strconv.Atoi(measures["width"].(string))
	length, lengthErr := strconv.Atoi(measures["length"].(string))
	if widthErr != nil || lengthErr != nil {
		return domain.Pan{}, errors.New("invalid width or length")
	}

	area := float64(width * length)

	return domain.Pan{
		Shape: "rectangular",
		Measures: domain.Measures{
			Width:  &width,
			Length: &length,
		},
		DoughWeight: area / 2,
	}, nil
}
