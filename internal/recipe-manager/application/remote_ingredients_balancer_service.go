package application

import (
	bdomain "github.com/cfioretti/recipe-manager/internal/ingredients-balancer/domain"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
)

type IngredientsBalancerClient interface {
	Balance(recipe domain.Recipe, pans bdomain.Pans) (*domain.RecipeAggregate, error)
	Close() error
}

type RemoteIngredientsBalancerService struct {
	client IngredientsBalancerClient
}

func NewRemoteIngredientsBalancerService(client IngredientsBalancerClient) *RemoteIngredientsBalancerService {
	return &RemoteIngredientsBalancerService{
		client: client,
	}
}

func (bs *RemoteIngredientsBalancerService) Balance(recipe domain.Recipe, pans bdomain.Pans) (*domain.RecipeAggregate, error) {
	return bs.client.Balance(recipe, pans)
}
