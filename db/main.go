package main

import (
	"db/db"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestBody struct {
	Collection string  `json:"collection"`
	Data       db.Data `json:"data"`
}

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
		data := db.FetchData()
		c.JSON(200, data)
	})

	router.POST("/insert", func(c *gin.Context) {
		var requestBody RequestBody

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("got JSON")
		db.InsertData(requestBody.Data, requestBody.Collection)
		c.JSON(200, requestBody)
	})

	router.Run("127.0.0.1:8080")
}
