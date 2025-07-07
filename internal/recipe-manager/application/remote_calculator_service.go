package application

import (
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
)

type CalculatorClient interface {
	TotalDoughWeightByPans(pans domain.Pans) (*domain.Pans, error)
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

func (dc *RemoteCalculatorService) TotalDoughWeightByPans(pans domain.Pans) (*domain.Pans, error) {
	return dc.client.TotalDoughWeightByPans(pans)
}
