package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"recipe-manager/configs"
	"recipe-manager/internal/recipe-manager/application"
	"recipe-manager/internal/recipe-manager/infrastructure/http"
	"recipe-manager/internal/recipe-manager/infrastructure/mysql"
	"recipe-manager/internal/recipe-manager/infrastructure/mysql/migrations"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func main() {
	setConfigs()
	config := configs.NewDBConfig()
	db, err := newDBConnection(config)
	if err != nil {
		panic(err)
	}

	router := makeRouter(db)
	port := viper.GetInt("server.port")
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		panic("failed to start server")
	}
}

func makeRouter(dB *sql.DB) *gin.Engine {
	recipeController := http.NewRecipeController(
		application.NewRecipeService(mysql.NewMysqlRecipeRepository(dB)),
	)

	router := gin.Default()
	router.POST("/recipes/:uuid/aggregate", recipeController.RetrieveRecipeAggregate)
	return router
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
