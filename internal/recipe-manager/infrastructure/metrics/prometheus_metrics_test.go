package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrometheusMetrics(t *testing.T) {
	metrics := NewPrometheusMetrics()

	assert.NotNil(t, metrics)
	assert.NotNil(t, metrics.recipeRetrievalsTotal)
	assert.NotNil(t, metrics.recipeRetrievalDuration)
	assert.NotNil(t, metrics.httpRequestsTotal)
	assert.NotNil(t, metrics.databaseOperationsTotal)
	assert.NotNil(t, metrics.calculatorServiceCallsTotal)
	assert.NotNil(t, metrics.balancerServiceCallsTotal)

	metrics.IncrementRecipeRetrievals("test-uuid")
	metrics.IncrementRecipeRetrievalErrors("database_error")
	metrics.IncrementRecipeAggregations("napoletana")
	metrics.IncrementCalculatorServiceCalls(true)
	metrics.IncrementBalancerServiceCalls(false)
	metrics.IncrementDatabaseOperations("SELECT", true)
	metrics.IncrementHTTPRequests("POST", "/recipes/:uuid/aggregate", 200)
	metrics.SetActiveHTTPConnections(5)
	metrics.IncrementRecipesByAuthor("chef-mario")
	metrics.RecordRecipeComplexity(5)
	metrics.RecordIngredientVariations(10)

	assert.Equal(t, "true", boolToString(true))
	assert.Equal(t, "false", boolToString(false))
}
