package application

import (
	"encoding/json"
	"errors"

	"github.com/cfioretti/recipe-manager/internal/ingredients-balancer/domain"
	"github.com/cfioretti/recipe-manager/internal/ingredients-balancer/domain/strategies"
)

type DoughCalculatorService struct{}

func NewDoughCalculatorService() *DoughCalculatorService {
	return &DoughCalculatorService{}
}

type Input struct {
	Pans []PanInput `json:"pans"`
}

type PanInput struct {
	Shape    string                 `json:"shape"`
	Measures map[string]interface{} `json:"measures"`
}

func (dc DoughCalculatorService) TotalDoughWeightByPans(body []byte) (*domain.Pans, error) {
	var input Input
	if err := json.Unmarshal(body, &input); err != nil {
		return nil, errors.New("invalid JSON format")
	}

	var result domain.Pans
	for _, item := range input.Pans {
		strategy, err := strategies.GetStrategy(item.Shape)
		if err != nil {
			return nil, errors.New("unsupported shape")
		}

		pan, err := strategy.Calculate(item.Measures)
		if err != nil {
			return nil, errors.New("error processing pan")
		}

		result.Pans = append(result.Pans, pan)
		result.TotalDoughWeight += pan.DoughWeight
	}
	return &result, nil
}
