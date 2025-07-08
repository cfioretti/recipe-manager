package application_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cfioretti/recipe-manager/internal/recipe-manager/application"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
)

type StubCalculatorClient struct {
	TotalDoughWeightByPansFunc func(ctx context.Context, pans domain.Pans) (*domain.Pans, error)
}

func (s *StubCalculatorClient) TotalDoughWeightByPans(ctx context.Context, pans domain.Pans) (*domain.Pans, error) {
	return s.TotalDoughWeightByPansFunc(ctx, pans)
}

func (s *StubCalculatorClient) Close() error {
	return nil
}

func TestTotalDoughWeightByPans(t *testing.T) {
	stubClient := createStubCalculatorClient()
	service := application.NewRemoteDoughCalculatorService(stubClient)

	diameter := 28
	pans := domain.Pans{
		Pans: []domain.Pan{
			{
				Shape: "round",
				Measures: domain.Measures{
					Diameter: &diameter,
				},
				Name: "round 28 cm",
				Area: 615.75,
			},
		},
		TotalArea: 615.75,
	}

	result, err := service.TotalDoughWeightByPans(context.Background(), pans)

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
		TotalDoughWeightByPansFunc: func(ctx context.Context, pans domain.Pans) (*domain.Pans, error) {
			return &pans, nil
		},
	}
}
