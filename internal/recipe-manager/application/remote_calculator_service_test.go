package application_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	bdomain "github.com/cfioretti/recipe-manager/internal/ingredients-balancer/domain"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/application"
)

type StubCalculatorClient struct {
	TotalDoughWeightByPansFunc func(pans bdomain.Pans) (*bdomain.Pans, error)
}

func (s *StubCalculatorClient) TotalDoughWeightByPans(pans bdomain.Pans) (*bdomain.Pans, error) {
	return s.TotalDoughWeightByPansFunc(pans)
}

func (s *StubCalculatorClient) Close() error {
	return nil
}

func TestTotalDoughWeightByPans(t *testing.T) {
	stubClient := createStubCalculatorClient()
	service := application.NewRemoteDoughCalculatorService(stubClient)

	diameter := 28
	pans := bdomain.Pans{
		Pans: []bdomain.Pan{
			{
				Shape: "round",
				Measures: bdomain.Measures{
					Diameter: &diameter,
				},
				Name: "round 28 cm",
				Area: 615.75,
			},
		},
		TotalArea: 615.75,
	}

	result, err := service.TotalDoughWeightByPans(pans)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Pans))
	assert.Equal(t, "round", result.Pans[0].Shape)
	assert.Equal(t, "round 28 cm", result.Pans[0].Name)
	assert.Equal(t, 615.75, result.Pans[0].Area)
	assert.Equal(t, 615.75, result.TotalArea)
}

func createStubCalculatorClient() *StubCalculatorClient {
	return &StubCalculatorClient{
		TotalDoughWeightByPansFunc: func(pans bdomain.Pans) (*bdomain.Pans, error) {
			// Return a copy of the input pans
			return &pans, nil
		},
	}
}
