package configs

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type GRPCConfig struct {
	Address string
	Timeout time.Duration
}

func LoadCalculatorGRPCConfig() GRPCConfig {
	calculatorAddr := os.Getenv("CALCULATOR_ADDR")
	if calculatorAddr == "" {
		host := viper.GetString("grpc.host")
		port := viper.GetString("grpc.calculator.port")
		calculatorAddr = host + ":" + port
	}
	timeout := 5 * time.Second
	log.Println("calculatorAddr: ", calculatorAddr)

	return GRPCConfig{
		Address: calculatorAddr,
		Timeout: timeout,
	}
}

func LoadBalancerGRPCConfig() GRPCConfig {
	balancerAddr := os.Getenv("INGREDIENTS_BALANCER_ADDR")
	if balancerAddr == "" {
		host := viper.GetString("grpc.host")
		port := viper.GetString("grpc.balancer.port")
		balancerAddr = host + ":" + port
	}
	timeout := 5 * time.Second
	log.Println("balancerAddr: ", balancerAddr)

	return GRPCConfig{
		Address: balancerAddr,
		Timeout: timeout,
	}
}
