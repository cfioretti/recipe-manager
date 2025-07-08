package client

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/cfioretti/recipe-manager/internal/infrastructure/logging"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
	pb "github.com/cfioretti/recipe-manager/internal/recipe-manager/infrastructure/grpc/proto/generated"
)

type IngredientsBalancerClient struct {
	client     pb.IngredientsBalancerClient
	conn       *grpc.ClientConn
	serverAddr string
	timeout    time.Duration
}

func NewIngredientsBalancerClient(serverAddr string, timeout time.Duration) (*IngredientsBalancerClient, error) {
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewIngredientsBalancerClient(conn)

	return &IngredientsBalancerClient{
		client:     client,
		conn:       conn,
		serverAddr: serverAddr,
		timeout:    timeout,
	}, nil
}

func (c *IngredientsBalancerClient) Close() error {
	return c.conn.Close()
}

func (c *IngredientsBalancerClient) Balance(ctx context.Context, recipe domain.Recipe, pans domain.Pans) (*domain.RecipeAggregate, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	correlationID := logging.GetCorrelationID(ctx)
	md := metadata.Pairs("x-correlation-id", correlationID)
	timeoutCtx = metadata.NewOutgoingContext(timeoutCtx, md)

	protoRecipe := toProtoRecipe(recipe)
	protoPans := toProtoPans(pans)

	response, err := c.client.Balance(timeoutCtx, &pb.BalanceRequest{
		Recipe: protoRecipe,
		Pans:   protoPans,
	})
	if err != nil {
		return nil, err
	}

	result := toDomainRecipeAggregate(response.RecipeAggregate)
	return result, nil
}

func toProtoRecipe(recipe domain.Recipe) *pb.Recipe {
	return &pb.Recipe{
		Id:          int32(recipe.Id),
		Uuid:        recipe.Uuid.String(),
		Name:        recipe.Name,
		Description: recipe.Description,
		Author:      recipe.Author,
		Dough:       toProtoDough(recipe.Dough),
		Topping:     toProtoTopping(recipe.Topping),
		Steps:       toProtoSteps(recipe.Steps),
	}
}

func toProtoDough(dough domain.Dough) *pb.Dough {
	ingredients := make([]*pb.Ingredient, 0, len(dough.Ingredients))
	for _, ingredient := range dough.Ingredients {
		ingredients = append(ingredients, &pb.Ingredient{
			Name:   ingredient.Name,
			Amount: ingredient.Amount,
		})
	}

	return &pb.Dough{
		Name:             dough.Name,
		PercentVariation: dough.PercentVariation,
		Ingredients:      ingredients,
	}
}

func toProtoTopping(topping domain.Topping) *pb.Topping {
	ingredients := make([]*pb.Ingredient, 0, len(topping.Ingredients))
	for _, ingredient := range topping.Ingredients {
		ingredients = append(ingredients, &pb.Ingredient{
			Name:   ingredient.Name,
			Amount: ingredient.Amount,
		})
	}

	return &pb.Topping{
		Name:          topping.Name,
		ReferenceArea: topping.ReferenceArea,
		Ingredients:   ingredients,
	}
}

func toProtoSteps(steps domain.Steps) *pb.Steps {
	protoSteps := make([]*pb.Step, 0, len(steps.Steps))
	for _, step := range steps.Steps {
		protoSteps = append(protoSteps, &pb.Step{
			Id:          int32(step.Id),
			StepNumber:  int32(step.StepNumber),
			Description: step.Description,
		})
	}

	return &pb.Steps{
		RecipeId: int32(steps.RecipeId),
		Steps:    protoSteps,
	}
}

func toProtoPans(pans domain.Pans) *pb.Pans {
	protoPans := make([]*pb.Pan, 0, len(pans.Pans))
	for _, pan := range pans.Pans {
		protoPans = append(protoPans, &pb.Pan{
			Shape: pan.Shape,
			Measures: &pb.Measures{
				Diameter: toProtoInt32Pointer(pan.Measures.Diameter),
				Edge:     toProtoInt32Pointer(pan.Measures.Edge),
				Width:    toProtoInt32Pointer(pan.Measures.Width),
				Length:   toProtoInt32Pointer(pan.Measures.Length),
			},
			Name: pan.Name,
			Area: pan.Area,
		})
	}

	return &pb.Pans{
		Pans:      protoPans,
		TotalArea: pans.TotalArea,
	}
}

func toDomainRecipeAggregate(protoAggregate *pb.RecipeAggregate) *domain.RecipeAggregate {
	if protoAggregate == nil {
		return nil
	}

	return &domain.RecipeAggregate{
		Recipe:           toDomainRecipe(protoAggregate.Recipe),
		SplitIngredients: toDomainSplitIngredients(protoAggregate.SplitIngredients),
	}
}

func toDomainRecipe(protoRecipe *pb.Recipe) domain.Recipe {
	if protoRecipe == nil {
		return domain.Recipe{}
	}

	recipeUuid, _ := uuid.Parse(protoRecipe.Uuid)
	return domain.Recipe{
		Id:          int(protoRecipe.Id),
		Uuid:        recipeUuid,
		Name:        protoRecipe.Name,
		Description: protoRecipe.Description,
		Author:      protoRecipe.Author,
		Dough:       toDomainDough(protoRecipe.Dough),
		Topping:     toDomainTopping(protoRecipe.Topping),
		Steps:       toDomainSteps(protoRecipe.Steps),
	}
}

func toDomainDough(protoDough *pb.Dough) domain.Dough {
	if protoDough == nil {
		return domain.Dough{}
	}

	ingredients := make([]domain.Ingredient, 0, len(protoDough.Ingredients))
	for _, ingredient := range protoDough.Ingredients {
		ingredients = append(ingredients, domain.Ingredient{
			Name:   ingredient.Name,
			Amount: ingredient.Amount,
		})
	}

	return domain.Dough{
		Name:             protoDough.Name,
		PercentVariation: protoDough.PercentVariation,
		Ingredients:      ingredients,
	}
}

func toDomainTopping(protoTopping *pb.Topping) domain.Topping {
	if protoTopping == nil {
		return domain.Topping{}
	}

	ingredients := make([]domain.Ingredient, 0, len(protoTopping.Ingredients))
	for _, ingredient := range protoTopping.Ingredients {
		ingredients = append(ingredients, domain.Ingredient{
			Name:   ingredient.Name,
			Amount: ingredient.Amount,
		})
	}

	return domain.Topping{
		Name:          protoTopping.Name,
		ReferenceArea: protoTopping.ReferenceArea,
		Ingredients:   ingredients,
	}
}

func toDomainSteps(protoSteps *pb.Steps) domain.Steps {
	if protoSteps == nil {
		return domain.Steps{}
	}

	steps := make([]domain.Step, 0, len(protoSteps.Steps))
	for _, step := range protoSteps.Steps {
		steps = append(steps, domain.Step{
			Id:          int(step.Id),
			StepNumber:  int(step.StepNumber),
			Description: step.Description,
		})
	}

	return domain.Steps{
		RecipeId: int(protoSteps.RecipeId),
		Steps:    steps,
	}
}

func toDomainSplitIngredients(protoSplitIngredients *pb.SplitIngredients) domain.SplitIngredients {
	if protoSplitIngredients == nil {
		return domain.SplitIngredients{}
	}

	splitDough := make([]domain.Dough, 0, len(protoSplitIngredients.SplitDough))
	for _, dough := range protoSplitIngredients.SplitDough {
		splitDough = append(splitDough, toDomainDough(dough))
	}

	splitTopping := make([]domain.Topping, 0, len(protoSplitIngredients.SplitTopping))
	for _, topping := range protoSplitIngredients.SplitTopping {
		splitTopping = append(splitTopping, toDomainTopping(topping))
	}

	return domain.SplitIngredients{
		SplitDough:   splitDough,
		SplitTopping: splitTopping,
	}
}

func toProtoInt32Pointer(value *int) *int32 {
	if value == nil {
		return nil
	}
	val := int32(*value)
	return &val
}
