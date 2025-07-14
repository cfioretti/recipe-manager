package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	domainMetrics "github.com/cfioretti/recipe-manager/internal/recipe-manager/domain/metrics"
)

type PrometheusMetrics struct {
	// Recipe Operations
	recipeRetrievalsTotal      *prometheus.CounterVec
	recipeRetrievalDuration    prometheus.Histogram
	recipeRetrievalErrorsTotal *prometheus.CounterVec

	// Recipe Aggregation
	recipeAggregationsTotal      *prometheus.CounterVec
	recipeAggregationDuration    prometheus.Histogram
	recipeAggregationErrorsTotal *prometheus.CounterVec

	// External Service Calls
	calculatorServiceCallsTotal *prometheus.CounterVec
	calculatorServiceDuration   prometheus.Histogram
	balancerServiceCallsTotal   *prometheus.CounterVec
	balancerServiceDuration     prometheus.Histogram

	// Database Operations
	databaseOperationsTotal   *prometheus.CounterVec
	databaseOperationDuration *prometheus.HistogramVec

	// HTTP Request Metrics
	httpRequestsTotal     *prometheus.CounterVec
	httpRequestDuration   *prometheus.HistogramVec
	activeHTTPConnections prometheus.Gauge

	// Business Domain Metrics
	recipesByAuthor      *prometheus.CounterVec
	recipeComplexity     prometheus.Histogram
	pansSizes            *prometheus.CounterVec
	ingredientVariations prometheus.Histogram
}

func NewPrometheusMetrics() *PrometheusMetrics {
	return NewPrometheusMetricsWithRegistry(prometheus.DefaultRegisterer)
}

func NewPrometheusMetricsWithRegistry(reg prometheus.Registerer) *PrometheusMetrics {
	factory := promauto.With(reg)
	return &PrometheusMetrics{
		recipeRetrievalsTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "recipe_manager_recipe_retrievals_total",
				Help: "Total number of recipe retrievals by UUID",
			},
			[]string{"recipe_uuid"},
		),
		recipeRetrievalDuration: factory.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "recipe_manager_recipe_retrieval_duration_seconds",
				Help:    "Duration of recipe retrieval operations",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0},
			},
		),
		recipeRetrievalErrorsTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "recipe_manager_recipe_retrieval_errors_total",
				Help: "Total number of recipe retrieval errors by type",
			},
			[]string{"error_type"},
		),

		recipeAggregationsTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "recipe_manager_recipe_aggregations_total",
				Help: "Total number of recipe aggregations by type",
			},
			[]string{"recipe_type"},
		),
		recipeAggregationDuration: factory.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "recipe_manager_recipe_aggregation_duration_seconds",
				Help:    "Duration of recipe aggregation operations",
				Buckets: []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0, 30.0, 60.0},
			},
		),
		recipeAggregationErrorsTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "recipe_manager_recipe_aggregation_errors_total",
				Help: "Total number of recipe aggregation errors by type",
			},
			[]string{"error_type"},
		),

		calculatorServiceCallsTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "recipe_manager_calculator_service_calls_total",
				Help: "Total number of calculator service calls",
			},
			[]string{"success"},
		),
		calculatorServiceDuration: factory.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "recipe_manager_calculator_service_duration_seconds",
				Help:    "Duration of calculator service calls",
				Buckets: []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0, 30.0},
			},
		),
		balancerServiceCallsTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "recipe_manager_balancer_service_calls_total",
				Help: "Total number of balancer service calls",
			},
			[]string{"success"},
		),
		balancerServiceDuration: factory.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "recipe_manager_balancer_service_duration_seconds",
				Help:    "Duration of balancer service calls",
				Buckets: []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0, 30.0},
			},
		),

		databaseOperationsTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "recipe_manager_database_operations_total",
				Help: "Total number of database operations",
			},
			[]string{"operation", "success"},
		),
		databaseOperationDuration: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "recipe_manager_database_operation_duration_seconds",
				Help:    "Duration of database operations",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0},
			},
			[]string{"operation"},
		),

		httpRequestsTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "recipe_manager_http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status_code"},
		),
		httpRequestDuration: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "recipe_manager_http_request_duration_seconds",
				Help:    "Duration of HTTP requests",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0},
			},
			[]string{"method", "endpoint"},
		),
		activeHTTPConnections: factory.NewGauge(
			prometheus.GaugeOpts{
				Name: "recipe_manager_active_http_connections",
				Help: "Number of active HTTP connections",
			},
		),

		recipesByAuthor: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "recipe_manager_recipes_by_author_total",
				Help: "Total number of recipes by author",
			},
			[]string{"author"},
		),
		recipeComplexity: factory.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "recipe_manager_recipe_complexity",
				Help:    "Recipe complexity score",
				Buckets: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
		),
		pansSizes: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "recipe_manager_pans_sizes_total",
				Help: "Total number of pans by size",
			},
			[]string{"pan_size"},
		),
		ingredientVariations: factory.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "recipe_manager_ingredient_variations",
				Help:    "Number of ingredient variations per recipe",
				Buckets: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 15, 20, 25, 30},
			},
		),
	}
}

func (p *PrometheusMetrics) IncrementRecipeRetrievals(recipeUuid string) {
	p.recipeRetrievalsTotal.WithLabelValues(recipeUuid).Inc()
}

func (p *PrometheusMetrics) RecordRecipeRetrievalDuration(duration time.Duration) {
	p.recipeRetrievalDuration.Observe(duration.Seconds())
}

func (p *PrometheusMetrics) IncrementRecipeRetrievalErrors(errorType string) {
	p.recipeRetrievalErrorsTotal.WithLabelValues(errorType).Inc()
}

func (p *PrometheusMetrics) IncrementRecipeAggregations(recipeType string) {
	p.recipeAggregationsTotal.WithLabelValues(recipeType).Inc()
}

func (p *PrometheusMetrics) RecordRecipeAggregationDuration(duration time.Duration) {
	p.recipeAggregationDuration.Observe(duration.Seconds())
}

func (p *PrometheusMetrics) IncrementRecipeAggregationErrors(errorType string) {
	p.recipeAggregationErrorsTotal.WithLabelValues(errorType).Inc()
}

func (p *PrometheusMetrics) IncrementCalculatorServiceCalls(success bool) {
	p.calculatorServiceCallsTotal.WithLabelValues(boolToString(success)).Inc()
}

func (p *PrometheusMetrics) RecordCalculatorServiceDuration(duration time.Duration) {
	p.calculatorServiceDuration.Observe(duration.Seconds())
}

func (p *PrometheusMetrics) IncrementBalancerServiceCalls(success bool) {
	p.balancerServiceCallsTotal.WithLabelValues(boolToString(success)).Inc()
}

func (p *PrometheusMetrics) RecordBalancerServiceDuration(duration time.Duration) {
	p.balancerServiceDuration.Observe(duration.Seconds())
}

func (p *PrometheusMetrics) IncrementDatabaseOperations(operation string, success bool) {
	p.databaseOperationsTotal.WithLabelValues(operation, boolToString(success)).Inc()
}

func (p *PrometheusMetrics) RecordDatabaseOperationDuration(operation string, duration time.Duration) {
	p.databaseOperationDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

func (p *PrometheusMetrics) IncrementHTTPRequests(method string, endpoint string, statusCode int) {
	p.httpRequestsTotal.WithLabelValues(method, endpoint, strconv.Itoa(statusCode)).Inc()
}

func (p *PrometheusMetrics) RecordHTTPRequestDuration(method string, endpoint string, duration time.Duration) {
	p.httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
}

func (p *PrometheusMetrics) SetActiveHTTPConnections(count int) {
	p.activeHTTPConnections.Set(float64(count))
}

func (p *PrometheusMetrics) IncrementRecipesByAuthor(author string) {
	p.recipesByAuthor.WithLabelValues(author).Inc()
}

func (p *PrometheusMetrics) RecordRecipeComplexity(complexity int) {
	p.recipeComplexity.Observe(float64(complexity))
}

func (p *PrometheusMetrics) IncrementPansSizes(panSize string) {
	p.pansSizes.WithLabelValues(panSize).Inc()
}

func (p *PrometheusMetrics) RecordIngredientVariations(variationCount int) {
	p.ingredientVariations.Observe(float64(variationCount))
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

var _ domainMetrics.RecipeMetrics = (*PrometheusMetrics)(nil)
