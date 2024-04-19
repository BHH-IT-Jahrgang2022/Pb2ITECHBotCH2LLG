package tokenizer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

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
	// parses the query parameter from the URL into a string and checks if it is empty
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Missing query parameter", http.StatusBadRequest)
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
	resp, err := http.Get(matcher_URL + "/match?input=" + joinedTokens)

	// checks if the request to the matcher failed
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// defers the closing of the response body
	defer resp.Body.Close()

	// reads the response from the matcher and prints it to the console
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Received response from matcher:", string(body))

	// unmarshals the JSON response from the matcher and checks if it is valid
	var jsonData map[string]interface{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		http.Error(w, "Received invalid JSON from matcher", http.StatusInternalServerError)
		return
	}

	// checks if the JSON contains the expected attribute "response"
	if _, ok := jsonData["response"]; !ok {
		http.Error(w, "Missing expected attribute in JSON from matcher", http.StatusBadRequest)
		return
	}

	// checks the number of attributes in the JSON (expected: 1)
	if len(jsonData) != 1 {
		http.Error(w, fmt.Sprintf("Expected %d attributes, got %d", 1, len(jsonData)), http.StatusBadRequest)
		return
	}

	// marshals the tokens into a JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func TokenInit() {
	http.HandleFunc("/tokenize", handleRequest)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
