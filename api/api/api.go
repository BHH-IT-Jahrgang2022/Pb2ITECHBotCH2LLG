package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	websocket "github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Response string `json:"message"`
}

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
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
		log_entry := LogEntry{
			Timestamp: time.Now().String(),
			Level:     "Error",
			Message:   err.Error(),
		}
		// Send the log entry to the logging API
		log(log_entry)
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
			log_entry := LogEntry{
				Timestamp: time.Now().String(),
				Level:     "Error",
				Message:   err.Error(),
			}
			// Send the log entry to the logging API
			log(log_entry)
		} else {
			log_entry := LogEntry{
				Timestamp: time.Now().String(),
				Level:     "Error",
				Message:   err.Error(),
			}
			// Send the log entry to the logging API
			log(log_entry)
		}
	}
	return string(message)
}

func writeToSocket(c *websocket.Conn, message string) {
	err := c.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log_entry := LogEntry{
			Timestamp: time.Now().String(),
			Level:     "Error",
			Message:   err.Error(),
		}
		// Send the log entry to the logging API
		log(log_entry)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log_entry := LogEntry{
			Timestamp: time.Now().String(),
			Level:     "Error",
			Message:   err.Error(),
		}
		// Send the log entry to the logging API
		log(log_entry)
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
				answer := chatFunc(message, analyzer_route)
				writeToSocket(conn, answer)
			}
		}
	}
}

func log(entry LogEntry) {
	// Send the log entry to the logging API
	if os.Getenv("LOGGING_ENABLED") == "true" {
		logging_API_route := os.Getenv("LOGGING_API_ROUTE")
		jsonEntry, _ := json.Marshal(entry)
		_, err := http.Post(logging_API_route, "application/json", bytes.NewBuffer(jsonEntry))
		if err != nil {
			log_entry := LogEntry{
				Timestamp: time.Now().String(),
				Level:     "Error",
				Message:   err.Error(),
			}
			// Send the log entry to the logging API
			log(log_entry)
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
