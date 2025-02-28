package configs

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type GRPCConfig struct {
	IngredientsBalancerAddress string
	Timeout                    time.Duration
}

func LoadGRPCConfig() GRPCConfig {
	ingredientsBalancerAddr := os.Getenv("INGREDIENTS_BALANCER_ADDR")
	if ingredientsBalancerAddr == "" {
		host := viper.GetString("grpc.host")
		port := viper.GetString("grpc.port")
		ingredientsBalancerAddr = host + ":" + port
	}
	timeout := 5 * time.Second
	log.Println("ingredientsBalancerAddr: ", ingredientsBalancerAddr)

	return GRPCConfig{
		IngredientsBalancerAddress: ingredientsBalancerAddr,
		Timeout:                    timeout,
	}
}
