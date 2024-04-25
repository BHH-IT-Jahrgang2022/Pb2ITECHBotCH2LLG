package tokenizer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

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

type Message struct {
	Text string `json:"text"`
}

func Tokenize(s string) []string {
	re := regexp.MustCompile(`[\p{L}\d_]+`)
	tokens := re.FindAllString(s, -1)
	if tokens == nil {
		tokens = []string{}
	}
	return tokens
}

type TokenResponse struct {
	Tokens []string `json:"tokens"`
}

func handleRequest(w http.ResponseWriter, r *http.Request) {

	// checks if the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "ERROR", Message: "Invalid request method", Service: "tokenizer"})
		return
	}
	// parses the query parameter from the URL into a string and checks if it is empty
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Missing query parameter", http.StatusBadRequest)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "ERROR", Message: "Missing query parameter", Service: "tokenizer"})
		return
	}

	// tokenizes the query
	tokens := Tokenize(query)

	// joins the tokens into a single string
	joinedTokens := strings.Join(tokens, "")

	// sets the matcher URL from the environment variable MATCHER_URL
	matcher_URL := os.Getenv("MATCHER_URL")

	// checks if the matcher URL is empty and sets a default value
	if matcher_URL == "" {
		matcher_URL = "http://localhost:8080" // (default value)
	}

	// sends a GET request to the matcher with the joined tokens as the input
	resp, err := http.Get(matcher_URL + "/match?input=" + joinedTokens + "&query=" + query)

	// checks if the request to the matcher failed
	if err != nil || resp == nil {
		http.Error(w, "GET request to matcher failed", http.StatusInternalServerError)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "ERROR", Message: "GET request to matcher failed", Service: "tokenizer"})
		return
	}

	// defers the closing of the response body
	defer resp.Body.Close()

	// reads the response from the matcher and prints it to the console
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "ERROR", Message: err.Error(), Service: "tokenizer"})
		return
	}

	fmt.Println("Received response from matcher:", string(body))

	// unmarshals the JSON response from the matcher and checks if it is valid
	var jsonData map[string]interface{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		http.Error(w, "Received invalid JSON from matcher", http.StatusInternalServerError)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "ERROR", Message: err.Error(), Service: "tokenizer"})
		return
	}

	// checks if the JSON contains the expected attribute "response"
	if _, ok := jsonData["response"]; !ok {
		http.Error(w, "Missing expected attribute in JSON from matcher", http.StatusBadRequest)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "ERROR", Message: "Missing expected attribute in JSON from matcher", Service: "tokenizer"})
		return
	}

	// checks the number of attributes in the JSON (expected: 1)
	if len(jsonData) != 1 {
		http.Error(w, fmt.Sprintf("Expected %d attributes, got %d", 1, len(jsonData)), http.StatusBadRequest)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "ERROR", Message: fmt.Sprintf("Expected %d attributes, got %d", 1, len(jsonData)), Service: "tokenizer"})
		return
	}

	// marshals the tokens into a JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func TokenInit() {
	http.HandleFunc("/tokenize", handleRequest)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe("0.0.0.0:8080", nil)
}
