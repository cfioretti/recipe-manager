package main

import (
	"fmt"
	"os"
	"strings"

	"recipe-manager/internal/recipe-manager/application"
	"recipe-manager/internal/recipe-manager/infrastructure/http"
	"recipe-manager/internal/recipe-manager/infrastructure/mysql"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	setConfig()

	router := makeRouter()
	port := viper.GetInt("server.port")
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		panic("failed to start server")
	}
}

func makeRouter() *gin.Engine {
	recipeController := http.NewRecipeController(
		application.NewRecipeService(mysql.NewMysqlRecipeRepository()),
	)

	router := gin.Default()
	router.POST("/recipes/:uuid/aggregate", recipeController.RetrieveRecipeAggregate)
	return router
}

func setConfig() {
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
