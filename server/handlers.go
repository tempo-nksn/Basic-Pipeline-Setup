package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func hello(c *gin.Context) {
	c.String(200, "Hello User, your taxi is booked")
}

func dummyPost(c *gin.Context) {
	var str struct {
		Name string `json:"name"`
	}
	c.BindJSON(&str)
	fmt.Printf("I am posting to ios")
	c.JSON(200, str)
}
