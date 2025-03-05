package application

import (
	bdomain "github.com/cfioretti/recipe-manager/internal/ingredients-balancer/domain"
)

type CalculatorClient interface {
	TotalDoughWeightByPans(pans bdomain.Pans) (*bdomain.Pans, error)
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

func (dc *RemoteCalculatorService) TotalDoughWeightByPans(pans bdomain.Pans) (*bdomain.Pans, error) {
	return dc.client.TotalDoughWeightByPans(pans)
}
