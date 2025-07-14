package metrics

import (
	"context"
	"time"
)

type RecipeMetrics interface {
	// Recipe Operations
	IncrementRecipeRetrievals(recipeUuid string)
	RecordRecipeRetrievalDuration(duration time.Duration)
	IncrementRecipeRetrievalErrors(errorType string)

	// Recipe Aggregation
	IncrementRecipeAggregations(recipeType string)
	RecordRecipeAggregationDuration(duration time.Duration)
	IncrementRecipeAggregationErrors(errorType string)

	// External Service Calls
	IncrementCalculatorServiceCalls(success bool)
	RecordCalculatorServiceDuration(duration time.Duration)
	IncrementBalancerServiceCalls(success bool)
	RecordBalancerServiceDuration(duration time.Duration)

	// Database Operations
	IncrementDatabaseOperations(operation string, success bool)
	RecordDatabaseOperationDuration(operation string, duration time.Duration)

	// HTTP Request Metrics
	IncrementHTTPRequests(method string, endpoint string, statusCode int)
	RecordHTTPRequestDuration(method string, endpoint string, duration time.Duration)
	SetActiveHTTPConnections(count int)

	// Business Domain Metrics
	IncrementRecipesByAuthor(author string)
	RecordRecipeComplexity(complexity int)
	IncrementPansSizes(panSize string)
	RecordIngredientVariations(variationCount int)
}

type RecipeOperationResult struct {
	Type         string
	Duration     time.Duration
	Success      bool
	ErrorType    string
	RecipeUuid   string
	RecipeAuthor string
	PansCount    int
	Complexity   int
}

type MetricsRecorder struct {
	metrics RecipeMetrics
}

func NewMetricsRecorder(metrics RecipeMetrics) *MetricsRecorder {
	return &MetricsRecorder{
		metrics: metrics,
	}
}

func (m *MetricsRecorder) RecordRecipeOperation(ctx context.Context, result RecipeOperationResult) {
	m.metrics.RecordRecipeRetrievalDuration(result.Duration)
	m.metrics.IncrementRecipeRetrievals(result.RecipeUuid)

	if result.Success {
		m.metrics.IncrementRecipeAggregations(result.Type)
		m.metrics.IncrementRecipesByAuthor(result.RecipeAuthor)
		m.metrics.RecordRecipeComplexity(result.Complexity)
		if result.PansCount > 0 {
			m.metrics.RecordIngredientVariations(result.PansCount)
		}
	} else {
		m.metrics.IncrementRecipeRetrievalErrors(result.ErrorType)
		m.metrics.IncrementRecipeAggregationErrors(result.ErrorType)
	}
}

func (m *MetricsRecorder) RecordCalculatorCall(ctx context.Context, duration time.Duration, success bool) {
	m.metrics.IncrementCalculatorServiceCalls(success)
	m.metrics.RecordCalculatorServiceDuration(duration)
}

func (m *MetricsRecorder) RecordBalancerCall(ctx context.Context, duration time.Duration, success bool) {
	m.metrics.IncrementBalancerServiceCalls(success)
	m.metrics.RecordBalancerServiceDuration(duration)
}

func (m *MetricsRecorder) RecordDatabaseOperation(ctx context.Context, operation string, duration time.Duration, success bool) {
	m.metrics.IncrementDatabaseOperations(operation, success)
	m.metrics.RecordDatabaseOperationDuration(operation, duration)
}
