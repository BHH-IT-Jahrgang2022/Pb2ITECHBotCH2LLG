package main

import (
	"db/db"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	data := db.FetchData()

	fmt.Println("Fetched data from MongoDB")
	fmt.Println((*data)[0].Keywords)

	router := gin.Default()
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": db.Test(),
		})
	})

	router.GET("/data", func(c *gin.Context) {
		c.JSON(200, data)
	})

	router.Run("127.0.0.1:8080")
}
