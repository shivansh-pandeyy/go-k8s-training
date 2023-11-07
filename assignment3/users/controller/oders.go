package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Order struct {
	OrderId string `json:"order_id"`
	ProductName string `json:"productName"`
	Description string `json:"description"`
	Price int `json:"price"`
}

var (
	orderServiceHost     = "localhost"
	orderServicePort     = "9090"
	orderServiceEndpoint = "http://" + orderServiceHost + ":" + orderServicePort
)

func sendRequestToOrderService(context *gin.Context, httpMethod string, url string) {
	fmt.Println("In the proxy, url = ", url)

	client := http.Client{}

	req, err := http.NewRequestWithContext(context, httpMethod, url, context.Request.Body)
	if err != nil {
		fmt.Println("Error occurred while new Request!", err)
	}

	var newOrder Order

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error during bindjson", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error during reading body", err)
	}

	// if httpMethod == "PUT" {
	// 	context.JSON(http.StatusOK, "Successfully Updated!!")
	// 	return
	// }

	fmt.Println(string(bodyBytes))

	// err = json.Unmarshal(bodyBytes, &newOrder)

	// if err != nil {
	// 	fmt.Println("Error during json unmarshalling", err)
	// }

	fmt.Println("resp code:", resp.StatusCode, ", neworder: ", newOrder)
	context.JSON(http.StatusOK, json.RawMessage(string(bodyBytes)))
}

func checkValidUser(userId string) bool {
	db := Db_connectivity()

	var user User

	err := db.QueryRow("SELECT * FROM users WHERE userId=?", userId).Scan(&user.UserId, &user.Username)

	if err != nil {
		fmt.Println("error while checking for valid user", err)
		return false
	}

	return true
}

func GetAllOrders(c *gin.Context) {

	userId := c.Param("userId")

	validUser := checkValidUser(userId)

	if !validUser {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
	} else {
		url := orderServiceEndpoint + "/users/" + userId + "/orders"
		sendRequestToOrderService(c, "GET", url)
	}
}

func GetOrderById(c *gin.Context) {

	userId := c.Param("userId")
	orderId := c.Param("orderId")

	validUser := checkValidUser(userId)
	
	if !validUser {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
	} else {
		url := orderServiceEndpoint + "/users/" + userId + "/orders/" + orderId
		sendRequestToOrderService(c, "GET", url)
	}

}

func DeleteAllOrders(c *gin.Context) {
	userId := c.Param("userId")

	validUser := checkValidUser(userId)
	
	if !validUser {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
	} else {
		url := orderServiceEndpoint + "/users/" + userId + "/orders"
		sendRequestToOrderService(c, "DELETE", url)
	}
}

func DeleteOrder(c *gin.Context) {
	userId := c.Param("userId")
	orderId := c.Param("orderId")

	validUser := checkValidUser(userId)
	
	if !validUser {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
	} else {
		url := orderServiceEndpoint + "/users/" + userId + "/orders/" + orderId
		sendRequestToOrderService(c, "DELETE", url)
	}
}


func UpdateOrder(c *gin.Context) {
	userId := c.Param("userId")
	orderId := c.Param("orderId")

	validUser := checkValidUser(userId)
	
	if !validUser {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
	} else {
		url := orderServiceEndpoint + "/users/" + userId + "/orders/" + orderId
		sendRequestToOrderService(c, "PUT", url)
	}
}

func PlaceOrder(c *gin.Context) {
	userId := c.Param("userId")

	validUser := checkValidUser(userId)
	
	if !validUser {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
	} else {
		url := orderServiceEndpoint + "/users/" + userId + "/orders"
		sendRequestToOrderService(c, "POST", url)
	}
}
