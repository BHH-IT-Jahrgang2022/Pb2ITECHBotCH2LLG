package main

import (
	"matcher/matcher"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	matches := matcher.LoadTable()

	/*
		fmt.Println("Starting Tests now...")

		test1 := "testy hat heute ??? leider blup"
		test2 := "bla mal sehen, ohne Treffer"
		test3 := "blup"
		fmt.Println(test1)
		fmt.Println(matcher.Match(test1, matches))
		fmt.Println(test2)
		fmt.Println(matcher.Match(test2, matches))
		fmt.Println(test3)
		fmt.Println(matcher.Match(test3, matches))
	*/

	router := gin.Default()
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": matcher.Test(),
		})
	})
	router.POST("/match", func(c *gin.Context) {
		var json struct {
			Input []string `json:"tokens"`
		}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		req := strings.Join(json.Input, "")
		c.JSON(200, gin.H{
			"query":    json.Input,
			"response": matcher.Match(req, matches),
		})
	})
	router.GET("/reload", func(c *gin.Context) {
		matches = matcher.LoadTable()
		c.JSON(200, gin.H{
			"message": "reloaded",
		})
	})

	router.Run("127.0.0.1:8080")

}
