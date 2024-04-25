package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	websocket "github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Response string `json:"response"`
}

type LogEntry struct {
	Timestamp int64  `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Service   string `json:"service"`
}

func logger(entry LogEntry) {
	// Send the log entry to the logging API
	if os.Getenv("LOGGING_ENABLED") == "true" {
		logging_API_route := os.Getenv("LOGGING_API_ROUTE")
		jsonEntry, _ := json.Marshal(entry)
		http.Post(logging_API_route+"/log", "application/json", bytes.NewBuffer(jsonEntry))
	}
}

// Helper function handling the chat logic for the chat endpoint
// This function is not exported and can only be used within the api package
// Call as Go routine
func chatFunc(query string, analyzer_route string) string {
	// Chat logic
	result := "Something went wrong"
	// Send the Text to the analyzer as query parameters and get the JSON response
	fixed_query := strings.ReplaceAll(query, " ", "%20")
	fixed_query = strings.ReplaceAll(fixed_query, "?", "%3F")
	fixed_query = strings.ReplaceAll(fixed_query, "!", "%21")
	fixed_query = strings.ReplaceAll(fixed_query, "#", "%23")
	fixed_query = strings.ReplaceAll(fixed_query, "&", "%26")
	fixed_query = strings.ReplaceAll(fixed_query, "=", "%3D")
	fixed_query = strings.ReplaceAll(fixed_query, "+", "%2B")
	fixed_query = strings.ReplaceAll(fixed_query, "/", "%2F")
	fixed_query = strings.ReplaceAll(fixed_query, ":", "%3A")
	fixed_query = strings.ReplaceAll(fixed_query, ";", "%3B")
	fixed_query = strings.ReplaceAll(fixed_query, "@", "%40")
	fixed_query = strings.ReplaceAll(fixed_query, "[", "%5B")
	fixed_query = strings.ReplaceAll(fixed_query, "]", "%5D")
	fixed_query = strings.ReplaceAll(fixed_query, "{", "%7B")
	fixed_query = strings.ReplaceAll(fixed_query, "}", "%7D")
	fixed_query = strings.ReplaceAll(fixed_query, "\"", "%22")
	fixed_query = strings.ReplaceAll(fixed_query, "'", "%27")
	fixed_query = strings.ReplaceAll(fixed_query, "<", "%3C")
	fixed_query = strings.ReplaceAll(fixed_query, ">", "%3E")
	fixed_query = strings.ReplaceAll(fixed_query, "\r", "")
	fixed_query = strings.ReplaceAll(fixed_query, "\n", "")
	response, err := http.Get(analyzer_route + "?query=" + fixed_query)
	if err != nil {
		log_entry := LogEntry{
			Timestamp: time.Now().Unix(),
			Level:     "Error",
			Message:   err.Error(),
			Service:   "api",
		}
		// Send the log entry to the logging API
		logger(log_entry)
		if os.Getenv("DEBUG") == "true" {
			result += " Error: " + err.Error()
		}
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
				Timestamp: time.Now().Unix(),
				Level:     "Error",
				Message:   err.Error(),
				Service:   "api",
			}
			// Send the log entry to the logging API
			logger(log_entry)
		} else {
			log_entry := LogEntry{
				Timestamp: time.Now().Unix(),
				Level:     "Error",
				Message:   err.Error(),
				Service:   "api",
			}
			// Send the log entry to the logging API
			logger(log_entry)
		}
	}
	return string(message)
}

func writeToSocket(c *websocket.Conn, message string) {
	err := c.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log_entry := LogEntry{
			Timestamp: time.Now().Unix(),
			Level:     "Error",
			Message:   err.Error(),
			Service:   "api",
		}
		// Send the log entry to the logging API
		logger(log_entry)
		fmt.Println("I BROKE DOWN")
	}
	fmt.Println("I AM DONE")
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log_entry := LogEntry{
			Timestamp: time.Now().Unix(),
			Level:     "Error",
			Message:   err.Error(),
			Service:   "api",
		}
		// Send the log entry to the logging API
		logger(log_entry)
		return
	}
	defer conn.Close()
	// Get environment variable for the analyzer route
	// If not set, use default value
	analyzer_route := os.Getenv("ANALYZER_ROUTE")
	for {
		message := readFromSocket(conn)
		answer := chatFunc(message, analyzer_route)
		fmt.Println("[DEBUG] ", answer)
		writeToSocket(conn, answer)
	}
}

// Starts the Api
func StartApi() {
	sessions := make(map[string]bool)

	r := gin.Default()
	r.GET("")
	//Create Websocket connection with Client
	r.GET("/chat", func(c *gin.Context) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		handler(c.Writer, c.Request)
	})
	// ENDPOINTS FOR DEBUGGING PURPOSES //
	// -------------------------------- //
	if os.Getenv("DEBUG") == "true" {
		// Ping test
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
		// Endpoint only for testing purposes
		r.GET("chat/socketless", func(c *gin.Context) {
			answer := chatFunc("Hello", os.Getenv("ANALYZER_ROUTE"))
			c.JSON(http.StatusOK, gin.H{
				"message": answer,
			})
		})
		// Endpoint for getting all sessions
		r.GET("sessions", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"sessions": sessions,
			})
		})
		r.GET("log/test", func(c *gin.Context) {
			fmt.Println("Test log entry")
			var testlog LogEntry
			testlog.Timestamp = time.Now().Unix()
			testlog.Level = "Info"
			testlog.Message = "Test log entry"
			testlog.Service = "api"
			c.JSON(http.StatusOK, gin.H{
				"message": "Log entry sent",
			})
		})
	}
	r.Run(":8080")
}
