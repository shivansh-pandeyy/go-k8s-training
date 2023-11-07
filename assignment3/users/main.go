package main

import (
	"ecommerce/users/controller"

	"github.com/gin-gonic/gin"
)



func main() {

	router := gin.Default();

	router.GET("/users", controller.GetAllUsers)
	router.GET("/users/:userId", controller.GetUser)
	router.DELETE("/users/:userId", controller.DeleteUser)
	router.POST("/users", controller.CreateUser)

	router.GET("/products", controller.ListAllProducts)

	router.GET("/users/:userId/orders", controller.GetAllOrders)
	router.GET("/users/:userId/orders/:orderId", controller.GetOrderById)
	router.POST("/users/:userId/orders", controller.PlaceOrder)
	router.DELETE("/users/:userId/orders", controller.DeleteAllOrders)
	router.DELETE("/users/:userId/orders/:orderId", controller.DeleteOrder)
	router.PUT("/users/:userId/orders/:orderId", controller.UpdateOrder)

	router.Run()

}
