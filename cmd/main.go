package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"

	"github.com/cfioretti/recipe-manager/configs"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/application"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/infrastructure/grpc/client"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/infrastructure/mysql"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/interfaces/api/http"
)

func main() {
	setEnvConfigs()
	db := loadDBConfig()

	router := setupRouter(db)
	startServer(router)
}

func setupRouter(db *sql.DB) *gin.Engine {
	calculatorService := initializeCalculatorService()
	balancerService := initializeBalancerService()
	recipeHandler := http.NewRecipeHandler(
		application.NewRecipeService(
			mysql.NewMySqlRecipeRepository(db),
			calculatorService,
			balancerService,
		),
	)

	router := gin.Default()
	router.Use(corsMiddleware())
	router.POST("/recipes/:uuid/aggregate", recipeHandler.RetrieveRecipeAggregate)
	return router
}

func startServer(router *gin.Engine) {
	port := viper.GetInt("server.port")
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		panic(fmt.Errorf("failed to start server: %v", err))
	}
}

func corsMiddleware() gin.HandlerFunc {
	feLocalHost := viper.GetString("local.fe-host")
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", feLocalHost)
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func initializeCalculatorService() application.CalculatorService {
	calculatorClient, err := loadCalculatorGrpcClient()
	if err != nil {
		panic(fmt.Errorf("failed to initialize calculator service: %v", err))
	}
	return application.NewRemoteDoughCalculatorService(calculatorClient)
}

func loadCalculatorGrpcClient() (*client.CalculatorClient, error) {
	config := configs.LoadCalculatorGRPCConfig()
	calculatorClient, err := client.NewDoughCalculatorClient(config.Address, config.Timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to create Calculator gRPC calculatorClient: %w", err)
	}
	return calculatorClient, nil
}

func initializeBalancerService() application.BalancerService {
	balancerClient, err := loadBalancerGrpcClient()
	if err != nil {
		panic(fmt.Errorf("failed to initialize balancer service: %v", err))
	}
	return application.NewRemoteIngredientsBalancerService(balancerClient)
}

func loadBalancerGrpcClient() (*client.IngredientsBalancerClient, error) {
	config := configs.LoadBalancerGRPCConfig()
	balancerClient, err := client.NewIngredientsBalancerClient(config.Address, config.Timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to create Balancer gRPC balancerClient: %w", err)
	}
	return balancerClient, nil
}

func loadDBConfig() *sql.DB {
	config := configs.NewDBConfig()
	db, err := newDBConnection(config)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %v", err))
	}
	return db
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

func setEnvConfigs() {
	configName := os.Getenv("CONFIG_NAME")
	if configName == "" {
		configName = "props"
	}
	viper.SetConfigName(configName)
	viper.SetConfigType("yml")
	viper.AddConfigPath("configs/")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed to read properties config: %w", err))
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}
