package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	domainMetrics "github.com/cfioretti/recipe-manager/internal/recipe-manager/domain/metrics"
	infraMetrics "github.com/cfioretti/recipe-manager/internal/recipe-manager/infrastructure/metrics"
)

type MetricsMiddleware struct {
	domainMetrics     domainMetrics.RecipeMetrics
	prometheusMetrics *infraMetrics.PrometheusMetrics
}

func NewMetricsMiddleware(
	domainMetrics domainMetrics.RecipeMetrics,
	prometheusMetrics *infraMetrics.PrometheusMetrics,
) *MetricsMiddleware {
	return &MetricsMiddleware{
		domainMetrics:     domainMetrics,
		prometheusMetrics: prometheusMetrics,
	}
}

func (m *MetricsMiddleware) HTTPMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		m.domainMetrics.SetActiveHTTPConnections(getCurrentActiveConnections() + 1)

		c.Next()

		duration := time.Since(start)
		method := c.Request.Method
		endpoint := c.FullPath()
		statusCode := c.Writer.Status()

		m.domainMetrics.IncrementHTTPRequests(method, endpoint, statusCode)
		m.domainMetrics.RecordHTTPRequestDuration(method, endpoint, duration)

		if isRecipeEndpoint(endpoint) {
			m.recordRecipeMetrics(method, endpoint, statusCode, duration)
		}

		m.domainMetrics.SetActiveHTTPConnections(getCurrentActiveConnections() - 1)
	}
}

func (m *MetricsMiddleware) recordRecipeMetrics(method, endpoint string, statusCode int, duration time.Duration) {
	switch {
	case isRecipeRetrievalEndpoint(endpoint):
		if statusCode >= 200 && statusCode < 300 {
			m.domainMetrics.RecordRecipeRetrievalDuration(duration)
		} else {
			m.domainMetrics.IncrementRecipeRetrievalErrors(getErrorType(statusCode))
		}
	case isRecipeAggregationEndpoint(endpoint):
		if statusCode >= 200 && statusCode < 300 {
			m.domainMetrics.RecordRecipeAggregationDuration(duration)
		} else {
			m.domainMetrics.IncrementRecipeAggregationErrors(getErrorType(statusCode))
		}
	}
}

func isRecipeEndpoint(endpoint string) bool {
	recipeEndpoints := []string{
		"/api/v1/recipes",
		"/api/v1/recipes/:id",
		"/api/v1/recipes/search",
		"/api/v1/recipes/aggregate",
		"/api/v1/recipes/author/:author",
	}

	for _, recipeEndpoint := range recipeEndpoints {
		if endpoint == recipeEndpoint {
			return true
		}
	}
	return false
}

func isRecipeRetrievalEndpoint(endpoint string) bool {
	retrievalEndpoints := []string{
		"/api/v1/recipes/:id",
		"/api/v1/recipes/search",
		"/api/v1/recipes/author/:author",
	}

	for _, retrievalEndpoint := range retrievalEndpoints {
		if endpoint == retrievalEndpoint {
			return true
		}
	}
	return false
}

func isRecipeAggregationEndpoint(endpoint string) bool {
	return endpoint == "/api/v1/recipes/aggregate"
}

func getErrorType(statusCode int) string {
	switch {
	case statusCode >= 400 && statusCode < 500:
		return "client_error"
	case statusCode >= 500:
		return "server_error"
	default:
		return "unknown_error"
	}
}

func getCurrentActiveConnections() int {
	// In a real implementation, this would be tracked in a service
	// For now, we'll use a simple counter
	return 0
}
