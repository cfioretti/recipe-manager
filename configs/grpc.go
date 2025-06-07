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
