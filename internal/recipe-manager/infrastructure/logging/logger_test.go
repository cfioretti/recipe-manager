package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		version     string
		logLevel    string
	}{
		{
			name:        "creates logger with default info level",
			serviceName: "test-service",
			version:     "1.0.0",
			logLevel:    "",
		},
		{
			name:        "creates logger with debug level",
			serviceName: "test-service",
			version:     "1.0.0",
			logLevel:    "debug",
		},
		{
			name:        "creates logger with error level",
			serviceName: "test-service",
			version:     "1.0.0",
			logLevel:    "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.logLevel != "" {
				t.Setenv("LOG_LEVEL", tt.logLevel)
			}

			logger := NewLogger(tt.serviceName, tt.version)

			assert.NotNil(t, logger)
			assert.Equal(t, tt.serviceName, logger.serviceName)
			assert.Equal(t, tt.version, logger.version)

			_, isJSON := logger.Formatter.(*logrus.JSONFormatter)
			assert.True(t, isJSON, "Logger should use JSON formatter")

			expectedLevel := logrus.InfoLevel
			if tt.logLevel != "" {
				var err error
				expectedLevel, err = logrus.ParseLevel(tt.logLevel)
				require.NoError(t, err)
			}
			assert.Equal(t, expectedLevel, logger.Level)
		})
	}
}

func TestLogger_WithContext(t *testing.T) {
	logger := NewLogger("test-service", "1.0.0")

	tests := []struct {
		name           string
		ctx            context.Context
		expectedFields map[string]interface{}
	}{
		{
			name: "context without correlation ID",
			ctx:  context.Background(),
			expectedFields: map[string]interface{}{
				ServiceNameKey: "test-service",
				VersionKey:     "1.0.0",
			},
		},
		{
			name: "context with correlation ID",
			ctx:  NewContextWithCorrelationID(context.Background(), "test-correlation-id"),
			expectedFields: map[string]interface{}{
				ServiceNameKey:   "test-service",
				VersionKey:       "1.0.0",
				CorrelationIDKey: "test-correlation-id",
			},
		},
		{
			name: "context with gRPC metadata",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs(
				CorrelationIDMetadataKey, "grpc-correlation-id",
			)),
			expectedFields: map[string]interface{}{
				ServiceNameKey:   "test-service",
				VersionKey:       "1.0.0",
				CorrelationIDKey: "grpc-correlation-id",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger.SetOutput(&buf)

			entry := logger.WithContext(tt.ctx)
			entry.Info("test message")

			var logData map[string]interface{}
			err := json.Unmarshal(buf.Bytes(), &logData)
			require.NoError(t, err)

			for key, expectedValue := range tt.expectedFields {
				assert.Equal(t, expectedValue, logData[key], "Field %s should match", key)
			}

			assert.Equal(t, "test message", logData["message"])
			assert.Equal(t, "info", logData["level"])
			assert.NotEmpty(t, logData["timestamp"])
		})
	}
}

func TestGetCorrelationID(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected string
	}{
		{
			name:     "generates new UUID for empty context",
			ctx:      context.Background(),
			expected: "",
		},
		{
			name:     "returns existing correlation ID from context",
			ctx:      NewContextWithCorrelationID(context.Background(), "existing-id"),
			expected: "existing-id",
		},
		{
			name: "returns correlation ID from gRPC metadata",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs(
				CorrelationIDMetadataKey, "grpc-id",
			)),
			expected: "grpc-id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetCorrelationID(tt.ctx)

			if tt.expected == "" {
				_, err := uuid.Parse(result)
				assert.NoError(t, err, "Should generate a valid UUID")
			} else {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestGRPCUnaryInterceptor(t *testing.T) {
	logger := NewLogger("test-service", "1.0.0")
	var buf bytes.Buffer
	logger.SetOutput(&buf)

	// Mock handler
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "response", nil
	}

	interceptor := logger.GRPCUnaryInterceptor()

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		CorrelationIDMetadataKey, "test-grpc-id",
	))

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/Method",
	}

	resp, err := interceptor(ctx, "request", info, handler)

	assert.NoError(t, err)
	assert.Equal(t, "response", resp)

	logs := strings.Split(strings.TrimSpace(buf.String()), "\n")
	require.Len(t, logs, 2, "Should have start and completion logs")

	var startLog map[string]interface{}
	err = json.Unmarshal([]byte(logs[0]), &startLog)
	require.NoError(t, err)

	var completionLog map[string]interface{}
	err = json.Unmarshal([]byte(logs[1]), &completionLog)
	require.NoError(t, err)

	assert.Equal(t, "gRPC request started", startLog["message"])
	assert.Equal(t, "/test.Service/Method", startLog["method"])
	assert.Equal(t, "grpc_request", startLog["type"])
	assert.Equal(t, "test-grpc-id", startLog[CorrelationIDKey])

	assert.Equal(t, "gRPC request completed", completionLog["message"])
	assert.Equal(t, "/test.Service/Method", completionLog["method"])
	assert.Equal(t, "grpc_request", completionLog["type"])
	assert.Equal(t, "test-grpc-id", completionLog[CorrelationIDKey])
	assert.NotNil(t, completionLog["duration_ms"])
}

func TestGinMiddleware(t *testing.T) {
	logger := NewLogger("test-service", "1.0.0")
	var buf bytes.Buffer
	logger.SetOutput(&buf)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(logger.GinMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	w := performRequest(router, "GET", "/test", nil)

	assert.Equal(t, 200, w.Code)

	logs := strings.Split(strings.TrimSpace(buf.String()), "\n")
	require.Len(t, logs, 2, "Should have start and completion logs")

	var startLog map[string]interface{}
	err := json.Unmarshal([]byte(logs[0]), &startLog)
	require.NoError(t, err)

	var completionLog map[string]interface{}
	err = json.Unmarshal([]byte(logs[1]), &completionLog)
	require.NoError(t, err)

	assert.Equal(t, "Request started", startLog["message"])
	assert.Equal(t, "GET", startLog["method"])
	assert.Equal(t, "/test", startLog["path"])
	assert.NotEmpty(t, startLog[CorrelationIDKey])

	assert.Equal(t, "Request completed", completionLog["message"])
	assert.Equal(t, "GET", completionLog["method"])
	assert.Equal(t, "/test", completionLog["path"])
	assert.Equal(t, float64(200), completionLog["status_code"])
	assert.NotNil(t, completionLog["duration_ms"])

	assert.Equal(t, startLog[CorrelationIDKey], completionLog[CorrelationIDKey])
}

func TestLoggerJSONFormat(t *testing.T) {
	logger := NewLogger("test-service", "1.0.0")
	var buf bytes.Buffer
	logger.SetOutput(&buf)

	logger.WithContext(context.Background()).Info("info message")
	logger.WithContext(context.Background()).Warn("warn message")
	logger.WithField("custom_field", "custom_value").Error("error message")

	logs := strings.Split(strings.TrimSpace(buf.String()), "\n")
	require.Len(t, logs, 3)

	for i, logLine := range logs {
		var logData map[string]interface{}
		err := json.Unmarshal([]byte(logLine), &logData)
		require.NoError(t, err, "Log line %d should be valid JSON", i)

		assert.NotEmpty(t, logData["timestamp"])
		assert.NotEmpty(t, logData["level"])
		assert.NotEmpty(t, logData["message"])
		assert.Equal(t, "test-service", logData[ServiceNameKey])
		assert.Equal(t, "1.0.0", logData[VersionKey])

		_, err = time.Parse(time.RFC3339, logData["timestamp"].(string))
		assert.NoError(t, err, "Timestamp should be in RFC3339 format")
	}

	var errorLog map[string]interface{}
	err := json.Unmarshal([]byte(logs[2]), &errorLog)
	require.NoError(t, err)
	assert.Equal(t, "custom_value", errorLog["custom_field"])
}

func performRequest(r *gin.Engine, method, path string, body []byte) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, bytes.NewBuffer(body))
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
