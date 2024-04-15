package api

import (
	"bytes"
	"encoding/json"
	fmt "fmt"
	"net/http"

	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

type Query struct {
	SessionToken string `json:"session_token"`
	Query string `json:"query"`
}

type Query_Text struct {
	Query string `json:"query"`
}

type Matcher_Response struct {
	Response string `json:"response"`
}

type Response struct {
	SessionToken string `json:"session_token"`
	Query string `json:"query"`
	Response string `json:"message"`
}

// Helper function handling the chat logic for the chat endpoint
// This function is not exported and can only be used within the api package
// Call as Go routine
func chatFunc(query string, tokenizer_route string, matcher_route string) string {
	// Chat logic
	result := "Something went wrong"
	// Build the query as json
	query_json, err := json.Marshal(Query_Text{Query: query})
	if err != nil {
		fmt.Println(err)
	} else {
		// Call the tokenizer
		tokenizer_response, err := http.Post(tokenizer_route, "application/json", bytes.NewBuffer(query_json))
		if err != nil {
			fmt.Println(err)
		} else {
			// Call the matcher
			matcher_response, err := http.Post(matcher_route, "application/json", tokenizer_response.Body)
			if err != nil {
				fmt.Println(err)
			} else {
				// Get the result
				decoded_matcher_response := Matcher_Response{}
				err := json.NewDecoder(matcher_response.Body).Decode(&decoded_matcher_response)
				if err != nil {
					fmt.Println(err)
				} else {
					result = decoded_matcher_response.Response
				}
			}
		}
	}
	return result
}

// Starts the Api
func StartApi() {
	sessions := make(map[string]bool)
	r := gin.Default()
	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	// Endpoint for connecting a session
	r.GET("connect", func(c *gin.Context) {
		session_token := uuid.NewV4().String()
		sessions[session_token] = true
		c.JSON(http.StatusOK, gin.H{
			"session_token": session_token,
		})
	})
	// Endpoint for disconnecting a session
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
	// Main Endpoint for chatting, calls all other services
	r.GET("chat", func(c *gin.Context) {
		var query Query
		err := c.BindJSON(&query)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid json",
			})
			return
		} else {
			if _, ok := sessions[query.SessionToken]; ok {
				if sessions[query.SessionToken] {
					answer := chatFunc(query.Query, "http://localhost:8081/tokenize", "http://localhost:8082/match")
					c.JSON(http.StatusOK, gin.H{
						"session_token": query.SessionToken,
						"query": query.Query,
						"message": answer,
					})

	// Endpoint only for testing purposes
	r.GET("chat/sessionless", func(c *gin.Context) {
		answer := chatFunc("Hello", "http://localhost:8081/tokenize", "http://localhost:8082/match")
		c.JSON(http.StatusOK, gin.H{
			"message": answer,
		})
	})
	r.GET("")
	r.Run(":8080")
	// Endpoint for getting all sessions
	r.GET("sessions", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"sessions": sessions,
		})
	})
}
