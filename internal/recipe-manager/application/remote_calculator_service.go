package application

import (
	bdomain "github.com/cfioretti/recipe-manager/internal/ingredients-balancer/domain"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/infrastructure/grpc/client"
)

type RemoteCalculatorService struct {
	client *client.CalculatorClient
}

func NewRemoteDoughCalculatorService(client *client.CalculatorClient) *RemoteCalculatorService {
	return &RemoteCalculatorService{
		client: client,
	}
}

func (dc *RemoteCalculatorService) TotalDoughWeightByPans(pans bdomain.Pans) (*bdomain.Pans, error) {
	return dc.client.TotalDoughWeightByPans(pans)
}
