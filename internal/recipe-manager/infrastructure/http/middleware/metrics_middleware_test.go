package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"

	infraMetrics "github.com/cfioretti/recipe-manager/internal/recipe-manager/infrastructure/metrics"
)

type MockRecipeMetrics struct {
	httpRequests               map[string]int
	httpRequestDurations       map[string][]time.Duration
	activeHTTPConnections      int
	recipeRetrievalDurations   []time.Duration
	recipeRetrievalErrors      map[string]int
	recipeAggregationDurations []time.Duration
	recipeAggregationErrors    map[string]int
}

func NewMockRecipeMetrics() *MockRecipeMetrics {
	return &MockRecipeMetrics{
		httpRequests:            make(map[string]int),
		httpRequestDurations:    make(map[string][]time.Duration),
		recipeRetrievalErrors:   make(map[string]int),
		recipeAggregationErrors: make(map[string]int),
	}
}

func (m *MockRecipeMetrics) IncrementRecipeRetrievals(recipeUuid string) {}
func (m *MockRecipeMetrics) RecordRecipeRetrievalDuration(duration time.Duration) {
	m.recipeRetrievalDurations = append(m.recipeRetrievalDurations, duration)
}
func (m *MockRecipeMetrics) IncrementRecipeRetrievalErrors(errorType string) {
	m.recipeRetrievalErrors[errorType]++
}
func (m *MockRecipeMetrics) IncrementRecipeAggregations(recipeType string) {}
func (m *MockRecipeMetrics) RecordRecipeAggregationDuration(duration time.Duration) {
	m.recipeAggregationDurations = append(m.recipeAggregationDurations, duration)
}
func (m *MockRecipeMetrics) IncrementRecipeAggregationErrors(errorType string) {
	m.recipeAggregationErrors[errorType]++
}
func (m *MockRecipeMetrics) IncrementCalculatorServiceCalls(success bool)               {}
func (m *MockRecipeMetrics) RecordCalculatorServiceDuration(duration time.Duration)     {}
func (m *MockRecipeMetrics) IncrementBalancerServiceCalls(success bool)                 {}
func (m *MockRecipeMetrics) RecordBalancerServiceDuration(duration time.Duration)       {}
func (m *MockRecipeMetrics) IncrementDatabaseOperations(operation string, success bool) {}
func (m *MockRecipeMetrics) RecordDatabaseOperationDuration(operation string, duration time.Duration) {
}
func (m *MockRecipeMetrics) IncrementHTTPRequests(method, endpoint string, statusCode int) {
	key := method + ":" + endpoint
	m.httpRequests[key]++
}
func (m *MockRecipeMetrics) RecordHTTPRequestDuration(method, endpoint string, duration time.Duration) {
	key := method + ":" + endpoint
	m.httpRequestDurations[key] = append(m.httpRequestDurations[key], duration)
}
func (m *MockRecipeMetrics) SetActiveHTTPConnections(count int) {
	m.activeHTTPConnections = count
}
func (m *MockRecipeMetrics) IncrementRecipesByAuthor(author string)        {}
func (m *MockRecipeMetrics) RecordRecipeComplexity(complexity int)         {}
func (m *MockRecipeMetrics) IncrementPansSizes(panSize string)             {}
func (m *MockRecipeMetrics) RecordIngredientVariations(variationCount int) {}

func TestMetricsMiddleware_HTTPMetricsMiddleware_Success(t *testing.T) {
	mockDomainMetrics := NewMockRecipeMetrics()
	registry := prometheus.NewRegistry()
	prometheusMetrics := infraMetrics.NewPrometheusMetricsWithRegistry(registry)

	middleware := NewMetricsMiddleware(mockDomainMetrics, prometheusMetrics)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.HTTPMetricsMiddleware())

	router.GET("/api/v1/recipes/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{"id": c.Param("id")})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/recipes/123", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if mockDomainMetrics.httpRequests["GET:/api/v1/recipes/:id"] != 1 {
		t.Errorf("Expected 1 HTTP request, got %d", mockDomainMetrics.httpRequests["GET:/api/v1/recipes/:id"])
	}

	if len(mockDomainMetrics.httpRequestDurations["GET:/api/v1/recipes/:id"]) != 1 {
		t.Errorf("Expected 1 HTTP request duration, got %d", len(mockDomainMetrics.httpRequestDurations["GET:/api/v1/recipes/:id"]))
	}

	if len(mockDomainMetrics.recipeRetrievalDurations) != 1 {
		t.Errorf("Expected 1 recipe retrieval duration, got %d", len(mockDomainMetrics.recipeRetrievalDurations))
	}
}

func TestMetricsMiddleware_HTTPMetricsMiddleware_Error(t *testing.T) {
	mockDomainMetrics := NewMockRecipeMetrics()
	registry := prometheus.NewRegistry()
	prometheusMetrics := infraMetrics.NewPrometheusMetricsWithRegistry(registry)

	middleware := NewMetricsMiddleware(mockDomainMetrics, prometheusMetrics)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.HTTPMetricsMiddleware())

	router.GET("/api/v1/recipes/:id", func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "not found"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/recipes/123", nil)
	router.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	if mockDomainMetrics.httpRequests["GET:/api/v1/recipes/:id"] != 1 {
		t.Errorf("Expected 1 HTTP request, got %d", mockDomainMetrics.httpRequests["GET:/api/v1/recipes/:id"])
	}

	if mockDomainMetrics.recipeRetrievalErrors["client_error"] != 1 {
		t.Errorf("Expected 1 recipe retrieval error, got %d", mockDomainMetrics.recipeRetrievalErrors["client_error"])
	}

	if len(mockDomainMetrics.recipeRetrievalDurations) != 0 {
		t.Errorf("Expected 0 recipe retrieval durations for error, got %d", len(mockDomainMetrics.recipeRetrievalDurations))
	}
}

func TestIsRecipeEndpoint(t *testing.T) {
	tests := []struct {
		endpoint string
		expected bool
	}{
		{"/api/v1/recipes", true},
		{"/api/v1/recipes/:id", true},
		{"/api/v1/recipes/search", true},
		{"/api/v1/recipes/aggregate", true},
		{"/api/v1/recipes/author/:author", true},
		{"/health", false},
		{"/metrics", false},
		{"/api/v1/other", false},
	}

	for _, test := range tests {
		result := isRecipeEndpoint(test.endpoint)
		if result != test.expected {
			t.Errorf("isRecipeEndpoint(%s) = %t, expected %t", test.endpoint, result, test.expected)
		}
	}
}

func TestGetErrorType(t *testing.T) {
	tests := []struct {
		statusCode int
		expected   string
	}{
		{200, "unknown_error"},
		{404, "client_error"},
		{400, "client_error"},
		{500, "server_error"},
		{503, "server_error"},
	}

	for _, test := range tests {
		result := getErrorType(test.statusCode)
		if result != test.expected {
			t.Errorf("getErrorType(%d) = %s, expected %s", test.statusCode, result, test.expected)
		}
	}
}
