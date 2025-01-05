package application

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"strconv"

	"github.com/cfioretti/recipe-manager/internal/dough-calculator/domain"
)

type DoughCalculatorService struct{}

func NewDoughCalculatorService() *DoughCalculatorService {
	return &DoughCalculatorService{}
}

func (dc DoughCalculatorService) TotalPansWeight(body []byte) (*domain.Pans, error) {
	var input struct {
		Pans []map[string]interface{} `json:"pans"`
	}
	if err := json.Unmarshal(body, &input); err != nil {
		return nil, errors.New("invalid JSON format")
	}

	var pans domain.Pans
	for _, item := range input.Pans {
		shape := item["shape"].(string)
		measures := item["measures"].(map[string]interface{})

		if strategy, exists := strategies[shape]; exists {
			pan := strategy(measures)
			pans.Pans = append(pans.Pans, pan)
			pans.Total += pan.Area
		} else {
			log.Printf("Shape '%s' not supported", shape)
		}
	}
	return &pans, nil
}

type Strategy func(data map[string]interface{}) domain.Pan

var strategies = map[string]Strategy{
	"round": func(data map[string]interface{}) domain.Pan {
		diameter, _ := strconv.Atoi(data["diameter"].(string))
		radius := float64(diameter) / 2
		area := math.Pi * radius * radius
		return domain.Pan{
			Shape: "round",
			Measures: domain.Measures{
				Diameter: &diameter,
			},
			Area: area,
		}
	},
	"square": func(data map[string]interface{}) domain.Pan {
		edge, _ := strconv.Atoi(data["edge"].(string))
		area := float64(edge * edge)
		return domain.Pan{
			Shape: "square",
			Measures: domain.Measures{
				Edge: &edge,
			},
			Area: area,
		}
	},
	"rectangular": func(data map[string]interface{}) domain.Pan {
		width, _ := strconv.Atoi(data["width"].(string))
		length, _ := strconv.Atoi(data["length"].(string))
		area := float64(width * length)
		return domain.Pan{
			Shape: "rectangular",
			Measures: domain.Measures{
				Width:  &width,
				Length: &length,
			},
			Area: area,
		}
	},
}
