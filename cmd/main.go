package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/cfioretti/recipe-manager/configs"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/application"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/infrastructure/grpc/client"
	httpHandlers "github.com/cfioretti/recipe-manager/internal/recipe-manager/infrastructure/http"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/infrastructure/logging"
	prometheusMetrics "github.com/cfioretti/recipe-manager/internal/recipe-manager/infrastructure/metrics"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/infrastructure/mysql"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/infrastructure/tracing"
	apihttp "github.com/cfioretti/recipe-manager/internal/recipe-manager/interfaces/api/http"
)

const (
	serviceName = "recipe-manager"
	version     = "1.0.0"
)

var logger *logging.Logger

func main() {
	logger = logging.NewLogger(serviceName, version)

	ctx := context.Background()
	logger.WithContext(ctx).Info("Starting recipe-manager service", logging.ServiceNameKey, serviceName)

	if err := tracing.InitTracing(nil); err != nil {
		logger.WithError(err).Fatal("Failed to initialize tracing")
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tracing.ShutdownTracing(ctx); err != nil {
			logger.WithError(err).Error("Failed to shutdown tracing")
		}
	}()

	if err := setEnvConfigs(); err != nil {
		logger.WithError(err).Fatal("Failed to load configuration")
	}

	db, err := loadDBConfig()
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.WithError(err).Error("Failed to close database connection")
		}
	}()

	prometheusMetrics := prometheusMetrics.NewPrometheusMetrics()

	router := setupRouter(db, prometheusMetrics)

	startServerWithGracefulShutdown(router)
}

func setupRouter(db *sql.DB, metrics *prometheusMetrics.PrometheusMetrics) *gin.Engine {
	calculatorService, err := initializeCalculatorService()
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize calculator service")
	}

	balancerService, err := initializeBalancerService()
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize balancer service")
	}

	recipeHandler := apihttp.NewRecipeHandler(
		application.NewRecipeService(
			mysql.NewMySqlRecipeRepository(db),
			calculatorService,
			balancerService,
		),
	)

	router := gin.New()

	router.Use(otelgin.Middleware(serviceName))
	router.Use(logger.GinMiddleware())

	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.WithContext(c.Request.Context()).WithFields(map[string]interface{}{
			"panic": recovered,
		}).Error("Panic recovered")
		c.AbortWithStatus(500)
	}))

	router.Use(corsMiddleware())

	metricsHandler := httpHandlers.NewMetricsHandler()
	metricsHandler.RegisterRoutes(router)

	healthHandler := httpHandlers.NewHealthHandler()
	healthHandler.RegisterRoutes(router)

	router.POST("/recipes/:uuid/aggregate", recipeHandler.RetrieveRecipeAggregate)

	return router
}

func startServerWithGracefulShutdown(router *gin.Engine) {
	port := viper.GetInt("server.port")
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	go func() {
		logger.WithField("port", port).Info("Starting HTTP server")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("Server forced to shutdown")
	}

	logger.Info("Server stopped")
}

func corsMiddleware() gin.HandlerFunc {
	feLocalHost := viper.GetString("local.fe-host")
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", feLocalHost)
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Correlation-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func initializeCalculatorService() (application.CalculatorService, error) {
	calculatorClient, err := loadCalculatorGrpcClient()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize calculator service: %w", err)
	}

	logger.Info("Calculator service initialized successfully")
	return application.NewRemoteDoughCalculatorService(calculatorClient), nil
}

func loadCalculatorGrpcClient() (*client.CalculatorClient, error) {
	config := configs.LoadCalculatorGRPCConfig()
	calculatorClient, err := client.NewDoughCalculatorClient(config.Address, config.Timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to create Calculator gRPC client: %w", err)
	}

	logger.WithField("address", config.Address).Info("Calculator gRPC client created")
	return calculatorClient, nil
}

func initializeBalancerService() (application.BalancerService, error) {
	balancerClient, err := loadBalancerGrpcClient()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize balancer service: %w", err)
	}

	logger.Info("Balancer service initialized successfully")
	return application.NewRemoteIngredientsBalancerService(balancerClient), nil
}

func loadBalancerGrpcClient() (*client.IngredientsBalancerClient, error) {
	config := configs.LoadBalancerGRPCConfig()
	balancerClient, err := client.NewIngredientsBalancerClient(config.Address, config.Timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to create Balancer gRPC client: %w", err)
	}

	logger.WithField("address", config.Address).Info("Balancer gRPC client created")
	return balancerClient, nil
}

func loadDBConfig() (*sql.DB, error) {
	config := configs.NewDBConfig()
	db, err := newDBConnection(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	logger.WithField("database", config.DBName).Info("Database connection established")
	return db, nil
}

func newDBConnection(config *configs.DBConfig) (*sql.DB, error) {
	db, err := sql.Open("mysql", config.DSN())
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	return db, nil
}

func setEnvConfigs() error {
	configName := os.Getenv("CONFIG_NAME")
	if configName == "" {
		configName = "props"
	}
	viper.SetConfigName(configName)
	viper.SetConfigType("yml")
	viper.AddConfigPath("configs/")
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read properties config: %w", err)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	logger.WithField("config", configName).Info("Configuration loaded")
	return nil
}
