package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xeipuuv/gojsonschema"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World!",
		})
	})

	// // Not valid request
	// curl --location --request POST 'http://localhost:8080/user' \
	// --header 'Content-Type: application/json' \
	// --header 'Cookie: messages="[[\"__json_message\"\0540\05425\054\"The lead outlet \\\"leadtesting\\\": synchronizing to Odoo. Refresh to see update.\"]]:1oSxRv:ll0lbcUXIWtvEqqaTOW0k7jVJELid0DTKwj0AEWIg78"' \
	// --data-raw '{
	// 	"name":"",
	// 	"email":"",
	// 	"phone_number":"0811223344",
	// 	"address":"Jalan Merdeka Barat"
	// }'

	// // Valid request
	// curl --location --request POST 'http://localhost:8080/user' \
	// --header 'Content-Type: application/json' \
	// --header 'Cookie: messages="[[\"__json_message\"\0540\05425\054\"The lead outlet \\\"leadtesting\\\": synchronizing to Odoo. Refresh to see update.\"]]:1oSxRv:ll0lbcUXIWtvEqqaTOW0k7jVJELid0DTKwj0AEWIg78"' \
	// --data-raw '{
	// 	"name":"ado",
	// 	"email":"adopabianko@gmail.com",
	// 	"phone_number":"0811223344",
	// 	"address":"Jalan Merdeka Barat"
	// }'
	r.POST("/user", CreateUser)

	r.Run()
}

type User struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
}

func CreateUser(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err,
		})
		return
	}

	userJson, _ := json.Marshal(user)
	// Schema validation
	userSchema := gojsonschema.NewReferenceLoader("file://./validations/user_schema.json")
	// Request data from client
	userDocument := gojsonschema.NewBytesLoader(userJson)

	// Validate request data
	validate, err := gojsonschema.Validate(userSchema, userDocument)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err,
		})
		return
	}

	// if not valid
	if !validate.Valid() {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": fmt.Sprintf("%s", validate.Errors()),
		})

		return
	}

	// success reponse
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "user successfully created",
	})
}
