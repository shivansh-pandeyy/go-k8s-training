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

type User struct {
	Username string `json:"username"`
	UserId string `json:"userId"`
}

type Product struct {
	ProductId string `json:"productID"`
	Name string `json:"name"`
	Description string `json:"description"`
	Price string `json:"price"`
}

func Db_connectivity() *sql.DB {
	// add data into database
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/ecommerce")
	
	if err != nil {
		fmt.Println("Error during connection establishment!", err)
		return nil;
	}

	fmt.Println("DB connected successfully", db)

	return db
}

func GetAllUsers(c *gin.Context) {
	db := Db_connectivity()

	if db == nil {
		fmt.Println("DB is nil")
	}

	rows, e := db.Query("SELECT * FROM users")

	if e != nil {	
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Some error occured",
		})

		return
	}

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.UserId, &user.Username)
		if err != nil {
			fmt.Println("error while getting all users : ", err)
		}

		users = append(users, user)
	}

		c.JSON(200, gin.H{
			"data": users,
		})
}

func GetUser(c *gin.Context) {
	id := c.Param("userId")

	db := Db_connectivity()

	var user User

	err := db.QueryRow("SELECT * from users where userID=?", id).Scan(&user.UserId, &user.Username)

	// userData, _ = json.Marshal(user)

	fmt.Println(user)
	
	if err != nil {
		fmt.Println("Error from DB:", err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found.",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": user,
		})
	}

	defer db.Close()
}

func DeleteUser(c *gin.Context) {

	id := c.Param("userId")

	db := Db_connectivity()

	_, err := db.Exec("DELETE FROM users WHERE userID=?", id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Some error occured",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": "Deletion successful",
		})
	}

	db.Close()

}

func CreateUser(c *gin.Context) {

	body := c.Request.Body
	userData, _ := ioutil.ReadAll(body)

	db := Db_connectivity()

	var newUser User

	json.Unmarshal(userData, &newUser)

	
	newUUID, _ := exec.Command("uuidgen").Output()
	
	newUser.UserId = strings.TrimSuffix(string(newUUID), "\n")

	fmt.Println("new user data", newUser)
	
	_, err := db.Exec("INSERT INTO users VALUES('" + newUser.UserId +"', '" + newUser.Username + "')")

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error" : "Some error occured while creating a new user",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": "data added successfully",
		})
	}

	defer db.Close()

}

func ListAllProducts(c *gin.Context) {
	db := Db_connectivity()

	rows, err := db.Query("SELECT * FROM product")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		fmt.Println(err)
	}

	var productList []Product

	for rows.Next() {
		var product Product

		err := rows.Scan(&product.ProductId,  &product.Name, &product.Description, &product.Price)

		if err != nil {
			fmt.Println("Error in getting product list")
		}

		fmt.Println(product)

		productList = append(productList, product)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":productList,
	})

}

