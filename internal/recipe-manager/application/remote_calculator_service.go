package application

import (
	"context"

	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
)

type CalculatorClient interface {
	TotalDoughWeightByPans(context.Context, domain.Pans) (*domain.Pans, error)
	Close() error
}

type RemoteCalculatorService struct {
	client CalculatorClient
}

func NewRemoteDoughCalculatorService(client CalculatorClient) *RemoteCalculatorService {
	return &RemoteCalculatorService{
		client: client,
	}
}

func (dc *RemoteCalculatorService) TotalDoughWeightByPans(ctx context.Context, pans domain.Pans) (*domain.Pans, error) {
	return dc.client.TotalDoughWeightByPans(ctx, pans)
}
