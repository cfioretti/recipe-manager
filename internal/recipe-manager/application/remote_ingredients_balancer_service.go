package application

import (
	"context"

	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
)

type IngredientsBalancerClient interface {
	Balance(context.Context, domain.Recipe, domain.Pans) (*domain.RecipeAggregate, error)
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

func (bs *RemoteIngredientsBalancerService) Balance(ctx context.Context, recipe domain.Recipe, pans domain.Pans) (*domain.RecipeAggregate, error) {
	return bs.client.Balance(ctx, recipe, pans)
}
