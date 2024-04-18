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

	matcherURL := os.Getenv("MATCHER_URL")

	if matcherURL == "" {
		matcherURL = "http://localhost:8080" // Default-Wert
	}

	resp, err := http.Get(matcherURL + "/match?input=" + joinedTokens)

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

	var jsonData map[string]interface{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		http.Error(w, "Received invalid JSON from matcher", http.StatusInternalServerError)
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
