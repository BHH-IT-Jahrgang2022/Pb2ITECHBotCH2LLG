package api

import (
	"encoding/json"
	fmt "fmt"
	"io"
	"net/http"
	"os"

	websocket "github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Response string `json:"message"`
}

// Helper function handling the chat logic for the chat endpoint
// This function is not exported and can only be used within the api package
// Call as Go routine
func chatFunc(query string, analyzer_route string) string {
	// Chat logic
	result := "Something went wrong"
	// Send the Text to the analyzer as query parameters and get the JSON response
	response, err := http.Get(analyzer_route + "?query=" + query)
	if err != nil {
		fmt.Println(err)
	} else {
		// Decode the JSON response
		var res Response
		json.NewDecoder(response.Body).Decode(&res)
		result = res.Response
	}
	return result
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func readFromSocket(c *websocket.Conn) string {
	_, message, err := c.ReadMessage()
	if err != nil {
		if err == io.EOF {
			fmt.Println("Connection closed")
		} else {
			fmt.Println(err)
		}
	}
	return string(message)
}

func writeToSocket(c *websocket.Conn, message string) {
	err := c.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		fmt.Println(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		// Get environment variable for the analyzer route
		// If not set, use default value
		analyzer_route := os.Getenv("ANALYZER_ROUTE")
		if analyzer_route == "" {
			analyzer_route = "http://localhost:8081/getanswer"
		} else {
			for {
				message := readFromSocket(conn)
				fmt.Println("Message received: " + message)
				answer := chatFunc(message, analyzer_route)
				writeToSocket(conn, answer)
			}
		}
	}
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
	//Create Websocket connection with Client
	r.GET("/chat", func(c *gin.Context) {
		go handler(c.Writer, c.Request)
	})
	// Endpoint only for testing purposes
	r.GET("chat/socketless", func(c *gin.Context) {
		answer := chatFunc("Hello", "http://localhost:8081/getanswer")
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
