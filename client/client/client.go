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
	"github.com/google/uuid"
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

	// Send a POST request with the log data
	//http.Post("https://api.bot.demo.pinguin-it.de", "application/json", bytes.NewBuffer(jsonLogData))
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

// String examiner function used for matching keywords later in the text
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
	//fmt.Println("JSON to be sent:", string(jsonInput))

	// Send the request via the websocket
	err = ws.WriteJSON(input)
	if err != nil {
		logError(err)
		return
	}

	
	var response Input
	err = ws.ReadJSON(&response)
	if err != nil {
		logError(err)
		return
	}
	// Print the response to the console
	fmt.Printf(response)

}

func Initializer() {
	// Connect to the WebSocket server
	wsURL := "ws://" + os.Getenv("WSHOST") + ":" + os.Getenv("WSPORT") + "/chat"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("dial: ", err)
		logError(err)
	}
	defer ws.Close()

	// Creation of the session token
	sessionToken := uuid.New().String()

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
		text = strings.Replace(text, "\n", "", -1)

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
		sendRequest(text, ws, sessionToken) // Call the function to send the request
	}
}