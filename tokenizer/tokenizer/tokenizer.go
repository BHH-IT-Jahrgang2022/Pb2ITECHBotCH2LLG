package tokenizer

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

type TokenResponse struct {
	Tokens []string `json:"tokens"`
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var msg Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tokens := Tokenize(msg.Text)
	resp := TokenResponse{Tokens: tokens}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func TokenInit() {
	http.HandleFunc("/tokenize", handleRequest)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
