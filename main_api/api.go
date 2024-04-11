package main

import (
	fmt "fmt"
	"net/http"

	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

func StartApi() {
	sessions := make(map[string]bool)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("connect", func(c *gin.Context) {
		session_token := uuid.NewV4().String()
		sessions[session_token] = true
		c.JSON(http.StatusOK, gin.H{
			"session_token": session_token,
		})
	})
	r.GET("disconnect", func(c *gin.Context) {
		session_token := c.Query("session_token")
		fmt.Println(sessions)
		fmt.Println(session_token)
		if _, ok := sessions[session_token]; ok {
			if sessions[session_token] {
				sessions[session_token] = false
				c.JSON(http.StatusOK, gin.H{
					"message": "disconnected",
				})
			} else {
				c.JSON(http.StatusNotFound, gin.H{
					"message": "session already disconnected",
				})
			}
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "session not found",
			})
		}
	})
	r.GET("chat", func(c *gin.Context) {
		session_token := c.Query("session_token")
		if _, ok := sessions[session_token]; ok {
			c.JSON(http.StatusOK, gin.H{
				"message": "chat",
			})
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "session not found",
			})
		}
	})
	r.GET("chat/sessionless", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "chat",
		})
	})
	r.GET("")
	r.Run(":8080")
}

func main() {
	StartApi()
}
