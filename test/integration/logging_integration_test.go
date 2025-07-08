package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cfioretti/recipe-manager/internal/infrastructure/logging"
)

type MockGRPCService struct {
	logger *logging.Logger
	logBuf *bytes.Buffer
}

func NewMockGRPCService() *MockGRPCService {
	logger := logging.NewLogger("mock-grpc-service", "1.0.0")
	logBuf := &bytes.Buffer{}
	logger.SetOutput(logBuf)

	return &MockGRPCService{
		logger: logger,
		logBuf: logBuf,
	}
}

func (m *MockGRPCService) ProcessRequest(ctx context.Context, data string) (string, error) {
	m.logger.WithContext(ctx).WithField("request_data", data).Info("Processing gRPC request")

	time.Sleep(10 * time.Millisecond)

	m.logger.WithContext(ctx).WithField("response_data", "processed_"+data).Info("gRPC request completed")

	return "processed_" + data, nil
}

func (m *MockGRPCService) GetLogs() []map[string]interface{} {
	logLines := strings.Split(strings.TrimSpace(m.logBuf.String()), "\n")
	var logs []map[string]interface{}

	for _, line := range logLines {
		if line == "" {
			continue
		}
		var logData map[string]interface{}
		if err := json.Unmarshal([]byte(line), &logData); err == nil {
			logs = append(logs, logData)
		}
	}

	return logs
}

func TestCorrelationIDPropagation(t *testing.T) {
	httpLogger := logging.NewLogger("http-service", "1.0.0")
	httpLogBuf := &bytes.Buffer{}
	httpLogger.SetOutput(httpLogBuf)

	grpcService := NewMockGRPCService()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(httpLogger.GinMiddleware())

	router.POST("/process", func(c *gin.Context) {
		var requestBody map[string]string
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		ctx := c.Request.Context()
		correlationID := logging.GetCorrelationID(ctx)

		grpcCtx := logging.NewContextWithCorrelationID(context.Background(), correlationID)
		result, err := grpcService.ProcessRequest(grpcCtx, requestBody["data"])

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"result": result})
	})

	requestBody := `{"data": "test_data"}`
	req := httptest.NewRequest("POST", "/process", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Correlation-ID", "custom-correlation-id")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "processed_test_data", response["result"])

	httpLogs := parseLogsFromBuffer(httpLogBuf)
	require.Len(t, httpLogs, 2, "Should have HTTP start and completion logs")

	grpcLogs := grpcService.GetLogs()
	require.Len(t, grpcLogs, 2, "Should have gRPC start and completion logs")

	httpCorrelationID := httpLogs[0]["correlation_id"]
	grpcCorrelationID := grpcLogs[0]["correlation_id"]

	assert.Equal(t, "custom-correlation-id", httpCorrelationID)
	assert.Equal(t, "custom-correlation-id", grpcCorrelationID)

	for _, log := range httpLogs {
		assert.Equal(t, "custom-correlation-id", log["correlation_id"])
	}
	for _, log := range grpcLogs {
		assert.Equal(t, "custom-correlation-id", log["correlation_id"])
	}

	assert.Equal(t, "http-service", httpLogs[0]["service_name"])
	assert.Equal(t, "mock-grpc-service", grpcLogs[0]["service_name"])
}

func TestAutomaticCorrelationIDGeneration(t *testing.T) {
	logger := logging.NewLogger("test-service", "1.0.0")
	logBuf := &bytes.Buffer{}
	logger.SetOutput(logBuf)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(logger.GinMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	logs := parseLogsFromBuffer(logBuf)
	require.Len(t, logs, 2)

	correlationID, exists := logs[0]["correlation_id"].(string)
	require.True(t, exists, "Correlation ID should exist")
	assert.Len(t, correlationID, 36, "Should be UUID format")
	assert.Contains(t, correlationID, "-", "Should contain UUID hyphens")

	assert.Equal(t, logs[0]["correlation_id"], logs[1]["correlation_id"])
}

func TestLoggingPerformance(t *testing.T) {
	logger := logging.NewLogger("perf-test", "1.0.0")
	logBuf := &bytes.Buffer{}
	logger.SetOutput(logBuf)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(logger.GinMiddleware())
	router.GET("/perf", func(c *gin.Context) {
		time.Sleep(1 * time.Millisecond)
		c.JSON(200, gin.H{"status": "ok"})
	})

	start := time.Now()
	numRequests := 100

	for i := 0; i < numRequests; i++ {
		req := httptest.NewRequest("GET", "/perf", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
	}

	duration := time.Since(start)
	avgDuration := duration / time.Duration(numRequests)

	assert.Less(t, avgDuration, 5*time.Millisecond, "Logging overhead should be minimal")

	logs := parseLogsFromBuffer(logBuf)
	assert.Equal(t, numRequests*2, len(logs), "Should have start and completion logs for each request")
}

func TestErrorLogging(t *testing.T) {
	logger := logging.NewLogger("error-test", "1.0.0")
	logBuf := &bytes.Buffer{}
	logger.SetOutput(logBuf)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(logger.GinMiddleware())
	router.GET("/error", func(c *gin.Context) {
		logger.WithContext(c.Request.Context()).WithError(fmt.Errorf("test error")).Error("Request failed")
		c.JSON(500, gin.H{"error": "Internal server error"})
	})

	req := httptest.NewRequest("GET", "/error", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)

	logs := parseLogsFromBuffer(logBuf)
	require.Len(t, logs, 3, "Should have start, error, and completion logs")

	var errorLog map[string]interface{}
	for _, log := range logs {
		if log["level"] == "error" {
			errorLog = log
			break
		}
	}

	require.NotNil(t, errorLog, "Should have error log")
	assert.Equal(t, "Request failed", errorLog["message"])
	assert.Equal(t, "test error", errorLog["error"])
	assert.NotEmpty(t, errorLog["correlation_id"])
}

func parseLogsFromBuffer(buf *bytes.Buffer) []map[string]interface{} {
	logLines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	var logs []map[string]interface{}

	for _, line := range logLines {
		if line == "" {
			continue
		}
		var logData map[string]interface{}
		if err := json.Unmarshal([]byte(line), &logData); err == nil {
			logs = append(logs, logData)
		}
	}

	return logs
}
