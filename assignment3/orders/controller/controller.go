package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Order struct {
	OrderId string `json:"order_id"`
	ProductId string `json:"productId"`
	OrderName string `json:"orderName"`
	Price string `json:"price"`
	UserId string `json:"userId"`
}

func ConnectToDatabase() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/ecommerce")
	
	if err != nil {
		fmt.Println("Error during connection establishment!", err)
		return nil;
	}

	fmt.Println("DB connected successfully", db)

	return db
}


func GetAllOrders(c *gin.Context) {

	userId := c.Param("userId")

	db := ConnectToDatabase()

	rows, err := db.Query("SELECT * FROM orders WHERE userId=?", userId)

	var orders []Order
	

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":"Some error occured",
		})
	} else {
		for rows.Next() {
			var order Order
			err := rows.Scan(&order.OrderId, &order.ProductId, &order.UserId, &order.Price, &order.OrderName)

			if err != nil {
				fmt.Println("error while getting all orders : ", err)
			}

			orders = append(orders, order)
		} 

		c.JSON(http.StatusOK, gin.H{
			"data": orders,
		})
	}
}

func GetOrderById(c *gin.Context) {
	orderId := c.Param("orderId")

	db := ConnectToDatabase()

	var order Order

	err := db.QueryRow("SELECT * FROM orders WHERE orderId=?", orderId).Scan(&order.OrderId, &order.ProductId, &order.UserId, &order.Price, &order.OrderName)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":"Order Not found",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": order,
		})
	}
}

func PlaceOrder(c *gin.Context) {
	body := c.Request.Body
	userId := c.Param("userId")

	orderData, _ := ioutil.ReadAll(body)

	db := ConnectToDatabase()
	
	var newOrder Order

	json.Unmarshal(orderData, &newOrder)
	newUUID, _ := exec.Command("uuidgen").Output()

	newOrder.OrderId = strings.TrimSuffix(string(newUUID), "\n")
	newOrder.UserId = userId

	_, err := db.Exec("INSERT INTO orders VALUES('" + newOrder.OrderId + "', '" + newOrder.ProductId + "', '" + newOrder.UserId + "', '" + newOrder.Price + "', '" + newOrder.OrderName + "')")

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Some error occured",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data":"Order done!!",
		})
	}

}

func DeleteAllOrders(c *gin.Context) {
	userId := c.Param("userId")

	db := ConnectToDatabase()

	_, err := db.Query("DELETE FROM orders WHERE userId=?", userId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Some error occured",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data":"Orders deleted successfully",
		})
	}

}

func DeleteOrder(c *gin.Context) {
	userId := c.Param("userId")
	orderId := c.Param("orderId")

	db := ConnectToDatabase()

	_, err := db.Query("DELETE FROM orders WHERE userId=? AND orderId=?", userId, orderId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Some error occured",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data":"Orders deleted successfully",
		})
	}

}

func UpdateOrder(c *gin.Context) {
	orderId := c.Param("orderId")
	userId := c.Param("userId")
	body := c.Request.Body

	bodyData, _ := ioutil.ReadAll(body)

	var updatedOrderData Order

	json.Unmarshal(bodyData, &updatedOrderData)
	updatedOrderData.OrderId = orderId

	db := ConnectToDatabase()

	fmt.Println(updatedOrderData)

	// var order Order

	res, err := db.Exec("UPDATE orders SET orderName=? WHERE orderId=? AND userId=?", updatedOrderData.OrderName, updatedOrderData.OrderId, userId)



	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Some error occured while updating order.",
		})
	} else {

		ans,_ := res.LastInsertId()
		fmt.Println(ans)

		// rows.Scan(&order.OrderId, &order.OrderName, &order.Price, &order.ProductId, &order.UserId)

		// fmt.Println(order)
		c.JSON(http.StatusOK, gin.H{
			"data": "Order updated Successfully!!",
		})
	}

}

