package application_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	bdomain "github.com/cfioretti/recipe-manager/internal/ingredients-balancer/domain"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/application"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/infrastructure/grpc/client"
)

func TestTotalDoughWeightByPans(t *testing.T) {
	t.Skip("this test requires a running gRPC server")

	grpcClient, err := client.NewDoughCalculatorClient("localhost:50051", 5*time.Second)
	require.NoError(t, err)
	defer grpcClient.Close()

	service := application.NewRemoteDoughCalculatorService(grpcClient)

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
