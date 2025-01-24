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
	calculatorapplication "github.com/cfioretti/recipe-manager/internal/ingredients-balancer/application"
	recipeapplication "github.com/cfioretti/recipe-manager/internal/recipe-manager/application"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/infrastructure/http"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/infrastructure/mysql"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/infrastructure/mysql/migrations"
)

func main() {
	setConfigs()
	config := configs.NewDBConfig()
	db, err := newDBConnection(config)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %v", err))
	}

	router := makeRouter(db)
	port := viper.GetInt("server.port")
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		panic(fmt.Errorf("failed to start server: %v", err))
	}
}

func makeRouter(dB *sql.DB) *gin.Engine {
	recipeController := http.NewRecipeController(
		recipeapplication.NewRecipeService(
			mysql.NewMySqlRecipeRepository(dB),
			calculatorapplication.NewIngredientsCalculatorService(),
			calculatorapplication.NewIngredientsBalancerService(),
		),
	)

	router := gin.Default()
	router.Use(corsMiddleware())
	router.POST("/recipes/:uuid/aggregate", recipeController.RetrieveRecipeAggregate)
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

func newDBConnection(config *configs.DBConfig) (*sql.DB, error) {
	if err := migrations.RunMigrations(config.DSN(), "migrations", viper.GetString("database.dbName")); err != nil {
		return nil, fmt.Errorf("error executing db migrations: %w", err)
	}

	db, err := sql.Open("mysql", config.DSN())
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	return db, nil
}

func setConfigs() {
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
