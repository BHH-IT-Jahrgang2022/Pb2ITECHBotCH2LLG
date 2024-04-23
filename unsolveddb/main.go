package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"unsolveddb/unsolveddb"

	"github.com/gin-gonic/gin"
)

type RequestBody struct {
	Problem string   `json:"problem"`
	Tags    []string `json:"tags"`
}

func main() {

	tickets := unsolveddb.FetchTicket()

	fmt.Println("Fetched data from MongoDB")
	fmt.Println((*tickets)[0].Problem)

	router := gin.Default()

	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": unsolveddb.Test(),
		})
	})

	router.GET("/data", func(c *gin.Context) {
		tickets := unsolveddb.FetchTicket()
		c.JSON(200, tickets)
	})

	router.POST("/insert", func(c *gin.Context) {
		var requestBody RequestBody

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		err2 := json.Unmarshal(body, &requestBody)
		if err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		unsolveddb.InsertTicket(unsolveddb.Ticket{Tags: requestBody.Tags, Problem: requestBody.Problem})
		c.JSON(200, requestBody)
	})
	router.Run("0.0.0.0:8080")
}
