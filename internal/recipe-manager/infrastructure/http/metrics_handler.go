package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsHandler struct {
	handler http.Handler
}

func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{
		handler: promhttp.Handler(),
	}
}

func (h *MetricsHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/metrics", h.handleMetrics)
}

func (h *MetricsHandler) handleMetrics(c *gin.Context) {
	h.handler.ServeHTTP(c.Writer, c.Request)
}

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", h.handleHealth)
}

func (h *HealthHandler) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "recipe-manager",
		"version": "1.0.0",
	})
}
