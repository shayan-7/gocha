package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shayan-7/gocha/internal/services"
)

type OrderHandler struct {
	orderService *services.OrderService
}

func NewOrderHandler(o *services.OrderService) *OrderHandler {
	return &OrderHandler{orderService: o}
}

func (oh *OrderHandler) PostOrderHandler(c *gin.Context) {
	rawData, err := c.GetRawData()
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	resp := oh.orderService.Publish(string(rawData))
	log.Println("Cache response:", resp)
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}
