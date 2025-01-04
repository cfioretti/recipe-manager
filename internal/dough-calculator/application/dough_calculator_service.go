package application

import (
	"github.com/gin-gonic/gin"
)

type DoughCalculatorService struct{}

func NewDoughCalculatorService() *DoughCalculatorService {
	return &DoughCalculatorService{}
}

func (dc DoughCalculatorService) TotalPansWeight(params gin.Params) (int, error) {
	return 0, nil
}
