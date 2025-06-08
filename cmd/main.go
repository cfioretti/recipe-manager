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
	capplication "github.com/cfioretti/recipe-manager/internal/ingredients-balancer/application"
	rapplication "github.com/cfioretti/recipe-manager/internal/recipe-manager/application"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/infrastructure/grpc/client"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/infrastructure/mysql"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/interfaces/api/http"
)

func main() {
	setEnvConfigs()
	db := loadDBConfig()
	calculatorService, err := initializeCalculatorService()
	if err != nil {
		calculatorService = capplication.NewCalculatorService()
	}
	balancerService, err := initializeBalancerService()
	if err != nil {
		balancerService = capplication.NewIngredientsBalancerService()
	}
	router := makeRouter(db, calculatorService, balancerService)
	port := viper.GetInt("server.port")
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		panic(fmt.Errorf("failed to start server: %v", err))
	}
}

func makeRouter(dB *sql.DB, calculatorService rapplication.CalculatorService, balanceService rapplication.BalancerService) *gin.Engine {
	recipeHandler := http.NewRecipeHandler(
		rapplication.NewRecipeService(
			mysql.NewMySqlRecipeRepository(dB),
			calculatorService,
			balanceService,
		),
	)

	router := gin.Default()
	router.Use(corsMiddleware())
	router.POST("/recipes/:uuid/aggregate", recipeHandler.RetrieveRecipeAggregate)
	return router
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

func initializeCalculatorService() (rapplication.CalculatorService, error) {
	calculatorClient, err := loadCalculatorGrpcClient()
	if err != nil {
		return nil, err
	}
	return rapplication.NewRemoteDoughCalculatorService(calculatorClient), nil
}

func loadCalculatorGrpcClient() (*client.CalculatorClient, error) {
	calculatorGRPCConfig := configs.LoadCalculatorGRPCConfig()

	calculatorClient, err := client.NewDoughCalculatorClient(
		calculatorGRPCConfig.Address,
		calculatorGRPCConfig.Timeout,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Calculator grpc client: %w", err)
	}
	return calculatorClient, nil
}

func initializeBalancerService() (rapplication.BalancerService, error) {
	balancerClient, err := loadBalancerGrpcClient()
	if err != nil {
		return nil, err
	}
	return rapplication.NewRemoteIngredientsBalancerService(balancerClient), nil
}

func loadBalancerGrpcClient() (*client.IngredientsBalancerClient, error) {
	balancerGRPCConfig := configs.LoadBalancerGRPCConfig()

	balancerClient, err := client.NewIngredientsBalancerClient(
		balancerGRPCConfig.Address,
		balancerGRPCConfig.Timeout,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Balancer grpc client: %w", err)
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
