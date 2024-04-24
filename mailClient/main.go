package main

import (
	"mailClient/mailClient"

	"github.com/gin-gonic/gin"
)

type Unresolved struct {
	Query    string `json:"query"`
	Response string `json:"response"`
}

func start_API() {
	router := gin.Default()
	router.POST("/ticket/all", func(c *gin.Context) {
		mailClient.FetchAndPrintTickets()
		c.JSON(200, gin.H{
			"message": "Email sent",
		})
	})
	router.POST("/ticket/matchfailed", func(c *gin.Context) {
		var unresolved Unresolved
		c.BindJSON(&unresolved)
		problem_text := "Got a match failed for: " + unresolved.Query + " result: " + unresolved.Response
		mailClient.SendEmail(&mailClient.Ticket{Tags: []string{"unresolved"}, Problem: problem_text})
		c.JSON(200, gin.H{
			"message": "Email sent",
		})
	})
	router.Run("0.0.0.0:8080")
}

func main() {
	start_API()
	//mailClient.FetchAndEmailTicket()
}
