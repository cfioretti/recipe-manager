package application

import (
	balancerdomain "github.com/cfioretti/recipe-manager/internal/ingredients-balancer/domain"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/infrastructure/grpc/client"
)

type RemoteDoughCalculatorService struct {
	client *client.DoughCalculatorClient
}

func NewRemoteDoughCalculatorService(client *client.DoughCalculatorClient) *RemoteDoughCalculatorService {
	return &RemoteDoughCalculatorService{
		client: client,
	}
}

func (dc *RemoteDoughCalculatorService) TotalDoughWeightByPans(pans balancerdomain.Pans) (*balancerdomain.Pans, error) {
	return dc.client.TotalDoughWeightByPans(pans)
}
