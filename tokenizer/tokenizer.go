package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

type Message struct {
	Text string `json:"text"`
}

func Tokenize(s string) []string {
	re := regexp.MustCompile(`[\p{L}\d_]+|[!?]`)
	tokens := re.FindAllString(s, -1)
	return tokens
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var msg Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tokens := Tokenize(msg.Text)
	err = json.NewEncoder(w).Encode(tokens)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/tokenize", handleRequest)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
