package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
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

	// Open the log file
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a new logger
	logger := log.New(file, "", log.LstdFlags)

	// Write the log data to the file
	logger.Println(string(jsonLogData))

}

func handleConnections(ws *websocket.Conn) {
	for {
		// Read in a new message as bytes
		_, msgBytes, err := ws.ReadMessage()
		if err != nil {
			logError(err)
			break
		}

		// Convert the bytes to a string
		msg := string(msgBytes)

		// Print the message to the console
		log.Printf("%s\n", msg)
	}
}

// String examiner function used for matching keywords later in the text
func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func sendRequest(text string, ws *websocket.Conn) {
	// Convert the text to bytes
	textBytes := []byte(text)

	// Send the request via the websocket
	err := ws.WriteMessage(websocket.TextMessage, textBytes)
	if err != nil {
		logError(err)
		return
	}

	// Create a timer that will trigger after 30 seconds
	timer := time.NewTimer(30 * time.Second)

	// Create a channel to signal when a response is received
	done := make(chan bool)

	go func() {
		_, responseBytes, err := ws.ReadMessage()
		if err != nil {
			logError(err)
			return
		}
		// Convert the response bytes to a string
		response := string(responseBytes)

		// Print the response to the console
		fmt.Printf("%s\n", response)

		// Signal that a response was received
		done <- true
	}()

	// Wait for either the response or the timer
	select {
	case <-done:
		// A response was received, so stop the timer
		timer.Stop()
	case <-timer.C:
		// The timer triggered, so print an error message
		fmt.Println("Error: No response received after 30 seconds")
	}
}

func Initializer() {
	// Connect to the WebSocket server
	wsURL := "wss://" + os.Getenv("WSHOST") + ":" + os.Getenv("WSPORT") + "/chat"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("dial: ", err)
		logError(err)
	}
	defer ws.Close()

	// loop this and end this with keyword
	reader := bufio.NewReader(os.Stdin)
	exitKeywords := []string{"Ende", "Exit", "Quit", "Beenden", "Bye", "Das war alles"} // Add more keywords as needed

	firstIteration := true
	for {
		if firstIteration {
			fmt.Println("Hello, how can I help you today?")
			firstIteration = false
		}
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Trim(text, "\n")

		//text = strings.TrimSpace(text)

		// Check if the string only contains printable characters
		/*for _, r := range text {
			if !unicode.IsPrint(r) {
				fmt.Println("Invalid input. Please enter a string.")
				continue
			}
		}*/

		if contains(exitKeywords, text) {
			fmt.Println("Danke für ihre Anfrage")
			break
		}
		sendRequest(text, ws) // Call the function to send the request
	}
}
