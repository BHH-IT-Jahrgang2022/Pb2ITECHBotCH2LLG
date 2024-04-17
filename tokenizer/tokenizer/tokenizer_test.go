package tokenizer

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func TestTokenize(t *testing.T) {
	// Mock the matcher server
	matcherServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"response": "hallo auch"}`)
	}))
	defer matcherServer.Close()

	os.Setenv("MATCHER_URL", matcherServer.URL)
	// Replace the matcher URL with the mock server's URL
	var matcherURL string
	oldMatcherURL := matcherURL
	matcherURL = matcherServer.URL
	defer func() { matcherURL = oldMatcherURL }()

	// Create a request to pass to our handler
	query := "hello world!@#$%^&*()"
	escapedQuery := url.QueryEscape(query)
	req, err := http.NewRequest("GET", "http://localhost:8080/tokenize?query="+escapedQuery, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleRequest)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := `{"response": "hallo auch"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
