package tokenizer

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type Message struct {
	Text string `json:"text"`
}

func Tokenize(s string) []string {
	re := regexp.MustCompile(`[\p{L}\d_]+`)
	tokens := re.FindAllString(s, -1)
	return tokens
}

type TokenResponse struct {
	Tokens []string `json:"tokens"`
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Missing query parameter", http.StatusBadRequest)
		return
	}

	tokens := Tokenize(query)

	joinedTokens := strings.Join(tokens, "")

	resp, err := http.Get("/match?input=" + joinedTokens)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func TokenInit() {
	http.HandleFunc("/tokenize", handleRequest)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
