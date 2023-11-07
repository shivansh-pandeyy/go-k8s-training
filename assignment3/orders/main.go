package main

import (
	"ecommerce/orders/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/users/:userId/orders", controller.GetAllOrders)
	router.GET("/users/:userId/orders/:orderId", controller.GetOrderById)
	router.POST("/users/:userId/orders", controller.PlaceOrder)
	router.DELETE("/users/:userId/orders", controller.DeleteAllOrders)
	router.DELETE("/users/:userId/orders/:orderId", controller.DeleteOrder)
	router.PUT("/users/:userId/orders/:orderId", controller.UpdateOrder)

	router.Run(":9090")
}
