package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/gorilla/websocket"
)

type Input struct {
	Text      string `json:"text"`
	SessionID string `json:"token"`
}


type LogData struct {
	Error     string `json:"error"`
	Timestamp int64  `json:"timestamp"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func logError(err error) {
	// Get the current Unix timestamp
	timestamp := time.Now().Unix()

	// Create the log data
	logData := LogData{Error: err.Error(), Timestamp: timestamp}

	// Convert the log data to JSON
	jsonLogData, _ := json.Marshal(logData)

	// Send a POST request with the log data
	http.Post("http://your-logging-database-url", "application/json", bytes.NewBuffer(jsonLogData))
}

func handleConnections(ws *websocket.Conn) {
	for {
		// Read in a new message as JSON and map it to a Message object
		var msg Input
		err := ws.ReadJSON(&msg)
		if err != nil {
			logError(err)
			break
		}

		// Print the message to the console
		log.Printf("%+v\n", msg)
	}
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func sendRequest(text string, ws *websocket.Conn, sessionToken string) {
	// Store the input in a struct
	input := Input{Text: text, SessionID: sessionToken}

	// Convert the struct to JSON
	jsonInput, err := json.Marshal(input)
	if err != nil {
		logError(err)
		return
	}
	// Print debug to check if the output is correct
	fmt.Println("JSON to be sent:", string(jsonInput))

	// Send the request via the websocket
	err = ws.WriteJSON(input)
	if err != nil {
		logError(err)
		return
	}
}

func Initializer() {
	// Connect to the WebSocket server
	ws, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:8080/chat", nil)
	if err != nil {
		log.Fatal("dial: ", err)
		logError(err)
	}
	defer ws.Close()

	// Mock of the session token
	sessionToken := "mockSessionToken"

	// loop this and end this with keyword
	reader := bufio.NewReader(os.Stdin)
	exitKeywords := []string{"Ende", "Exit", "Quit", "Beenden"} // Add more keywords as needed

	firstIteration := true
	for {
		if firstIteration {
			fmt.Println("Hello, how can I help you today?")
			firstIteration = false
		}
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		if contains(exitKeywords, text) {
			fmt.Println("Ending the program.")
			break
		}
		sendRequest(text, ws, sessionToken) // Call the function to send the request
	}
}